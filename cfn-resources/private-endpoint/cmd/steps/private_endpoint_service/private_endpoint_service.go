package private_endpoint_service

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
	"net/http"
)

const (
	ProviderName     = "AWS"
	AvailableStatus  = "AVAILABLE"
	InitiatingStatus = "INITIATING"
)

type privateEndpointCreationCallBackContext struct {
	StateName constants.EventStatus
	Id        string
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

func CreatePrivateEndpoint(mongodbClient mongodbatlas.Client, currentModel *resource.Model) handler.ProgressEvent {

	privateEndpointRequest := &mongodbatlas.PrivateEndpointConnection{
		ProviderName: ProviderName,
		Region:       *currentModel.Region,
	}

	privateEndpointResponse, response, err := mongodbClient.PrivateEndpoints.Create(context.Background(),
		*currentModel.GroupId,
		privateEndpointRequest)

	if response.Response.StatusCode == http.StatusConflict {
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          "Resource already exists",
			HandlerErrorCode: cloudformation.HandlerErrorCodeAlreadyExists}
	}

	if err != nil {
		return progress_events.GetFailedEventByResponse(fmt.Sprintf("Error creating resource : %s", err.Error()),
			response.Response)
	}

	callBackContext := privateEndpointCreationCallBackContext{
		StateName: constants.CreatingPrivateEndpointService,
		Id:        privateEndpointResponse.ID,
	}

	var callBackMap map[string]interface{}
	data, _ := json.Marshal(callBackContext)
	json.Unmarshal(data, &callBackMap)

	return progress_events.GetInProgressProgressEvent("Creating private endpoint service", currentModel, callBackMap)
}

func ValidateCreationCompletion(mongodbClient *mongodbatlas.Client, currentModel *resource.Model, req handler.Request) (*mongodbatlas.PrivateEndpointConnection, *handler.ProgressEvent) {

	PrivateEndpointCallBackContext := privateEndpointCreationCallBackContext{}

	err := PrivateEndpointCallBackContext.FillStruct(req.CallbackContext)
	if err != nil {
		ev := progress_events.GetFailedEventByCode(fmt.Sprintf("Error parsing PrivateEndpointCallBackContext : %s", err.Error()), cloudformation.HandlerErrorCodeServiceInternalError)
		return nil, &ev
	}

	privateEndpointResponse, response, err := mongodbClient.PrivateEndpoints.Get(context.Background(), *currentModel.GroupId, ProviderName, PrivateEndpointCallBackContext.Id)
	if err != nil {
		ev := progress_events.GetFailedEventByResponse(fmt.Sprintf("Error getting resource : %s", err.Error()),
			response.Response)
		return nil, &ev
	}

	if privateEndpointResponse.Status == InitiatingStatus {
		callBackContext := privateEndpointCreationCallBackContext{
			StateName: constants.CreatingPrivateEndpointService,
			Id:        privateEndpointResponse.ID,
		}

		var callBackMap map[string]interface{}
		data, _ := json.Marshal(callBackContext)

		json.Unmarshal(data, &callBackMap)
		ev := progress_events.GetInProgressProgressEvent("Private endpoint service initiating", currentModel, callBackMap)
		return nil, &ev
	} else if privateEndpointResponse.Status == AvailableStatus {
		return privateEndpointResponse, nil
	} else {
		ev := progress_events.GetFailedEventByCode(fmt.Sprintf("Error creating private endpoint in status : %s", privateEndpointResponse.Status),
			cloudformation.HandlerErrorCodeInvalidRequest)
		return nil, &ev
	}
}
