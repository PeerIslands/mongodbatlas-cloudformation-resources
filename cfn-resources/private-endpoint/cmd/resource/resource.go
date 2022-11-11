package resource

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/private-endpoint/cmd/validator_def"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util/constants"
	progress_events "github.com/mongodb/mongodbatlas-cloudformation-resources/util/progress_event"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util/validator"
	"go.mongodb.org/atlas/mongodbatlas"
)

const (
	providerName                   = "AWS"
	creatingPrivateEndpointService = "CREATING_PRIVATE_ENDPOINT_SERVICE"
	creatingVcpConnection          = "CREATING_PRIVATE_ENDPOINT_SERVICE"
	creatingPrivateEndpoint        = "CREATING_PRIVATE_ENDPOINT"
)

func setup() {
	util.SetupLogger("mongodb-atlas-private-endpoint")
}

func validateModel(event constants.Event, model *Model) *handler.ProgressEvent {
	return validator.ValidateModel(event, validator_def.ModelValidator{}, model)
}

func (m *Model) completeByConnection(c mongodbatlas.PrivateEndpointConnection) {
	m.Id = &c.ID
	m.EndpointServiceName = &c.EndpointServiceName
	m.ErrorMessage = &c.ErrorMessage
	m.InterfaceEndpoints = c.InterfaceEndpoints
	m.Status = &c.Status
}

// Create handles the Create event from the Cloudformation service.
func Create(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup()
	modelValidation := validateModel(constants.Create, currentModel)
	if modelValidation != nil {
		return *modelValidation, nil
	}

	mongodbClient, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		return progress_events.GetFailedEventByCode(fmt.Sprintf("Error creating mongoDB client : %s", err.Error()),
			cloudformation.HandlerErrorCodeInvalidRequest), nil
	}

	callback, _ := req.CallbackContext["stateName"]
	callbackValue := fmt.Sprintf("%v", callback)

	log.Infof("Callback value state name %v", callbackValue)

	switch callbackValue {
	case creatingPrivateEndpointService:
		callbackId, _ := req.CallbackContext["id"]
		serviceId := fmt.Sprintf("%v", callbackId)

		vcpe, fpe := validateAndCreateVCPConnection(mongodbClient, req, currentModel, serviceId)
		if fpe != nil {
			fpe.ResourceModel = vcpe
			return *fpe, nil
		}

		vpcEndpointId := *vcpe.VpcEndpoint.VpcEndpointId
		log.Infof("Attaching private endpoint interfaceID %v", vpcEndpointId)
		return attachPrivateEndpoint(mongodbClient, *currentModel, vpcEndpointId, serviceId), nil

	case creatingPrivateEndpoint:

		callbackId, _ := req.CallbackContext["id"]
		serviceId := fmt.Sprintf("%v", callbackId)

		callbackInterfaceId, _ := req.CallbackContext["interfaceID"]
		interfaceId := fmt.Sprintf("%v", callbackInterfaceId)

		privateEndpointResponse, response, err := mongodbClient.PrivateEndpoints.GetOnePrivateEndpoint(context.Background(),
			*currentModel.GroupId,
			providerName,
			serviceId,
			interfaceId)
		if err != nil {
			return progress_events.GetFailedEventByResponse(fmt.Sprintf("Error getting resource : %s", err.Error()),
				response.Response), nil
		}

		switch privateEndpointResponse.AWSConnectionStatus {
		case "PENDING_ACCEPTANCE", "PENDING":
			{
				return getInProgressProgressEvent(fmt.Sprintf("Adding private endpoint, status: %v", privateEndpointResponse.AWSConnectionStatus), currentModel,
					creatingPrivateEndpoint, serviceId, &interfaceId), nil
			}
		case "AVAILABLE":
			{
				currentModel.Id = &serviceId
				currentModel.InterfaceEndpoints = []string{interfaceId}
				return handler.ProgressEvent{
					OperationStatus: handler.Success,
					Message:         "Create Completed",
					ResourceModel:   currentModel}, nil
			}
		}

		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          fmt.Sprintf("Resource is in status : %s", privateEndpointResponse.AWSConnectionStatus),
			HandlerErrorCode: cloudformation.HandlerErrorCodeAlreadyExists}, nil
	}

	privateEndpointRequest := &mongodbatlas.PrivateEndpointConnection{
		ProviderName: providerName,
		Region:       *currentModel.Region,
	}

	privateEndpointResponse, response, err := mongodbClient.PrivateEndpoints.Create(context.Background(),
		*currentModel.GroupId,
		privateEndpointRequest)

	if response.Response.StatusCode == http.StatusConflict {
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          "Resource already exists",
			HandlerErrorCode: cloudformation.HandlerErrorCodeAlreadyExists}, nil
	}

	if err != nil {
		return progress_events.GetFailedEventByResponse(fmt.Sprintf("Error creating resource : %s", err.Error()),
			response.Response), nil
	}

	return getInProgressProgressEvent("Creating private endpoint service", currentModel,
		creatingPrivateEndpointService, privateEndpointResponse.ID, nil), nil
}

func getInProgressProgressEvent(message string, currentModel *Model, stateName string, privateEndpointServiceID string, interfaceConnection *string) handler.ProgressEvent {
	return handler.ProgressEvent{
		OperationStatus:      handler.InProgress,
		Message:              message,
		ResourceModel:        currentModel,
		CallbackDelaySeconds: 10,
		CallbackContext: map[string]interface{}{
			"stateName":   stateName,
			"id":          privateEndpointServiceID,
			"interfaceID": interfaceConnection,
		}}
}

func validateAndCreateVCPConnection(mongodbClient *mongodbatlas.Client, req handler.Request, currentModel *Model, serviceId string) (*ec2.CreateVpcEndpointOutput, *handler.ProgressEvent) {
	peCon, completionValidation := validatePrivateEndpointServiceCreationCompletion(mongodbClient, req, currentModel, serviceId)
	if completionValidation != nil {
		return nil, completionValidation
	}

	vcpEndpoint, progressEvent := createVcpEndpoint(*peCon, currentModel)
	if progressEvent != nil {
		return nil, progressEvent
	}

	return vcpEndpoint, nil
}

func createVcpEndpoint(peCon mongodbatlas.PrivateEndpointConnection, currentModel *Model) (*ec2.CreateVpcEndpointOutput, *handler.ProgressEvent) {
	mySession := session.Must(session.NewSession())

	// Create a EC2 client from just a session.
	svc := ec2.New(mySession, aws.NewConfig().WithRegion("us-east-1"))

	subnetIds := []*string{currentModel.SubnetId}

	vcpType := "Interface"

	connection := ec2.CreateVpcEndpointInput{
		VpcId:           currentModel.VpcId,
		ServiceName:     &peCon.EndpointServiceName,
		VpcEndpointType: &vcpType,
		SubnetIds:       subnetIds,
	}

	log.Infof("VpcId: %v", *connection.VpcId)
	log.Infof("ServiceName: %v", *connection.ServiceName)
	log.Infof("VpcEndpointType: %v", *connection.VpcEndpointType)
	log.Infof("SubnetIds: %v", *connection.SubnetIds[0])
	log.Infof("region: %v", *currentModel.Region)

	vpcE, err := svc.CreateVpcEndpoint(&connection)
	if err != nil {
		fpe := handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          fmt.Sprintf("Error creating vcp Endpoint: %s", err.Error()),
			HandlerErrorCode: cloudformation.HandlerErrorCodeGeneralServiceException}
		return nil, &fpe
	}

	log.Info("PRUEB: se creo correctamente el vcp endpoint")

	return vpcE, nil
}

func attachPrivateEndpoint(mongodbClient *mongodbatlas.Client, currentModel Model, interfaceEndpointID string, endpointServiceID string) handler.ProgressEvent {

	interfaceEndpointRequest := &mongodbatlas.InterfaceEndpointConnection{
		ID: interfaceEndpointID,
	}

	_, response, err := mongodbClient.PrivateEndpoints.AddOnePrivateEndpoint(context.Background(),
		*currentModel.GroupId,
		providerName,
		endpointServiceID,
		interfaceEndpointRequest)
	if err != nil {
		return progress_events.GetFailedEventByResponse(fmt.Sprintf("Error creating resource : %s", err.Error()),
			response.Response)
	}

	return getInProgressProgressEvent("Add private endpoint in progress", &currentModel,
		creatingPrivateEndpoint, endpointServiceID, &interfaceEndpointID)
}

func validatePrivateEndpointServiceCreationCompletion(mongodbClient *mongodbatlas.Client, req handler.Request, currentModel *Model, serviceId string) (*mongodbatlas.PrivateEndpointConnection, *handler.ProgressEvent) {

	privateEndpointResponse, response, err := mongodbClient.PrivateEndpoints.Get(context.Background(), *currentModel.GroupId, providerName, serviceId)
	if err != nil {
		ev := progress_events.GetFailedEventByResponse(fmt.Sprintf("Error getting resource : %s", err.Error()),
			response.Response)
		return nil, &ev
	}

	if privateEndpointResponse.Status == "INITIATING" {
		ev := getInProgressProgressEvent("Private endpoint service initiating", currentModel,
			creatingPrivateEndpointService, privateEndpointResponse.ID, nil)
		return nil, &ev
	} else if privateEndpointResponse.Status == "AVAILABLE" {
		return privateEndpointResponse, nil
	} else {
		ev := progress_events.GetFailedEventByCode(fmt.Sprintf("Error creating private endpoint in status : %s", privateEndpointResponse.Status),
			cloudformation.HandlerErrorCodeInvalidRequest)
		return nil, &ev
	}
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
	// Add your code here:
	// * Make API calls (use req.Session)
	// * Mutate the model
	// * Check/set any callback context (req.CallbackContext / response.CallbackContext)

	/*
	   // Construct a new handler.ProgressEvent and return it
	   response := handler.ProgressEvent{
	       OperationStatus: handler.Success,
	       Message: "Update complete",
	       ResourceModel: currentModel,
	   }

	   return response, nil
	*/

	// Not implemented, return an empty handler.ProgressEvent
	// and an error
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
	if err != nil {
		return progress_events.GetFailedEventByResponse(fmt.Sprintf("Error getting resource : %s", err.Error()),
			response.Response), nil
	}

	callback, _ := req.CallbackContext["stateName"]
	if callback != nil {
		callbackValue := fmt.Sprintf("%v", callback)
		if callbackValue == "DELETING" {
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

	currentModel.completeByConnection(*privateEndpointResponse)

	if len(currentModel.InterfaceEndpoints) == 0 {
		response, err := mongodbClient.PrivateEndpoints.Delete(context.Background(), *currentModel.GroupId,
			providerName,
			*currentModel.Id)

		if err != nil {

			return progress_events.GetFailedEventByResponse(fmt.Sprintf("Error getting resource : %s", err.Error()),
				response.Response), nil
		}

		return handler.ProgressEvent{
			OperationStatus:      handler.InProgress,
			Message:              "Delete in progress",
			ResourceModel:        currentModel,
			CallbackDelaySeconds: 20,
			CallbackContext: map[string]interface{}{
				"stateName": "DELETING",
			}}, nil
	} else {

	}
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
