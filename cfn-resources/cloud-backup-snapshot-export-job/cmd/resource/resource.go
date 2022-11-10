package resource

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util/logger"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util/progressevent"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util/validator"
	"github.com/openlyinc/pointy"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/atlas/mongodbatlas"
)

const (
	publicKey              = "ApiKeys.PublicKey"
	privateKey             = "ApiKeys.PrivateKey"
	projectId              = "GroupId"
	userName               = "UserName"
	exportBucketId         = "ExportBucketId"
	snapshotId             = "SnapshotId"
	exportId               = "ExportId"
	clusterName            = "ClusterName"
	errorCreateMongoClient = "Error - Create MongoDB Client- Details: %+v"
	errorExportJobCreate   = "error creating Export Job for the project(%s) : %s"
	errorExportJobRead     = "error reading export job for the projects(%s) : Job Id : %s with error :%+v"
	errorExportJobDelete   = "error deleting Export Job for the projects(%s)(%s): %s"
)

var CreateRequiredFields = []string{publicKey, privateKey, projectId, exportBucketId, snapshotId}
var ReadRequiredFields = []string{publicKey, privateKey, projectId, exportId, clusterName}
var DeleteRequiredFields = []string{publicKey, privateKey, projectId}
var ListRequiredFields = []string{publicKey, privateKey, projectId, userName}

func validateModel(fields []string, model *Model) *handler.ProgressEvent {
	return validator.ValidateModel(fields, model)
}

func setup() {
	util.SetupLogger("mongodb-atlas-CloudBackupSnapshotExportJob")
}

func Create(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup() // logger setup
	_, _ = logger.Debugf("Create snapshot for Request() currentModel:%+v", currentModel)

	// Validate required fields in the request
	if modelValidation := validateModel(CreateRequiredFields, currentModel); modelValidation != nil {
		return *modelValidation, nil
	}

	// Create MongoDb Atlas Client using keys
	client, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		_, _ = logger.Warnf(errorCreateMongoClient, err)
		return progressevents.GetFailedEventByCode(fmt.Sprintf("Failed to Create Client : %s", err.Error()),
			cloudformation.HandlerErrorCodeInvalidRequest), nil
	}
	log.Info("111111111111111111")
	projectID := *currentModel.GroupId
	clusterName := *currentModel.ClusterName
	log.Info("222222222222222222222222")
	customData := expandExportJobCustomData(currentModel)
	log.Info("3333333333333333333333")
	request := &mongodbatlas.CloudProviderSnapshotExportJob{
		SnapshotID:     *currentModel.SnapshotId,
		ExportBucketID: *currentModel.ExportBucketId,
		CustomData:     customData,
	}
	log.Info("444444444444444444444444444444")
	// progress callback setup
	if _, ok := req.CallbackContext["status"]; ok {
		sid := req.CallbackContext["export_id"].(string)
		currentModel.ExportId = &sid
		return validateProgress(client, currentModel, "Successful")
	}
	jobResponse, resp, err := client.CloudProviderSnapshotExportJobs.Create(context.Background(), projectID, clusterName, request)
	log.Info("555555555555555555555555555555555555")
	if err != nil {
		log.Infof(errorExportJobCreate, projectID, err)
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          err.Error(),
			HandlerErrorCode: cloudformation.HandlerErrorCodeInvalidRequest}, nil

	}
	log.Info("6666666666666666666666666666666666")
	log.Info(resp)
	log.Info(jobResponse)
	log.Info(jobResponse.ID)
	currentModel.ExportId = &jobResponse.ID
	log.Info("6666666666666666666666666666666666")
	log.Debugf("Atlas Client %v", &jobResponse.State)

	// track progress
	event := handler.ProgressEvent{
		OperationStatus:      handler.InProgress,
		Message:              fmt.Sprintf("Create export snapshots : %s", jobResponse.ID),
		ResourceModel:        currentModel,
		CallbackDelaySeconds: 65,
		CallbackContext: map[string]interface{}{
			"status":    jobResponse.State,
			"export_id": jobResponse.ID,
		},
	}
	log.Debugf("Create() return event:%+v", event)

	return event, nil
}

func Read(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup() // logger setup
	_, _ = logger.Debugf("Read snapshot for Request() currentModel:%+v", currentModel)

	// Validate required fields in the request
	if modelValidation := validateModel(ReadRequiredFields, currentModel); modelValidation != nil {
		return *modelValidation, nil
	}

	// Create MongoDb Atlas Client using keys
	client, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		_, _ = logger.Warnf(errorCreateMongoClient, err)
		return progressevents.GetFailedEventByCode(fmt.Sprintf("Failed to Create Client : %s", err.Error()),
			cloudformation.HandlerErrorCodeInvalidRequest), nil
	}

	projectID := *currentModel.GroupId
	clusterName := *currentModel.ClusterName
	exportJobID := *currentModel.ExportId

	if !isExist(client, projectID, clusterName, clusterName) {
		_, _ = logger.Warnf(errorExportJobRead, projectID, exportJobID, errors.New("resource Not Found"))
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          "Resource Not Found",
			HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound}, nil
	}

	var res mongodbatlas.Response
	log.Info("555555555555555555555555555555555555555")
	exportJob, resp, err := client.CloudProviderSnapshotExportJobs.Get(context.Background(), projectId, clusterName, exportJobID)
	if err != nil {
		log.Debugf(errorExportJobRead, projectID, exportJobID, err)
		return progressevents.GetFailedEventByResponse(err.Error(), res.Response), nil
	}
	log.Info("66666666666666666666666666666666666666666")
	currentModel.ExportId = &exportJob.ID
	currentModel.ExportBucketId = &exportJob.ExportBucketID
	currentModel.CreatedAt = &exportJob.CreatedAt
	currentModel.FinishedAt = &exportJob.FinishedAt
	currentModel.CreatedAt = &exportJob.CreatedAt
	currentModel.Prefix = &exportJob.Prefix
	currentModel.State = &exportJob.State
	currentModel.SnapshotId = &exportJob.SnapshotID
	currentModel.Links = flattenLinks(resp.Links)
	if exportJob.ExportStatus != nil {
		currentModel.ExportStatus = &ApiExportStatusView{
			ExportedCollections: pointy.Int(exportJob.ExportStatus.ExportedCollections),
			TotalCollections:    pointy.Int(exportJob.ExportStatus.TotalCollections),
		}

	}
	currentModel.ExportStatus = flattenStatus(exportJob.ExportStatus)
	currentModel.CustomDataSet = flattenExportJobsCustomData(exportJob.CustomData)
	currentModel.Components = flattenExportComponent(exportJob.Components)
	log.Info("77777777777777777777777777777777777777777777")

	log.Debugf("Read Result : %v", currentModel)

	// Response
	event := handler.ProgressEvent{
		OperationStatus: handler.Success,
		ResourceModel:   currentModel,
	}
	return event, nil
}
func flattenLinks(linksResult []*mongodbatlas.Link) []Link {
	if len(linksResult) == 0 {
		return nil
	}
	log.Info("KKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKkkkk")
	log.Info(len(linksResult))
	log.Info("sssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssss")

	links := make([]Link, 0)
	for _, link := range linksResult {
		var lin Link
		lin.Href = &link.Href
		lin.Rel = &link.Rel
		links = append(links, lin)
	}
	return links
}
func flattenStatus(v *mongodbatlas.CloudProviderSnapshotExportJobStatus) *ApiExportStatusView {
	log.Info("statusstatusstatusstatusstatusstatusstatusstatusstatusstatusstatus")
	log.Info(v)

	log.Info("statusstatusstatusstatusstatusstatucccccccccccccccccccsstatusstatusstatusstatusstatus")
	log.Info("statusstatusstatusstatusstatusstatusstatusstatusstatusstatusstatus")

	log.Info(pointy.Int(v.ExportedCollections))
	log.Info(pointy.Int(v.TotalCollections))
	log.Info("vbvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv")

	status := ApiExportStatusView{
		ExportedCollections: pointy.Int(v.ExportedCollections),
		TotalCollections:    pointy.Int(v.TotalCollections),
	}

	return &status
}

func flattenExportJobsCustomData(m []*mongodbatlas.CloudProviderSnapshotExportJobCustomData) []CustomData {

	statusList := make(
		[]CustomData,
		len(m),
	)

	for i, _ := range m {
		v := m[i]
		role := CustomData{
			Key:   pointy.String(v.Key),
			Value: pointy.String(v.Value),
		}

		statusList = append(statusList, role)
	}
	return statusList
}
func flattenExportComponent(m []*mongodbatlas.CloudProviderSnapshotExportJobComponent) []ApiAtlasDiskBackupBaseRestoreMemberView {

	statusList := make(
		[]ApiAtlasDiskBackupBaseRestoreMemberView,
		len(m),
	)

	for i, _ := range m {
		v := m[i]
		role := ApiAtlasDiskBackupBaseRestoreMemberView{
			ReplicaSetName: pointy.String(v.ReplicaSetName),
			ExportID:       pointy.String(v.ExportID),
		}

		statusList = append(statusList, role)
	}
	return statusList
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

	//
	/*
	    // Pseudocode:
	    res , resModel, err := client.cluster-export-job.Update(context.Background(),&mongodbatlas.Cluster-export-job{
	   })

	*/

	if err != nil {
		log.Debugf("Update - error: %+v", err)
		return progressevents.GetFailedEventByResponse(err.Error(), res.Response), nil
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

	//
	/*
	    // Pseudocode:
	    res , resModel, err := client.cluster-export-job.Delete(context.Background(),&mongodbatlas.Cluster-export-job{
	   })

	*/

	if err != nil {
		log.Debugf("Delete - error: %+v", err)
		return progressevents.GetFailedEventByResponse(err.Error(), res.Response), nil
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

	// Create Atlas API Request Object
	projectId := *currentModel.GroupId
	clusterName := *currentModel.ClusterName

	params := &mongodbatlas.ListOptions{
		PageNum:      0,
		ItemsPerPage: 100,
	}
	// API call
	exportJobs, _, err := client.CloudProviderSnapshotExportJobs.List(context.Background(), projectId, clusterName, params)

	if err != nil {
		return handler.ProgressEvent{}, fmt.Errorf("error reading cloud provider snapshot restore job list with id(project: %s): %s", projectId, err)
	}

	var models []Model
	for _, exportJob := range exportJobs.Results {
		var model Model
		model.ExportId = &exportJob.ID
		model.ExportBucketId = &exportJob.ExportBucketID
		model.CreatedAt = &exportJob.CreatedAt
		model.FinishedAt = &exportJob.FinishedAt
		model.CreatedAt = &exportJob.CreatedAt
		model.Prefix = &exportJob.Prefix
		currentModel.Links = flattenLinks(exportJobs.Links)
		model.State = &exportJob.State
		model.SnapshotId = &exportJob.SnapshotID
		if exportJob.ExportStatus != nil {
			model.ExportStatus = &ApiExportStatusView{
				ExportedCollections: pointy.Int(exportJob.ExportStatus.ExportedCollections),
				TotalCollections:    pointy.Int(exportJob.ExportStatus.TotalCollections),
			}

		}
		model.CustomDataSet = flattenExportJobsCustomData(exportJob.CustomData)

		log.Info("uuuuuuuuuuuuuuuuuuuuuu")
		model.Components = flattenExportComponent(exportJob.Components)
		log.Info("vvvvvvvvvvvvv")
		models = append(models, model)
	}
	log.Debug("List cloud backup restore job ends")
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "List complete",
		ResourceModel:   models,
	}, nil
}

// function to track snapshot creation status
func validateProgress(client *mongodbatlas.Client, currentModel *Model, targetState string) (handler.ProgressEvent, error) {
	exportId := *currentModel.ExportId
	projectId := *currentModel.GroupId
	clusterName := *currentModel.ClusterName
	isReady, state, err := isJobInTargetState(client, projectId, exportId, clusterName, targetState)
	if err != nil || state == "Cancelled" {
		return handler.ProgressEvent{}, err
	}

	if !isReady {
		p := handler.NewProgressEvent()
		p.ResourceModel = currentModel
		p.OperationStatus = handler.InProgress
		p.CallbackDelaySeconds = 35
		p.Message = "Pending"
		p.CallbackContext = map[string]interface{}{
			"status":    state,
			"export_id": *currentModel.ExportId,
		}
		return p, nil
	}

	exportJob, resp, err := client.CloudProviderSnapshotExportJobs.Get(context.Background(), projectId, clusterName, exportId)
	currentModel.ExportId = &exportJob.ID
	currentModel.ExportBucketId = &exportJob.ExportBucketID
	currentModel.CreatedAt = &exportJob.CreatedAt
	currentModel.FinishedAt = &exportJob.FinishedAt
	currentModel.CreatedAt = &exportJob.CreatedAt
	currentModel.Prefix = &exportJob.Prefix
	currentModel.State = &exportJob.State
	currentModel.SnapshotId = &exportJob.SnapshotID
	currentModel.Links = flattenLinks(resp.Links)
	if exportJob.ExportStatus != nil {
		currentModel.ExportStatus = &ApiExportStatusView{
			ExportedCollections: pointy.Int(exportJob.ExportStatus.ExportedCollections),
			TotalCollections:    pointy.Int(exportJob.ExportStatus.TotalCollections),
		}

	}
	currentModel.ExportStatus = flattenStatus(exportJob.ExportStatus)
	currentModel.CustomDataSet = flattenExportJobsCustomData(exportJob.CustomData)
	currentModel.Components = flattenExportComponent(exportJob.Components)
	p := handler.NewProgressEvent()
	p.ResourceModel = currentModel
	p.OperationStatus = handler.Success
	p.Message = "Complete"
	return p, nil
}

// function to check if export job is in target state
func isJobInTargetState(client *mongodbatlas.Client, projectId, exportJobID, clusterName, targetState string) (bool, string, error) {
	exportJob, resp, err := client.CloudProviderSnapshotExportJobs.Get(context.Background(), projectId, clusterName, exportJobID)
	if err != nil {
		if exportJob == nil && resp == nil {
			return false, "", err
		}

		return false, "", err
	}
	return exportJob.State == targetState, exportJob.State, nil
}

// function to check if snapshot already exist in atlas
func isExist(client *mongodbatlas.Client, projectId, exportJobID, clusterName string) bool {
	exportJob, _, err := client.CloudProviderSnapshotExportJobs.Get(context.Background(), projectId, clusterName, exportJobID)
	if err != nil {
		return false
	} else if exportJob == nil {
		return false
	}
	return true
}

// function to convert custom metadata from request to mongodbatlas object
func expandExportJobCustomData(currentModel *Model) []*mongodbatlas.CloudProviderSnapshotExportJobCustomData {
	customData := currentModel.CustomDataSet
	if customData != nil {
		res := make([]*mongodbatlas.CloudProviderSnapshotExportJobCustomData, len(customData))

		for i, val := range customData {
			res[i] = &mongodbatlas.CloudProviderSnapshotExportJobCustomData{
				Key:   *val.Key,
				Value: *val.Value,
			}
		}
		return res
	}
	return nil
}
