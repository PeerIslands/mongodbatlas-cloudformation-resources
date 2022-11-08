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

var CreateRequiredFields = []string{"GroupId", "ClusterName", "CollectionName", "Database", "Name", "CollectionName", "Database", "Name", "ApiKeys.PrivateKey", "ApiKeys.PublicKey"}
var ReadRequiredFields = []string{"GroupId", "ClusterName", "IndexId", "ApiKeys.PrivateKey", "ApiKeys.PublicKey"}
var UpdateRequiredFields = []string{"GroupId", "ClusterName", "IndexId", "CollectionName", "Database", "Name", "ApiKeys.PrivateKey", "ApiKeys.PublicKey"}
var DeleteRequiredFields = []string{"GroupId", "ClusterName", "IndexId", "ApiKeys.PrivateKey", "ApiKeys.PublicKey"}
var ListRequiredFields = []string{"ApiKeys.PrivateKey", "ApiKeys.PublicKey"}

func validateModel(fields []string, model *Model) *handler.ProgressEvent {
	return validator.ValidateModel(fields, model)
}

func setup() {
	util.SetupLogger("mongodb-atlas-SearchIndex")
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
	   Analyzer, Analyzers, DatabaseName, Status, Database, ApiKeys, SearchAnalyzer, CollectionName, ClusterName, Name, Mappings, IndexID, GroupId, IndexId, Synonyms, ...
	*/
	/*
	    // Pseudocode:
	    res , resModel, err := client.SearchIndex.Create(context.Background(),&mongodbatlas.SearchIndex{
	   	Analyzer:currentModel.Analyzer,
	   	Analyzers:currentModel.Analyzers,
	   	DatabaseName:currentModel.DatabaseName,
	   	Status:currentModel.Status,
	   	Database:currentModel.Database,
	   	ApiKeys:currentModel.ApiKeys,
	   	SearchAnalyzer:currentModel.SearchAnalyzer,
	   	CollectionName:currentModel.CollectionName,
	   	ClusterName:currentModel.ClusterName,
	   	Name:currentModel.Name,
	   	Mappings:currentModel.Mappings,
	   	IndexID:currentModel.IndexID,
	   	GroupId:currentModel.GroupId,
	   	IndexId:currentModel.IndexId,
	   	Synonyms:currentModel.Synonyms,
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
	   GroupId, ClusterName, IndexId, ...
	*/
	/*
	    // Pseudocode:
	    res , resModel, err := client.SearchIndex.Read(context.Background(),&mongodbatlas.SearchIndex{
	   	GroupId:currentModel.GroupId,
	   	ClusterName:currentModel.ClusterName,
	   	IndexId:currentModel.IndexId,
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
	   ClusterName, Name, CollectionName, IndexID, GroupId, IndexId, Synonyms, Mappings, Analyzers, DatabaseName, Status, Analyzer, ApiKeys, SearchAnalyzer, Database, ...
	*/
	/*
	    // Pseudocode:
	    res , resModel, err := client.SearchIndex.Update(context.Background(),&mongodbatlas.SearchIndex{
	   	ClusterName:currentModel.ClusterName,
	   	Name:currentModel.Name,
	   	CollectionName:currentModel.CollectionName,
	   	IndexID:currentModel.IndexID,
	   	GroupId:currentModel.GroupId,
	   	IndexId:currentModel.IndexId,
	   	Synonyms:currentModel.Synonyms,
	   	Mappings:currentModel.Mappings,
	   	Analyzers:currentModel.Analyzers,
	   	DatabaseName:currentModel.DatabaseName,
	   	Status:currentModel.Status,
	   	Analyzer:currentModel.Analyzer,
	   	ApiKeys:currentModel.ApiKeys,
	   	SearchAnalyzer:currentModel.SearchAnalyzer,
	   	Database:currentModel.Database,
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
	   GroupId, ClusterName, IndexId, ...
	*/
	/*
	    // Pseudocode:
	    res , resModel, err := client.SearchIndex.Delete(context.Background(),&mongodbatlas.SearchIndex{
	   	GroupId:currentModel.GroupId,
	   	ClusterName:currentModel.ClusterName,
	   	IndexId:currentModel.IndexId,
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
	    res , resModel, err := client.SearchIndex.List(context.Background(),&mongodbatlas.SearchIndex{
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
