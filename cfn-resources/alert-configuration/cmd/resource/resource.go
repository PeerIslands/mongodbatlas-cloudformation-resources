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

var CreateRequiredFields = []string{"GroupId", "EventTypeName", "Links", "ApiKeys.PrivateKey", "ApiKeys.PublicKey"}
var ReadRequiredFields = []string{"GroupId", "AlertConfigId", "ApiKeys.PrivateKey", "ApiKeys.PublicKey"}
var UpdateRequiredFields = []string{"GroupId", "AlertConfigId", "EventTypeName", "Links", "ApiKeys.PrivateKey", "ApiKeys.PublicKey"}
var DeleteRequiredFields = []string{"GroupId", "AlertConfigId", "ApiKeys.PrivateKey", "ApiKeys.PublicKey"}
var ListRequiredFields = []string{"ApiKeys.PrivateKey", "ApiKeys.PublicKey"}

func validateModel(fields []string, model *Model) *handler.ProgressEvent {
	return validator.ValidateModel(fields, model)
}

func setup() {
	util.SetupLogger("mongodb-atlas-AlertConfiguration")
}

func Create(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup()

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
	var res *mongodbatlas.Response

	/*
	   Considerable params from currentModel:
	   Updated, EventTypeName, Results, ItemsPerPage, Id, Links, Threshold, Enabled, Matchers, TypeName, IncludeCount, PageNum, MetricThreshold, AlertConfigId, Created, GroupId, Notifications, ApiKeys, TotalCount, ...
	*/
	/*
	    // Pseudocode:
	    res , resModel, err := client.AlertConfiguration.Create(context.Background(),&mongodbatlas.AlertConfiguration{
	   	Updated:currentModel.Updated,
	   	EventTypeName:currentModel.EventTypeName,
	   	Results:currentModel.Results,
	   	ItemsPerPage:currentModel.ItemsPerPage,
	   	Id:currentModel.Id,
	   	Links:currentModel.Links,
	   	Threshold:currentModel.Threshold,
	   	Enabled:currentModel.Enabled,
	   	Matchers:currentModel.Matchers,
	   	TypeName:currentModel.TypeName,
	   	IncludeCount:currentModel.IncludeCount,
	   	PageNum:currentModel.PageNum,
	   	MetricThreshold:currentModel.MetricThreshold,
	   	AlertConfigId:currentModel.AlertConfigId,
	   	Created:currentModel.Created,
	   	GroupId:currentModel.GroupId,
	   	Notifications:currentModel.Notifications,
	   	ApiKeys:currentModel.ApiKeys,
	   	TotalCount:currentModel.TotalCount,
	   })

	*/

	if err != nil {
		log.Debugf("Create - error: %+v", err)
		return progress_events.GetFailedEventByResponse(err.Error(), res.Response), nil
	}
	log.Debugf("Atlas Client %v", client)

	// Response
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		ResourceModel:   currentModel,
	}, nil
}

func Read(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup()

	log.Debugf("Read() currentModel:%+v", currentModel)

	// Validation
	modelValidation := validateModel(ReadRequiredFields, currentModel)
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
	var res *mongodbatlas.Response

	/*
	   Considerable params from currentModel:
	   GroupId, AlertConfigId, ...
	*/
	/*
	    // Pseudocode:
	    res , resModel, err := client.AlertConfiguration.Read(context.Background(),&mongodbatlas.AlertConfiguration{
	   	GroupId:currentModel.GroupId,
	   	AlertConfigId:currentModel.AlertConfigId,
	   })

	*/

	if err != nil {
		log.Debugf("Read - error: %+v", err)
		return progress_events.GetFailedEventByResponse(err.Error(), res.Response), nil
	}
	log.Debugf("Atlas Client %v", client)

	// Response
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		ResourceModel:   currentModel,
	}, nil
}

func Update(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup()

	log.Debugf("Update() currentModel:%+v", currentModel)

	// Validation
	modelValidation := validateModel(UpdateRequiredFields, currentModel)
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
	var res *mongodbatlas.Response

	/*
	   Considerable params from currentModel:
	   Created, GroupId, Notifications, ApiKeys, TotalCount, Updated, EventTypeName, Results, Id, Links, Threshold, Enabled, Matchers, TypeName, IncludeCount, ItemsPerPage, PageNum, MetricThreshold, AlertConfigId, ...
	*/
	/*
	    // Pseudocode:
	    res , resModel, err := client.AlertConfiguration.Update(context.Background(),&mongodbatlas.AlertConfiguration{
	   	Created:currentModel.Created,
	   	GroupId:currentModel.GroupId,
	   	Notifications:currentModel.Notifications,
	   	ApiKeys:currentModel.ApiKeys,
	   	TotalCount:currentModel.TotalCount,
	   	Updated:currentModel.Updated,
	   	EventTypeName:currentModel.EventTypeName,
	   	Results:currentModel.Results,
	   	Id:currentModel.Id,
	   	Links:currentModel.Links,
	   	Threshold:currentModel.Threshold,
	   	Enabled:currentModel.Enabled,
	   	Matchers:currentModel.Matchers,
	   	TypeName:currentModel.TypeName,
	   	IncludeCount:currentModel.IncludeCount,
	   	ItemsPerPage:currentModel.ItemsPerPage,
	   	PageNum:currentModel.PageNum,
	   	MetricThreshold:currentModel.MetricThreshold,
	   	AlertConfigId:currentModel.AlertConfigId,
	   })

	*/

	if err != nil {
		log.Debugf("Update - error: %+v", err)
		return progress_events.GetFailedEventByResponse(err.Error(), res.Response), nil
	}
	log.Debugf("Atlas Client %v", client)

	// Response
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		ResourceModel:   currentModel,
	}, nil
}

func Delete(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup()

	log.Debugf("Delete() currentModel:%+v", currentModel)

	// Validation
	modelValidation := validateModel(DeleteRequiredFields, currentModel)
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
	var res *mongodbatlas.Response

	/*
	   Considerable params from currentModel:
	   GroupId, AlertConfigId, ...
	*/
	/*
	    // Pseudocode:
	    res , resModel, err := client.AlertConfiguration.Delete(context.Background(),&mongodbatlas.AlertConfiguration{
	   	GroupId:currentModel.GroupId,
	   	AlertConfigId:currentModel.AlertConfigId,
	   })

	*/

	if err != nil {
		log.Debugf("Delete - error: %+v", err)
		return progress_events.GetFailedEventByResponse(err.Error(), res.Response), nil
	}
	log.Debugf("Atlas Client %v", client)

	// Response
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		ResourceModel:   currentModel,
	}, nil
}

func List(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup()

	log.Debugf("List() currentModel:%+v", currentModel)

	// Validation
	modelValidation := validateModel(ListRequiredFields, currentModel)
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
	var res *mongodbatlas.Response

	//
	/*
	    // Pseudocode:
	    res , resModel, err := client.AlertConfiguration.List(context.Background(),&mongodbatlas.AlertConfiguration{
	   })

	*/

	if err != nil {
		log.Debugf("List - error: %+v", err)
		return progress_events.GetFailedEventByResponse(err.Error(), res.Response), nil
	}
	log.Debugf("Atlas Client %v", client)

	// Response
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		ResourceModel:   currentModel,
	}, nil
}
