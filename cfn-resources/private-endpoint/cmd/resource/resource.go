package resource

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/private-endpoint/cmd/resource/steps/aws_vpc_endpoint"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/private-endpoint/cmd/resource/steps/private_endpoint"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/private-endpoint/cmd/resource/steps/private_endpoint_service"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	resource_constats "github.com/mongodb/mongodbatlas-cloudformation-resources/private-endpoint/cmd/constants"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/private-endpoint/cmd/validator_def"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util/constants"
	progress_events "github.com/mongodb/mongodbatlas-cloudformation-resources/util/progress_event"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util/validator"
	"go.mongodb.org/atlas/mongodbatlas"
)

const (
	providerName = "AWS"
)

func setup() {
	util.SetupLogger("mongodb-atlas-private-endpoint")
}

func validateModel(event constants.Event, model *Model) *handler.ProgressEvent {
	return validator.ValidateModel(event, validator_def.ModelValidator{}, model)
}

func getProcessStatus(req handler.Request) (resource_constats.EventStatus, *handler.ProgressEvent) {
	callback, _ := req.CallbackContext["StateName"]
	if callback == nil {
		return resource_constats.CreationInit, nil
	}

	eventStatus, err := resource_constats.ParseEventStatus(fmt.Sprintf("%v", callback))

	if err != nil {
		pe := progress_events.GetFailedEventByCode(fmt.Sprintf("Error parsing callback status : %s", err.Error()), cloudformation.HandlerErrorCodeServiceInternalError)
		return "", &pe
	}

	return eventStatus, nil
}

func (m *Model) completeByConnection(c mongodbatlas.PrivateEndpointConnection) {
	m.Id = &c.ID
	m.EndpointServiceName = &c.EndpointServiceName
	m.ErrorMessage = &c.ErrorMessage
	m.InterfaceEndpoints = c.InterfaceEndpoints
	m.Status = &c.Status
}

func addModelToProgressEvent(progressEvent *handler.ProgressEvent, model *Model) handler.ProgressEvent {
	if progressEvent.OperationStatus == handler.InProgress {
		progressEvent.ResourceModel = model

		callbackId, _ := progressEvent.CallbackContext["Id"]

		if callbackId != nil {
			id := fmt.Sprint(callbackId)
			model.Id = &id
		}

	}

	return *progressEvent
}

// Create handles the Create event from the Cloudformation service.
func Create(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup()
	log.Info("Initiated Create")
	modelValidation := validateModel(constants.Create, currentModel)
	if modelValidation != nil {
		return *modelValidation, nil
	}

	mongodbClient, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		return progress_events.GetFailedEventByCode(fmt.Sprintf("Error creating mongoDB client : %s", err.Error()),
			cloudformation.HandlerErrorCodeInvalidRequest), nil
	}

	status, pe := getProcessStatus(req)
	if pe != nil {
		return *pe, nil
	}

	log.Infof("Status recieved %s", status)

	switch status {
	case resource_constats.CreationInit:
		pe := private_endpoint_service.CreatePrivateEndpoint(*mongodbClient, *currentModel.Region, *currentModel.GroupId)
		return addModelToProgressEvent(&pe, currentModel), nil
	case resource_constats.CreatingPrivateEndpointService:
		peConnection, completionValidation := private_endpoint_service.ValidateCreationCompletion(mongodbClient, *currentModel.GroupId, req)
		if completionValidation != nil {
			return addModelToProgressEvent(completionValidation, currentModel), nil
		}

		vpcEndpointId, progressEvent := aws_vpc_endpoint.CreateVpcEndpoint(*peConnection, *currentModel.Region, *currentModel.SubnetId, *currentModel.VpcId)
		if progressEvent != nil {
			return addModelToProgressEvent(progressEvent, currentModel), nil
		}

		pe := private_endpoint.CreatePrivateEndpoint(mongodbClient, *currentModel.GroupId, *vpcEndpointId, peConnection.ID)

		return addModelToProgressEvent(&pe, currentModel), nil
	default:
		ValidationOutput, progressEvent := private_endpoint.ValidateCreationCompletion(mongodbClient, *currentModel.GroupId, req)
		if progressEvent != nil {
			return addModelToProgressEvent(progressEvent, currentModel), nil
		}
		currentModel.Id = &ValidationOutput.Id
		currentModel.InterfaceEndpoints = ValidationOutput.InterfaceEndpoints
		return handler.ProgressEvent{
			OperationStatus: handler.Success,
			Message:         "Create Completed",
			ResourceModel:   currentModel}, nil
	}
}

func deleteVcpEndpoints(currentModel *Model) (*ec2.DeleteVpcEndpointsOutput, *handler.ProgressEvent) {
	mySession := session.Must(session.NewSession())

	// Create a EC2 client from just a session.
	svc := ec2.New(mySession, aws.NewConfig().WithRegion("us-east-1"))

	subnetIds := currentModel.InterfaceEndpoints
	vpcEndpointIds := make([]*string, 0)

	for _, i := range subnetIds {
		vpcEndpointIds = append(vpcEndpointIds, &i)
	}

	connection := ec2.DeleteVpcEndpointsInput{
		DryRun:         nil,
		VpcEndpointIds: vpcEndpointIds,
	}

	//vpcE, err := svc.CreateVpcEndpoint(&connection)
	vpcE, err := svc.DeleteVpcEndpoints(&connection)
	if err != nil {
		fpe := handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          fmt.Sprintf("Error deleting vcp Endpoint: %s", err.Error()),
			HandlerErrorCode: cloudformation.HandlerErrorCodeGeneralServiceException}
		return nil, &fpe
	}

	return vpcE, nil
}

// Read handles the Read event from the Cloudformation service.
func Read(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup()
	modelValidation := validateModel(constants.Read, currentModel)
	if modelValidation != nil {
		return *modelValidation, nil
	}
	mongodbClient, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		return progress_events.GetFailedEventByCode(fmt.Sprintf("Error creating mongoDB client : %s", err.Error()),
			cloudformation.HandlerErrorCodeInvalidRequest), nil
	}

	privateEndpointResponse, response, err := mongodbClient.PrivateEndpoints.Get(context.Background(), *currentModel.GroupId, providerName, *currentModel.Id)
	if err != nil {
		return progress_events.GetFailedEventByResponse(fmt.Sprintf("Error getting resource : %s", err.Error()),
			response.Response), nil
	}

	currentModel.completeByConnection(*privateEndpointResponse)

	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "Get successful",
		ResourceModel:   currentModel}, nil
}

// Update handles the Update event from the Cloudformation service.
func Update(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	return handler.ProgressEvent{}, errors.New("Not implemented: Update")
}

// Delete handles the Delete event from the Cloudformation service.
func Delete(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup()
	modelValidation := validateModel(constants.Delete, currentModel)
	if modelValidation != nil {
		return *modelValidation, nil
	}

	mongodbClient, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		return progress_events.GetFailedEventByCode(fmt.Sprintf("Error creating mongoDB client : %s", err.Error()),
			cloudformation.HandlerErrorCodeNotFound), nil
	}

	privateEndpointResponse, response, err := mongodbClient.PrivateEndpoints.Get(context.Background(), *currentModel.GroupId, providerName, *currentModel.Id)

	callback, _ := req.CallbackContext["stateName"]
	if callback != nil {
		callbackValue := fmt.Sprintf("%v", callback)
		if callbackValue == "DELETING" {
			if response.StatusCode == http.StatusNotFound {
				return handler.ProgressEvent{
					OperationStatus: handler.Success,
					Message:         "Delete success"}, nil
			}

			return handler.ProgressEvent{
				OperationStatus:      handler.InProgress,
				Message:              "Delete in progress",
				ResourceModel:        currentModel,
				CallbackDelaySeconds: 20,
				CallbackContext: map[string]interface{}{
					"stateName": "DELETING",
				}}, nil
		}
	}

	if err != nil {
		return progress_events.GetFailedEventByResponse(fmt.Sprintf("Error getting resource : %s", err.Error()),
			response.Response), nil
	}

	currentModel.completeByConnection(*privateEndpointResponse)

	if len(currentModel.InterfaceEndpoints) != 0 {
		for _, intEndpoints := range currentModel.InterfaceEndpoints {

			//delete the private endpoint
			response, err := mongodbClient.PrivateEndpoints.DeleteOnePrivateEndpoint(context.Background(),
				*currentModel.GroupId,
				providerName,
				*currentModel.Id,
				intEndpoints)
			if err != nil {
				return progress_events.GetFailedEventByResponse(fmt.Sprintf("Error deleting resource : %s", err.Error()),
					response.Response), nil
			}
		}
		_, epr := deleteVcpEndpoints(currentModel)

		if epr != nil {
			return *epr, nil
		}

	} else {
		response, err = mongodbClient.PrivateEndpoints.Delete(context.Background(), *currentModel.GroupId,
			providerName,
			*currentModel.Id)

		if err != nil {
			return progress_events.GetFailedEventByResponse(fmt.Sprintf("Error getting resource : %s", err.Error()),
				response.Response), nil
		}
	}

	return handler.ProgressEvent{
		OperationStatus:      handler.InProgress,
		Message:              "Delete in progress",
		ResourceModel:        currentModel,
		CallbackDelaySeconds: 20,
		CallbackContext: map[string]interface{}{
			"stateName":         "DELETING",
			"AwsVpcEndpointIds": currentModel.InterfaceEndpoints,
		}}, nil
}

// List handles the List event from the Cloudformation service.
func List(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup()
	modelValidation := validateModel(constants.List, currentModel)
	if modelValidation != nil {
		return *modelValidation, nil
	}

	mongodbClient, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		return progress_events.GetFailedEventByCode(fmt.Sprintf("Error creating mongoDB client : %s", err.Error()),
			cloudformation.HandlerErrorCodeInvalidRequest), nil
	}

	params := &mongodbatlas.ListOptions{
		PageNum:      0,
		ItemsPerPage: 100,
	}

	privateEndpointResponse, response, err := mongodbClient.PrivateEndpoints.List(context.Background(),
		*currentModel.GroupId,
		providerName,
		params)
	if err != nil {
		return progress_events.GetFailedEventByResponse(fmt.Sprintf("Error listing resource : %s", err.Error()),
			response.Response), nil
	}

	mm := make([]interface{}, 0)
	for _, privateEndpoint := range privateEndpointResponse {
		var m Model
		m.completeByConnection(privateEndpoint)
		mm = append(mm, m)
	}

	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "List successful",
		ResourceModels:  mm}, nil
}
