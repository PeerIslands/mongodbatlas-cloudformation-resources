package resource

import (
	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util/progress_event"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util/validator"
	log "github.com/sirupsen/logrus"
	mongodbatlas "go.mongodb.org/atlas/mongodbatlas"
)

var CreateRequiredFields = []string{"ApiKeys.PrivateKey", "ApiKeys.PublicKey"}
var ReadRequiredFields = []string{"currentModel.GroupId", "ApiKeys.PrivateKey", "ApiKeys.PublicKey"}
var UpdateRequiredFields = []string{"currentModel.DayOfWeek", "currentModel.HourOfDay", "ApiKeys.PrivateKey", "ApiKeys.PublicKey"}
var DeleteRequiredFields = []string{"currentModel.GroupId", "ApiKeys.PrivateKey", "ApiKeys.PublicKey"}
var ListRequiredFields = []string{"ApiKeys.PrivateKey", "ApiKeys.PublicKey"}

func validateModel(fields []string, model *Model) *handler.ProgressEvent {
	return validator.ValidateModel(fields, model)
}
func Create(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	log.Debugf("Create() currentModel:%+v", currentModel)

	// Validation
	modelValidation := validateModel(CreateRequiredFields, currentModel)
	if modelValidation != nil {
		return *modelValidation, nil
	}

	// Create atlas client
	client, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		log.Debugf("Create - error: %+v", err)
		return handler.ProgressEvent{
			HandlerErrorCode: cloudformation.HandlerErrorCodeInvalidRequest,
			Message:          err.Error(),
			OperationStatus:  handler.Failed,
		}, nil
	}
	var res mongodbatlas.Response

	//
	/*
	    // Pseudocode:
	    res , resModel, err := client.maintenancewindow.Create(context.Background(),&mongodbatlas.Maintenancewindow{
	   })

	*/

	if err != nil {
		log.Debugf("Create - error: %+v", err)
		return progress_events.GetFailedEventByResponse(err.Error(), res.Response), nil
	}
	log.Debugf("Atlas Client %v", client)

	// Response
	event := handler.ProgressEvent{
		OperationStatus: handler.InProgress,
		ResourceModel:   currentModel,
	}
	return event, nil
}

func Read(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	log.Debugf("Read() currentModel:%+v", currentModel)

	// Validation
	modelValidation := validateModel(CreateRequiredFields, currentModel)
	if modelValidation != nil {
		return *modelValidation, nil
	}

	// Create atlas client
	client, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		log.Debugf("Read - error: %+v", err)
		return handler.ProgressEvent{
			HandlerErrorCode: cloudformation.HandlerErrorCodeInvalidRequest,
			Message:          err.Error(),
			OperationStatus:  handler.Failed,
		}, nil
	}
	var res mongodbatlas.Response

	/*
	   Considerable params from currentModel:
	   GroupId, ...
	*/
	/*
	    // Pseudocode:
	    res , resModel, err := client.maintenancewindow.Read(context.Background(),&mongodbatlas.Maintenancewindow{
	   	GroupId:currentModel.GroupId,
	   })

	*/

	if err != nil {
		log.Debugf("Read - error: %+v", err)
		return progress_events.GetFailedEventByResponse(err.Error(), res.Response), nil
	}
	log.Debugf("Atlas Client %v", client)

	// Response
	event := handler.ProgressEvent{
		OperationStatus: handler.InProgress,
		ResourceModel:   currentModel,
	}
	return event, nil
}

func Update(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	log.Debugf("Update() currentModel:%+v", currentModel)

	// Validation
	modelValidation := validateModel(CreateRequiredFields, currentModel)
	if modelValidation != nil {
		return *modelValidation, nil
	}

	// Create atlas client
	client, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		log.Debugf("Update - error: %+v", err)
		return handler.ProgressEvent{
			HandlerErrorCode: cloudformation.HandlerErrorCodeInvalidRequest,
			Message:          err.Error(),
			OperationStatus:  handler.Failed,
		}, nil
	}
	var res mongodbatlas.Response

	/*
	   Considerable params from currentModel:
	   HourOfDay, StartASAP, ApiKeys, GroupId, AutoDeferOnceEnabled, DayOfWeek, ...
	*/
	/*
	    // Pseudocode:
	    res , resModel, err := client.maintenancewindow.Update(context.Background(),&mongodbatlas.Maintenancewindow{
	   	HourOfDay:currentModel.HourOfDay,
	   	StartASAP:currentModel.StartASAP,
	   	ApiKeys:currentModel.ApiKeys,
	   	GroupId:currentModel.GroupId,
	   	AutoDeferOnceEnabled:currentModel.AutoDeferOnceEnabled,
	   	DayOfWeek:currentModel.DayOfWeek,
	   })

	*/

	if err != nil {
		log.Debugf("Update - error: %+v", err)
		return progress_events.GetFailedEventByResponse(err.Error(), res.Response), nil
	}
	log.Debugf("Atlas Client %v", client)

	// Response
	event := handler.ProgressEvent{
		OperationStatus: handler.InProgress,
		ResourceModel:   currentModel,
	}
	return event, nil
}

func Delete(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	log.Debugf("Delete() currentModel:%+v", currentModel)

	// Validation
	modelValidation := validateModel(CreateRequiredFields, currentModel)
	if modelValidation != nil {
		return *modelValidation, nil
	}

	// Create atlas client
	client, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		log.Debugf("Delete - error: %+v", err)
		return handler.ProgressEvent{
			HandlerErrorCode: cloudformation.HandlerErrorCodeInvalidRequest,
			Message:          err.Error(),
			OperationStatus:  handler.Failed,
		}, nil
	}
	var res mongodbatlas.Response

	/*
	   Considerable params from currentModel:
	   GroupId, ...
	*/
	/*
	    // Pseudocode:
	    res , resModel, err := client.maintenancewindow.Delete(context.Background(),&mongodbatlas.Maintenancewindow{
	   	GroupId:currentModel.GroupId,
	   })

	*/

	if err != nil {
		log.Debugf("Delete - error: %+v", err)
		return progress_events.GetFailedEventByResponse(err.Error(), res.Response), nil
	}
	log.Debugf("Atlas Client %v", client)

	// Response
	event := handler.ProgressEvent{
		OperationStatus: handler.InProgress,
		ResourceModel:   currentModel,
	}
	return event, nil
}

func List(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	log.Debugf("List() currentModel:%+v", currentModel)

	// Validation
	modelValidation := validateModel(CreateRequiredFields, currentModel)
	if modelValidation != nil {
		return *modelValidation, nil
	}

	// Create atlas client
	client, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		log.Debugf("List - error: %+v", err)
		return handler.ProgressEvent{
			HandlerErrorCode: cloudformation.HandlerErrorCodeInvalidRequest,
			Message:          err.Error(),
			OperationStatus:  handler.Failed,
		}, nil
	}
	var res mongodbatlas.Response

	//
	/*
	    // Pseudocode:
	    res , resModel, err := client.maintenancewindow.List(context.Background(),&mongodbatlas.Maintenancewindow{
	   })

	*/

	if err != nil {
		log.Debugf("List - error: %+v", err)
		return progress_events.GetFailedEventByResponse(err.Error(), res.Response), nil
	}
	log.Debugf("Atlas Client %v", client)

	// Response
	event := handler.ProgressEvent{
		OperationStatus: handler.InProgress,
		ResourceModel:   currentModel,
	}
	return event, nil
}
