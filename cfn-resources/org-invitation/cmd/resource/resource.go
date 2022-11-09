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

var CreateRequiredFields = []string{"OrgId", "ApiKeys.PrivateKey", "ApiKeys.PublicKey"}
var ReadRequiredFields = []string{"OrgId", "InvitationId", "ApiKeys.PrivateKey", "ApiKeys.PublicKey"}
var UpdateRequiredFields = []string{"OrgId", "InvitationId", "ApiKeys.PrivateKey", "ApiKeys.PublicKey"}
var DeleteRequiredFields = []string{"OrgId", "InvitationId", "ApiKeys.PrivateKey", "ApiKeys.PublicKey"}
var ListRequiredFields = []string{"ApiKeys.PrivateKey", "ApiKeys.PublicKey"}

func validateModel(fields []string, model *Model) *handler.ProgressEvent {
	return validator.ValidateModel(fields, model)
}

func setup() {
	util.SetupLogger("mongodb-atlas-OrgInvitation")
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
	   TotalCount, ItemsPerPage, IncludeCount, PageNum, TeamIds, OrgId, Links, ApiKeys, InviterUsername, Username, ExpiresAt, Id, OrgName, CreatedAt, Results, InvitationId, Roles, ...
	*/
	/*
	    // Pseudocode:
	    res , resModel, err := client.OrgInvitation.Create(context.Background(),&mongodbatlas.OrgInvitation{
	   	TotalCount:currentModel.TotalCount,
	   	ItemsPerPage:currentModel.ItemsPerPage,
	   	IncludeCount:currentModel.IncludeCount,
	   	PageNum:currentModel.PageNum,
	   	TeamIds:currentModel.TeamIds,
	   	OrgId:currentModel.OrgId,
	   	Links:currentModel.Links,
	   	ApiKeys:currentModel.ApiKeys,
	   	InviterUsername:currentModel.InviterUsername,
	   	Username:currentModel.Username,
	   	ExpiresAt:currentModel.ExpiresAt,
	   	Id:currentModel.Id,
	   	OrgName:currentModel.OrgName,
	   	CreatedAt:currentModel.CreatedAt,
	   	Results:currentModel.Results,
	   	InvitationId:currentModel.InvitationId,
	   	Roles:currentModel.Roles,
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
	   OrgId, InvitationId, ...
	*/
	/*
	    // Pseudocode:
	    res , resModel, err := client.OrgInvitation.Read(context.Background(),&mongodbatlas.OrgInvitation{
	   	OrgId:currentModel.OrgId,
	   	InvitationId:currentModel.InvitationId,
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
	   ApiKeys, OrgId, Links, Username, InviterUsername, OrgName, CreatedAt, Results, InvitationId, Roles, ExpiresAt, Id, IncludeCount, PageNum, TeamIds, TotalCount, ItemsPerPage, ...
	*/
	/*
	    // Pseudocode:
	    res , resModel, err := client.OrgInvitation.Update(context.Background(),&mongodbatlas.OrgInvitation{
	   	ApiKeys:currentModel.ApiKeys,
	   	OrgId:currentModel.OrgId,
	   	Links:currentModel.Links,
	   	Username:currentModel.Username,
	   	InviterUsername:currentModel.InviterUsername,
	   	OrgName:currentModel.OrgName,
	   	CreatedAt:currentModel.CreatedAt,
	   	Results:currentModel.Results,
	   	InvitationId:currentModel.InvitationId,
	   	Roles:currentModel.Roles,
	   	ExpiresAt:currentModel.ExpiresAt,
	   	Id:currentModel.Id,
	   	IncludeCount:currentModel.IncludeCount,
	   	PageNum:currentModel.PageNum,
	   	TeamIds:currentModel.TeamIds,
	   	TotalCount:currentModel.TotalCount,
	   	ItemsPerPage:currentModel.ItemsPerPage,
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
	   OrgId, InvitationId, ...
	*/
	/*
	    // Pseudocode:
	    res , resModel, err := client.OrgInvitation.Delete(context.Background(),&mongodbatlas.OrgInvitation{
	   	OrgId:currentModel.OrgId,
	   	InvitationId:currentModel.InvitationId,
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
	    res , resModel, err := client.OrgInvitation.List(context.Background(),&mongodbatlas.OrgInvitation{
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
