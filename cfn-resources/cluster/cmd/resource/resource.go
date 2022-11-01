package resource

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	progress_events "github.com/mongodb/mongodbatlas-cloudformation-resources/util/progress_event"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util/validator"
	log "github.com/sirupsen/logrus"
	"strings"

	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util"
	"github.com/spf13/cast"
	"go.mongodb.org/atlas/mongodbatlas"
)

func setup() {
	util.SetupLogger("mongodb-atlas-cluster")

}

func castNO64(i *int64) *int {
	x := cast.ToInt(&i)
	return &x
}
func cast64(i *int) *int64 {
	x := cast.ToInt64(&i)
	return &x
}

// validateModel inputs based on the method
func validateModel(fields []string, model *Model) *handler.ProgressEvent {
	return validator.ValidateModel(fields, model)
}

// Create handles the Create event from the Cloudformation service.
func Create(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup()
	log.Debugf("Create() currentModel:%+v", currentModel)
	spew.Dump("Model : ", currentModel)
	modelValidation := validateModel([]string{}, currentModel)
	if modelValidation != nil {
		return *modelValidation, nil
	}

	client, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		return progress_events.GetFailedEventByCode(fmt.Sprintf("Error creating mongoDB client : %s", err.Error()),
			cloudformation.HandlerErrorCodeInvalidRequest), nil
	}

	projectID := *currentModel.ProjectId
	log.Infof("cluster Create projectID=%s", projectID)
	if len(currentModel.ReplicationSpecs) > 0 {
		if currentModel.ClusterType == nil {
			err := errors.New("error creating cluster: ClusterType should be set when `ReplicationSpecs` is set")
			log.Infof("Create - error: %+v", err)
			return handler.ProgressEvent{
				OperationStatus:  handler.Failed,
				Message:          err.Error(),
				HandlerErrorCode: cloudformation.HandlerErrorCodeInvalidRequest}, nil
		}

		if currentModel.NumShards == nil {
			err := errors.New("error creating cluster: NumShards should be set when `ReplicationSpecs` is set")
			log.Infof("Create - error: %+v", err)
			return handler.ProgressEvent{
				OperationStatus:  handler.Failed,
				Message:          err.Error(),
				HandlerErrorCode: cloudformation.HandlerErrorCodeInvalidRequest}, nil
		}
	}

	clusterRequest := &mongodbatlas.AdvancedCluster{
		Name:                     cast.ToString(currentModel.Name),
		EncryptionAtRestProvider: cast.ToString(currentModel.EncryptionAtRestProvider),
		ClusterType:              cast.ToString(currentModel.ClusterType),
		ReplicationSpecs:         expandReplicationSpecs(currentModel.ReplicationSpecs),
	}

	if currentModel.BackupEnabled != nil {
		clusterRequest.BackupEnabled = currentModel.BackupEnabled
	}

	if currentModel.DiskSizeGB != nil {
		currentModel.DiskSizeGB = clusterRequest.DiskSizeGB
	}

	if currentModel.MongoDBMajorVersion != nil {
		clusterRequest.MongoDBMajorVersion = formatMongoDBMajorVersion(*currentModel.MongoDBMajorVersion)
	}

	if currentModel.BiConnector != nil {
		clusterRequest.BiConnector = expandBiConnector(currentModel.BiConnector)
	}

	if currentModel.PitEnabled != nil {
		clusterRequest.PitEnabled = currentModel.PitEnabled
	}

	if currentModel.VersionReleaseSystem != nil {
		clusterRequest.VersionReleaseSystem = *currentModel.VersionReleaseSystem
	}

	if currentModel.RootCertType != nil {
		clusterRequest.RootCertType = *currentModel.RootCertType
	}
	resJson, err := json.Marshal(clusterRequest)
	if err == nil {
		fmt.Println("input:", string(resJson))
	}
	jsonStr, _ := json.Marshal(clusterRequest)
	fmt.Println(string(jsonStr))
	log.Printf("clusterRequest --- value:%s ", jsonStr)
	ctx := context.Background()
	log.Debugf("DEBUG: clusterRequest: %+v", clusterRequest)

	//Create Callback
	if _, ok := req.CallbackContext["stateName"]; ok {
		return callBackClusterCreate(req, client, currentModel, ctx, projectID)
	}

	//Create Cluster
	cluster, res, err := client.AdvancedClusters.Create(ctx, projectID, clusterRequest)
	if err != nil {
		log.Errorf("Create - Cluster.Create() - error: %+v", err)
		return progress_events.GetFailedEventByResponse(fmt.Sprintf("Error creating resource : %s", err.Error()),
			res.Response), nil
	}

	currentModel.StateName = &cluster.StateName
	event := handler.ProgressEvent{
		OperationStatus:      handler.InProgress,
		Message:              fmt.Sprintf("Create Cluster `%s`", cluster.StateName),
		ResourceModel:        currentModel,
		CallbackDelaySeconds: 65,
		CallbackContext: map[string]interface{}{
			"stateName":   cluster.StateName,
			"projectId":   projectID,
			"clusterName": *currentModel.Name,
		},
	}
	log.Debugf("Create() return event:%+v", event)
	return event, nil
}

func callBackClusterCreate(req handler.Request, client *mongodbatlas.Client, currentModel *Model, ctx context.Context, projectID string) (handler.ProgressEvent, error) {
	progressEvent, err := validateProgress(client, req, currentModel, "IDLE", "CREATING")
	if err != nil {
		return progressEvent, nil
	}
	if progressEvent.Message == "Complete" {
		log.Debugf("Compelted creation:%s", *currentModel.Name)

		cluster, res, err := client.AdvancedClusters.Get(ctx, projectID, *currentModel.Name)
		if err != nil {
			return progress_events.GetFailedEventByResponse(fmt.Sprintf("Error creating resource : %s", err.Error()),
				res.Response), nil
		}

		log.Debugf("Updating cluster :%s", *currentModel.Name)
		return updateCluster(currentModel, client, ctx, projectID, cluster)

	}
	return progressEvent, nil
}
func callBackClusterPause(req handler.Request, client *mongodbatlas.Client, currentModel *Model, ctx context.Context, projectID string) (handler.ProgressEvent, error) {
	progressEvent, err := validateProgress(client, req, currentModel, "IDLE", "CREATING")
	if err != nil {
		return progressEvent, nil
	}
	if progressEvent.Message == "Complete" {
		log.Debugf("Compelted pause:%s", *currentModel.Name)

		cluster, res, err := client.AdvancedClusters.Get(ctx, projectID, *currentModel.Name)
		if err != nil {
			return progress_events.GetFailedEventByResponse(fmt.Sprintf("Error creating resource : %s", err.Error()),
				res.Response), nil
		}

		progressEvent = handler.ProgressEvent{
			OperationStatus: handler.Success,
			Message:         fmt.Sprintf("Pause Cluster Status `%s`", cluster.StateName),
			ResourceModel:   currentModel,
		}
	}
	return progressEvent, nil
}

func updateCluster(currentModel *Model, client *mongodbatlas.Client, ctx context.Context, projectID string, cluster *mongodbatlas.AdvancedCluster) (handler.ProgressEvent, error) {
	//Update advanced configuration
	if currentModel.AdvancedSettings != nil {
		advancedConfig := processAdvancedSettings(*currentModel.AdvancedSettings)
		_, res, err := client.Clusters.UpdateProcessArgs(ctx, projectID, cluster.Name, advancedConfig)
		if err != nil {
			log.Errorf("Cluster UpdateProcessArgs - error: %+v", err)
			return progress_events.GetFailedEventByResponse(fmt.Sprintf("Error creating resource : %s", err.Error()),
				res.Response), err
		}
	}

	//Update pause
	if (currentModel.Paused != nil) && (*currentModel.Paused != *cluster.Paused) {
		_, res, err := updateAdvancedCluster(ctx, client, &mongodbatlas.AdvancedCluster{Paused: currentModel.Paused}, projectID, *currentModel.Name)
		if err != nil {
			log.Errorf("Cluster Pause - error: %+v", err)
			return progress_events.GetFailedEventByResponse(fmt.Sprintf("Error creating resource : %s", err.Error()),
				res.Response), err
		}
	}
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         fmt.Sprintf("Create Cluster `%s`", cluster.StateName),
		ResourceModel:   currentModel,
	}, nil
}

func Read(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	return handler.ProgressEvent{
		OperationStatus:  handler.Failed,
		Message:          "Not implemented",
		HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound}, nil
}
func Update(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	return handler.ProgressEvent{
		OperationStatus:  handler.Failed,
		Message:          "Not implemented",
		HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound}, nil
}
func Delete(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	return handler.ProgressEvent{
		OperationStatus:  handler.Failed,
		Message:          "Not implemented",
		HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound}, nil
}
func List(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	return handler.ProgressEvent{
		OperationStatus:  handler.Failed,
		Message:          "Not implemented",
		HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound}, nil
}

func expandBiConnector(biConnector *BiConnector) *mongodbatlas.BiConnector {
	return &mongodbatlas.BiConnector{
		Enabled:        biConnector.Enabled,
		ReadPreference: cast.ToString(biConnector.ReadPreference),
	}
}

const (
	tenant  = "TENANT"
	atlasM2 = "M2"
	atlasM5 = "M5"
	AWS     = "AWS"
)

func expandReplicationSpecs(replicationSpecs []ReplicationSpec) []*mongodbatlas.AdvancedReplicationSpec {
	var rSpecs []*mongodbatlas.AdvancedReplicationSpec

	for _, s := range replicationSpecs {
		var numShards int
		if s.NumShards != nil {
			numShards = *s.NumShards
		}
		rSpec := &mongodbatlas.AdvancedReplicationSpec{
			ID:            cast.ToString(s.ID),
			NumShards:     numShards,
			ZoneName:      cast.ToString(s.ZoneName),
			RegionConfigs: expandRegionsConfig(s.RegionsConfig),
		}

		rSpecs = append(rSpecs, rSpec)
	}
	spew.Dump(rSpecs)
	fmt.Printf("specs: len %d %+v", len(replicationSpecs), rSpecs)
	return rSpecs
}

func expandRegionsConfig(regionConfigs []RegionConfig) []*mongodbatlas.AdvancedRegionConfig {
	var regionsConfigs []*mongodbatlas.AdvancedRegionConfig
	for _, regionCfg := range regionConfigs {
		regionsConfigs = append(regionsConfigs, expandRegionConfig(regionCfg))
	}
	return regionsConfigs
}

func expandRegionConfig(regionCfg RegionConfig) *mongodbatlas.AdvancedRegionConfig {
	log.Debugf("expandRegionConfig: %+v", regionCfg)
	var region string
	if regionCfg.RegionName != nil {
		region = *regionCfg.RegionName
	}
	return &mongodbatlas.AdvancedRegionConfig{
		AutoScaling:    expandAutoScaling(regionCfg.AutoScaling),
		ProviderName:   AWS,
		RegionName:     region,
		Priority:       regionCfg.Priority,
		AnalyticsSpecs: expandRegionConfigSpec(regionCfg.AnalyticsSpecs),
		ElectableSpecs: expandRegionConfigSpec(regionCfg.ElectableSpecs),
		ReadOnlySpecs:  expandRegionConfigSpec(regionCfg.ReadOnlySpecs),
	}
}

func expandAutoScaling(scaling *AnalyticsAutoScaling) *mongodbatlas.AdvancedAutoScaling {
	log.Debugf("expandAutoScaling: %+v", scaling)

	var minInstanceSize string
	if scaling == nil {
		return nil
	}
	if scaling.Compute.MinInstanceSize != nil {
		minInstanceSize = *scaling.Compute.MinInstanceSize
	}
	var maxInstanceSize string
	if scaling.Compute.MaxInstanceSize != nil {
		maxInstanceSize = *scaling.Compute.MaxInstanceSize
	}
	return &mongodbatlas.AdvancedAutoScaling{
		DiskGB: &mongodbatlas.DiskGB{Enabled: scaling.DiskGB.Enabled},
		Compute: &mongodbatlas.Compute{
			Enabled:          scaling.Compute.Enabled,
			ScaleDownEnabled: scaling.Compute.ScaleDownEnabled,
			MinInstanceSize:  minInstanceSize,
			MaxInstanceSize:  maxInstanceSize,
		},
	}
}

func expandRegionConfigSpec(spec *Specs) *mongodbatlas.Specs {
	log.Debugf("expandRegionConfigSpec: %+v", spec)
	if spec == nil {
		return nil
	}
	var ebsVolumeType string
	var instanceSize string
	if spec.EbsVolumeType != nil {
		ebsVolumeType = *spec.EbsVolumeType
	}
	if spec.InstanceSize != nil {
		instanceSize = *spec.InstanceSize
	}
	return &mongodbatlas.Specs{
		DiskIOPS:      cast64(spec.DiskIOPS),
		EbsVolumeType: ebsVolumeType,
		InstanceSize:  instanceSize,
		NodeCount:     spec.NodeCount,
	}
}

func formatMongoDBMajorVersion(val interface{}) string {
	if strings.Contains(val.(string), ".") {
		return val.(string)
	}
	return fmt.Sprintf("%.1f", cast.ToFloat32(val))
}

func updateAdvancedCluster(ctx context.Context, conn *mongodbatlas.Client, request *mongodbatlas.AdvancedCluster, projectID, name string) (*mongodbatlas.AdvancedCluster, *mongodbatlas.Response, error) {
	cluster, resp, err := conn.AdvancedClusters.Update(ctx, projectID, name, request)
	if err != nil {
		return nil, nil, err
	}

	return cluster, resp, nil
}

func processAdvancedSettings(processArgs ProcessArgs) *mongodbatlas.ProcessArgs {
	var args mongodbatlas.ProcessArgs
	if processArgs.DefaultReadConcern != nil {
		args.DefaultReadConcern = *processArgs.DefaultReadConcern
	}
	args.FailIndexKeyTooLong = processArgs.FailIndexKeyTooLong
	if processArgs.DefaultWriteConcern != nil {
		args.DefaultWriteConcern = *processArgs.DefaultWriteConcern
	}
	args.JavascriptEnabled = processArgs.JavascriptEnabled
	if processArgs.MinimumEnabledTLSProtocol != nil {
		args.MinimumEnabledTLSProtocol = *processArgs.MinimumEnabledTLSProtocol
	}
	args.NoTableScan = processArgs.NoTableScan
	if processArgs.OplogSizeMB != nil {
		args.OplogSizeMB = cast64(processArgs.OplogSizeMB)
	}
	if processArgs.SampleSizeBIConnector != nil {
		args.SampleSizeBIConnector = cast64(processArgs.SampleSizeBIConnector)
	}
	if processArgs.SampleRefreshIntervalBIConnector != nil {
		args.SampleRefreshIntervalBIConnector = cast64(processArgs.SampleRefreshIntervalBIConnector)
	}
	//TODO: cross check with API
	//args.OplogMinRetentionHours = cast64(processArgs.)
	return &args
}

func validateProgress(client *mongodbatlas.Client, req handler.Request, currentModel *Model, targetState string, pendingState string) (handler.ProgressEvent, error) {
	log.Debugf(" Cluster validateProgress() currentModel:%+v", currentModel)
	isReady, state, cluster, err := isClusterInTargetState(client, *currentModel.ProjectId, *currentModel.Name, targetState)
	log.Debugf("Cluster validateProgress() isReady:%+v, state:+%v, cluster:%+v", isReady, state, cluster)
	if err != nil {
		log.Debugf("ERROR Cluster validateProgress() err:%+v", err)
		return handler.ProgressEvent{
			Message:          err.Error(),
			OperationStatus:  handler.Failed,
			HandlerErrorCode: cloudformation.HandlerErrorCodeServiceInternalError}, nil
	}

	if !isReady {
		p := handler.NewProgressEvent()
		p.ResourceModel = currentModel
		p.OperationStatus = handler.InProgress
		p.CallbackDelaySeconds = 60
		p.Message = "Pending"
		p.CallbackContext = map[string]interface{}{
			"stateName": state,
		}
		return p, nil
	}

	p := handler.NewProgressEvent()
	p.OperationStatus = handler.Success
	p.Message = "Complete"
	if targetState != "DELETED" {
		p.ResourceModel = currentModel
	}
	return p, nil
}

func isClusterInTargetState(client *mongodbatlas.Client, projectID, clusterName, targetState string) (bool, string, *mongodbatlas.Cluster, error) {
	cluster, resp, err := client.Clusters.Get(context.Background(), projectID, clusterName)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return "DELETED" == targetState, "DELETED", nil, nil
		}
		return false, "ERROR", nil, fmt.Errorf("error fetching cluster info (%s): %s", clusterName, err)
	}
	return cluster.StateName == targetState, cluster.StateName, cluster, nil
}
