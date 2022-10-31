package resource

import (
	"context"
	"errors"
	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/maintenance-window/cmd/validation"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util/constants"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util/progress_event"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util/validator"
	log "github.com/sirupsen/logrus"
	mongodbatlas "go.mongodb.org/atlas/mongodbatlas"
)

func validateModel(event constants.Event, model *Model) *handler.ProgressEvent {
	return validator.ValidateModel(event, validation.ModelValidator{}, model)
}

func setup() {
	util.SetupLogger("mongodb-atlas-maintenance-window")
}

func (m Model) toAtlasModel() mongodbatlas.MaintenanceWindow {
	return mongodbatlas.MaintenanceWindow{
		DayOfWeek:            *m.DayOfWeek,
		HourOfDay:            m.HourOfDay,
		StartASAP:            m.StartASAP,
		AutoDeferOnceEnabled: m.AutoDeferOnceEnabled,
	}
}

func Create(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup()
	log.Infof("Create() currentModel:%+v", *currentModel)

	// Validation
	modelValidation := validateModel(constants.Create, currentModel)
	if modelValidation != nil {
		log.Debugf("Validation Error")
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
	var res *mongodbatlas.Response

	atlasModel := currentModel.toAtlasModel()
	startASP := false
	atlasModel.StartASAP = &startASP

	res, err = client.MaintenanceWindows.Update(context.Background(), *currentModel.GroupId, &atlasModel)

	if err != nil {
		log.Debugf("Create - error: %+v", err)
		return progress_events.GetFailedEventByResponse(err.Error(), res.Response), nil
	}
	log.Debugf("Atlas Client %v", client)

	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		ResourceModel:   *currentModel,
	}, nil
}

func Read(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	log.Debugf("Read() currentModel:%+v", currentModel)

	// Validation
	modelValidation := validateModel(constants.Read, currentModel)
	if modelValidation != nil {
		log.Debugf("Validation Error")
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

	maintenanceWindow, errorProgressEvent := get(client, *currentModel)
	if errorProgressEvent != nil {
		return *errorProgressEvent, nil
	}

	currentModel.AutoDeferOnceEnabled = maintenanceWindow.AutoDeferOnceEnabled
	currentModel.DayOfWeek = &maintenanceWindow.DayOfWeek
	currentModel.HourOfDay = maintenanceWindow.HourOfDay
	// Response
	event := handler.ProgressEvent{
		OperationStatus: handler.Success,
		ResourceModel:   currentModel,
	}
	return event, nil
}

func get(client *mongodbatlas.Client, currentModel Model) (*mongodbatlas.MaintenanceWindow, *handler.ProgressEvent) {
	maintenanceWindow, res, err := client.MaintenanceWindows.Get(context.Background(), *currentModel.GroupId)
	if err != nil {
		log.Debugf("Read - error: %+v", err)
		ev := progress_events.GetFailedEventByResponse(err.Error(), res.Response)
		return nil, &ev
	}

	if isResponseEmpty(maintenanceWindow) {
		log.Debugf("Read - resource is empty: %+v", err)
		ev := progress_events.GetFailedEventByCode("resource not found", cloudformation.HandlerErrorCodeNotFound)
		return nil, &ev
	}

	return maintenanceWindow, nil
}

func isResponseEmpty(maintenanceWindow *mongodbatlas.MaintenanceWindow) bool {
	return (maintenanceWindow != nil) && (maintenanceWindow != nil && maintenanceWindow.DayOfWeek == 0)
}

func Update(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	log.Debugf("Update() currentModel:%+v", currentModel)

	// Validation
	modelValidation := validateModel(constants.Update, currentModel)
	if modelValidation != nil {
		log.Debugf("Validation Error")
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

	_, handlerError := get(client, *currentModel)
	if handlerError != nil {
		return *handlerError, nil
	}

	var res *mongodbatlas.Response

	atlasModel := currentModel.toAtlasModel()
	startASP := false
	atlasModel.StartASAP = &startASP

	res, err = client.MaintenanceWindows.Update(context.Background(), *currentModel.GroupId, &atlasModel)

	if err != nil {
		log.Debugf("Update - error: %+v", err)
		return progress_events.GetFailedEventByResponse(err.Error(), res.Response), nil
	}
	log.Debugf("Atlas Client %v", client)

	// Response
	event := handler.ProgressEvent{
		OperationStatus: handler.Success,
		ResourceModel:   currentModel,
	}
	return event, nil
}

func Delete(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	log.Debugf("Delete() currentModel:%+v", currentModel)

	// Validation
	modelValidation := validateModel(constants.Delete, currentModel)
	if modelValidation != nil {
		log.Debugf("Validation Error")
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

	_, handlerError := get(client, *currentModel)
	if handlerError != nil {
		return *handlerError, nil
	}

	var res *mongodbatlas.Response
	res, err = client.MaintenanceWindows.Reset(context.Background(), *currentModel.GroupId)

	if err != nil {
		log.Debugf("Delete - error: %+v", err)
		return progress_events.GetFailedEventByResponse(err.Error(), res.Response), nil
	}
	log.Debugf("Atlas Client %v", client)

	// Response
	event := handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "delete successful",
	}
	return event, nil
}

func List(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	return handler.ProgressEvent{}, errors.New("Not implemented: Update")
}
