package resource

import (
	"context"
	"errors"
	progress_events "github.com/mongodb/mongodbatlas-cloudformation-resources/util/progress_event"

	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/auditing/cmd/validation"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util/constants"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util/validator"
	log "github.com/sirupsen/logrus"
	mongodbatlas "go.mongodb.org/atlas/mongodbatlas"
)

func validateModel(event constants.Event, model *Model) *handler.ProgressEvent {
	return validator.ValidateModel(event, validation.ModelValidator{}, model)
}

func setup() {
	util.SetupLogger("mongodb-atlas-auditing")
}

func Create(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup()

	log.Infof("Create() currentModel:%+v", *currentModel.GroupId)

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

	enabled := true

	auditingInput := mongodbatlas.Auditing{
		Enabled: &enabled,
	}

	if currentModel.AuditAuthorizationSuccess != nil {
		auditingInput.AuditAuthorizationSuccess = currentModel.AuditAuthorizationSuccess
	}

	if currentModel.AuditFilter != nil {
		auditingInput.AuditFilter = *currentModel.AuditFilter
	}

	atlasAuditing, res, err := client.Auditing.Configure(context.Background(), *currentModel.GroupId, &auditingInput)

	if err != nil {
		log.Debugf("Create - error: %+v", err)
		return progress_events.GetFailedEventByResponse(err.Error(), res.Response), nil
	}

	currentModel.ConfigurationType = &atlasAuditing.ConfigurationType

	// Response
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		ResourceModel:   currentModel,
	}, nil
}

func Read(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup()

	log.Infof("Create() currentModel:%+v", *currentModel.GroupId)

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
	var res *mongodbatlas.Response

	atlasAuditing, res, err := client.Auditing.Get(context.Background(), *currentModel.GroupId)
	if err != nil {
		log.Debugf("Create - error: %+v", err)
		return progress_events.GetFailedEventByResponse(err.Error(), res.Response), nil
	}

	if *atlasAuditing.Enabled == false {
		return handler.ProgressEvent{
			HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound,
			OperationStatus:  handler.Failed,
		}, nil
	}

	currentModel.ConfigurationType = &atlasAuditing.ConfigurationType

	// Response
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "get successful",
		ResourceModel:   currentModel,
	}, nil
}

func Update(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup()

	log.Infof("Create() currentModel:%+v", *currentModel.GroupId)

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

	resourceEnabled, handlerEvent := isEnabled(*client, *currentModel)
	if handlerEvent != nil {
		return *handlerEvent, nil
	}
	if !resourceEnabled {
		return handler.ProgressEvent{
			HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound,
			OperationStatus:  handler.Failed,
			Message:          "resource not found",
		}, nil
	}

	var res *mongodbatlas.Response

	auditingInput := mongodbatlas.Auditing{}

	modified := false

	if currentModel.AuditAuthorizationSuccess != nil {
		modified = true
		auditingInput.AuditAuthorizationSuccess = currentModel.AuditAuthorizationSuccess
	}

	if currentModel.AuditFilter != nil {
		modified = true
		auditingInput.AuditFilter = *currentModel.AuditFilter
	}

	if !modified {
		return handler.ProgressEvent{
			OperationStatus: handler.Success,
			Message:         "Update success (no properties were changed)",
			ResourceModel:   currentModel,
		}, nil
	}

	atlasAuditing, res, err := client.Auditing.Configure(context.Background(), *currentModel.GroupId, &auditingInput)

	if err != nil {
		log.Debugf("Create - error: %+v", err)
		return progress_events.GetFailedEventByResponse(err.Error(), res.Response), nil
	}

	if err != nil {
		log.Debugf("Update - error: %+v", err)
		return progress_events.GetFailedEventByResponse(err.Error(), res.Response), nil
	}
	log.Debugf("Atlas Client %v", client)

	currentModel.ConfigurationType = &atlasAuditing.ConfigurationType

	// Response
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "Update success",
		ResourceModel:   currentModel,
	}, nil
}

func Delete(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup()

	log.Infof("Create() currentModel:%+v", *currentModel.GroupId)

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

	resourceEnabled, handlerEvent := isEnabled(*client, *currentModel)
	if handlerEvent != nil {
		return *handlerEvent, nil
	}

	if !resourceEnabled {
		return handler.ProgressEvent{
			HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound,
			OperationStatus:  handler.Failed,
		}, nil
	}

	var res *mongodbatlas.Response

	enabled := false

	auditingInput := mongodbatlas.Auditing{
		Enabled: &enabled,
	}

	_, res, err = client.Auditing.Configure(context.Background(), *currentModel.GroupId, &auditingInput)

	if err != nil {
		log.Debugf("Create - error: %+v", err)
		return progress_events.GetFailedEventByResponse(err.Error(), res.Response), nil
	}

	if err != nil {
		log.Debugf("Delete - error: %+v", err)
		return progress_events.GetFailedEventByResponse(err.Error(), res.Response), nil
	}
	log.Debugf("Atlas Client %v", client)

	// Response
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
	}, nil
}

func isEnabled(client mongodbatlas.Client, currentModel Model) (bool, *handler.ProgressEvent) {
	atlasAuditing, res, err := client.Auditing.Get(context.Background(), *currentModel.GroupId)

	if err != nil {
		log.Debugf("Create - error: %+v", err)
		er := progress_events.GetFailedEventByResponse(err.Error(), res.Response)
		return false, &er
	}

	return *atlasAuditing.Enabled, nil
}

func List(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	return handler.ProgressEvent{}, errors.New("Not implemented: List")
}
