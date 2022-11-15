package private_endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/private-endpoint/cmd/constants"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/private-endpoint/cmd/resource"
	progress_events "github.com/mongodb/mongodbatlas-cloudformation-resources/util/progress_event"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util/structs"
	"go.mongodb.org/atlas/mongodbatlas"
)

const (
	ProviderName            = "AWS"
	StatusPendingAcceptance = "PENDING_ACCEPTANCE"
	StatusPending           = "PENDING"
	StatusAvailable         = "AVAILABLE"
)

type privateEndpointCreationCallBackContext struct {
	StateName   constants.EventStatus
	Id          string
	InterfaceId string
}

func (s *privateEndpointCreationCallBackContext) FillStruct(m map[string]interface{}) error {
	for k, v := range m {
		err := structs.SetField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreatePrivateEndpoint(mongodbClient *mongodbatlas.Client, currentModel resource.Model, interfaceEndpointID string, endpointServiceID string) handler.ProgressEvent {
	interfaceEndpointRequest := &mongodbatlas.InterfaceEndpointConnection{
		ID: interfaceEndpointID,
	}

	_, response, err := mongodbClient.PrivateEndpoints.AddOnePrivateEndpoint(context.Background(),
		*currentModel.GroupId,
		ProviderName,
		endpointServiceID,
		interfaceEndpointRequest)
	if err != nil {
		return progress_events.GetFailedEventByResponse(fmt.Sprintf("Error creating resource : %s", err.Error()),
			response.Response)
	}

	callBackContext := privateEndpointCreationCallBackContext{
		StateName:   constants.CreatingPrivateEndpointService,
		Id:          endpointServiceID,
		InterfaceId: interfaceEndpointID,
	}

	var callBackMap map[string]interface{}
	data, _ := json.Marshal(callBackContext)
	json.Unmarshal(data, &callBackMap)

	return progress_events.GetInProgressProgressEvent("Creating private endpoint service", currentModel, callBackMap)
}

func ValidateCreationCompletion(mongodbClient *mongodbatlas.Client, currentModel *resource.Model, req handler.Request) handler.ProgressEvent {

	callBackContext := privateEndpointCreationCallBackContext{}

	err := callBackContext.FillStruct(req.CallbackContext)
	if err != nil {
		return progress_events.GetFailedEventByCode(fmt.Sprintf("Error parsing PrivateEndpointCallBackContext : %s", err.Error()), cloudformation.HandlerErrorCodeServiceInternalError)
	}

	privateEndpointResponse, response, err := mongodbClient.PrivateEndpoints.GetOnePrivateEndpoint(context.Background(),
		*currentModel.GroupId,
		ProviderName,
		callBackContext.Id,
		callBackContext.InterfaceId)
	if err != nil {
		return progress_events.GetFailedEventByResponse(fmt.Sprintf("Error getting resource : %s", err.Error()),
			response.Response)
	}

	switch privateEndpointResponse.AWSConnectionStatus {
	case StatusPendingAcceptance, StatusPending:
		{
			return progress_events.GetInProgressProgressEvent("Private endpoint service initiating", currentModel, req.CallbackContext)
		}
	case StatusAvailable:
		{
			currentModel.Id = &callBackContext.Id
			currentModel.InterfaceEndpoints = []string{callBackContext.InterfaceId}
			return handler.ProgressEvent{
				OperationStatus: handler.Success,
				Message:         "Create Completed",
				ResourceModel:   currentModel}
		}
	}

	return handler.ProgressEvent{
		OperationStatus:  handler.Failed,
		Message:          fmt.Sprintf("Resource is in status : %s", privateEndpointResponse.AWSConnectionStatus),
		HandlerErrorCode: cloudformation.HandlerErrorCodeAlreadyExists}
}
