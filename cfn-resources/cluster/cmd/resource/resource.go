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
	"reflect"
	"strings"

	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util"
	"github.com/spf13/cast"
	"go.mongodb.org/atlas/mongodbatlas"
)

const (
	tenant  = "TENANT"
	atlasM2 = "M2"
	atlasM5 = "M5"
	AWS     = "AWS"
)

var defaultLabel = mongodbatlas.Label{Key: "Infrastructure Tool", Value: "MongoDB Atlas Terraform Provider"}

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

// Read handles the Read event from the Cloudformation service.
func Read(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup()
	log.Debugf("Read() currentModel:%+v", currentModel)

	modelValidation := validateModel([]string{}, currentModel)
	if modelValidation != nil {
		return *modelValidation, nil
	}

	//Create Client
	client, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		return progress_events.GetFailedEventByCode(fmt.Sprintf("Error creating mongoDB client : %s", err.Error()),
			cloudformation.HandlerErrorCodeInvalidRequest), nil
	}

	//Read call
	model, res, err := readCluster(client, currentModel)
	if err != nil {
		log.Errorf("Cluster.Read() - error: %+v", err)
		return progress_events.GetFailedEventByResponse(fmt.Sprintf("Error in Read cluster : %s", err.Error()),
			res.Response), nil
	}
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "Read Complete",
		ResourceModel:   model,
	}, nil
}

// Update handles the Update event from the Cloudformation service.
func Update(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	return handler.ProgressEvent{
		OperationStatus:  handler.Failed,
		Message:          "Not implemented",
		HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound}, nil
}

// Delete handles the Delete event from the Cloudformation service.
func Delete(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	return handler.ProgressEvent{
		OperationStatus:  handler.Failed,
		Message:          "Not implemented",
		HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound}, nil
}

// List handles the List event from the Cloudformation service.
func List(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	return handler.ProgressEvent{
		OperationStatus:  handler.Failed,
		Message:          "Not implemented",
		HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound}, nil
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

func expandBiConnector(biConnector *BiConnector) *mongodbatlas.BiConnector {
	return &mongodbatlas.BiConnector{
		Enabled:        biConnector.Enabled,
		ReadPreference: cast.ToString(biConnector.ReadPreference),
	}
}

func expandReplicationSpecs(replicationSpecs []AdvancedReplicationSpec) []*mongodbatlas.AdvancedReplicationSpec {
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
			RegionConfigs: expandRegionsConfig(s.AdvancedRegionConfigs),
		}

		rSpecs = append(rSpecs, rSpec)
	}
	spew.Dump(rSpecs)
	fmt.Printf("specs: len %d %+v", len(replicationSpecs), rSpecs)
	return rSpecs
}

func expandAutoScaling(scaling *AdvancedAutoScaling) *mongodbatlas.AdvancedAutoScaling {
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

func expandRegionsConfig(regionConfigs []AdvancedRegionConfig) []*mongodbatlas.AdvancedRegionConfig {
	var regionsConfigs []*mongodbatlas.AdvancedRegionConfig
	for _, regionCfg := range regionConfigs {
		regionsConfigs = append(regionsConfigs, expandRegionConfig(regionCfg))
	}
	return regionsConfigs
}

func expandRegionConfig(regionCfg AdvancedRegionConfig) *mongodbatlas.AdvancedRegionConfig {
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

func flattenAutoScaling(scaling *mongodbatlas.AdvancedAutoScaling) *AdvancedAutoScaling {
	log.Debugf("expandAutoScaling: %+v", scaling)

	if scaling == nil {
		return nil
	}
	return &AdvancedAutoScaling{
		DiskGB: &DiskGB{Enabled: scaling.DiskGB.Enabled},
		Compute: &Compute{
			Enabled:          scaling.Compute.Enabled,
			ScaleDownEnabled: scaling.Compute.ScaleDownEnabled,
			MinInstanceSize:  &scaling.Compute.MinInstanceSize,
			MaxInstanceSize:  &scaling.Compute.MaxInstanceSize,
		},
	}
}

func flattenReplicationSpecs(replicationSpecs []*mongodbatlas.AdvancedReplicationSpec) []AdvancedReplicationSpec {
	var rSpecs []AdvancedReplicationSpec

	for ind, _ := range replicationSpecs {
		rSpec := AdvancedReplicationSpec{
			ID:                    &replicationSpecs[ind].ID,
			NumShards:             &replicationSpecs[ind].NumShards,
			ZoneName:              &replicationSpecs[ind].ZoneName,
			AdvancedRegionConfigs: flattenRegionsConfig(replicationSpecs[ind].RegionConfigs),
		}
		rSpecs = append(rSpecs, rSpec)
	}
	fmt.Printf("specs: len %d %+v", len(replicationSpecs), rSpecs)
	return rSpecs
}

func flattenRegionsConfig(regionConfigs []*mongodbatlas.AdvancedRegionConfig) []AdvancedRegionConfig {
	var regionsConfigs []AdvancedRegionConfig
	for _, regionCfg := range regionConfigs {
		regionsConfigs = append(regionsConfigs, flattenRegionConfig(regionCfg))
	}
	return regionsConfigs
}

func flattenRegionConfig(regionCfg *mongodbatlas.AdvancedRegionConfig) AdvancedRegionConfig {
	log.Debugf("expandRegionConfig: %+v", regionCfg)
	var region string

	return AdvancedRegionConfig{
		AutoScaling:    flattenAutoScaling(regionCfg.AutoScaling),
		RegionName:     &region,
		Priority:       regionCfg.Priority,
		AnalyticsSpecs: flattenRegionConfigSpec(regionCfg.AnalyticsSpecs),
		ElectableSpecs: flattenRegionConfigSpec(regionCfg.ElectableSpecs),
		ReadOnlySpecs:  flattenRegionConfigSpec(regionCfg.ReadOnlySpecs),
	}
}

func flattenRegionConfigSpec(spec *mongodbatlas.Specs) *Specs {
	log.Debugf("expandRegionConfigSpec: %+v", spec)
	if spec == nil {
		return nil
	}
	var ebsVolumeType string
	var instanceSize string
	return &Specs{
		DiskIOPS:      castNO64(spec.DiskIOPS),
		EbsVolumeType: &ebsVolumeType,
		InstanceSize:  &instanceSize,
		NodeCount:     spec.NodeCount,
	}
}

func flattenBiConnectorConfig(biConnector *mongodbatlas.BiConnector) *BiConnector {
	if biConnector == nil {
		return nil
	}

	return &BiConnector{
		ReadPreference: &biConnector.ReadPreference,
		Enabled:        biConnector.Enabled,
	}
}

func flattenConnectionStrings(clusterConnStrings *mongodbatlas.ConnectionStrings) *ConnectionStrings {
	var connStrings ConnectionStrings
	if clusterConnStrings != nil {
		connStrings = ConnectionStrings{
			Standard:        &clusterConnStrings.Standard,
			StandardSrv:     &clusterConnStrings.StandardSrv,
			Private:         &clusterConnStrings.Private,
			PrivateSrv:      &clusterConnStrings.PrivateSrv,
			PrivateEndpoint: flattenPrivateEndpoint(clusterConnStrings.PrivateEndpoint),
		}
	}
	return &connStrings
}

func flattenPrivateEndpoint(pes []mongodbatlas.PrivateEndpoint) []PrivateEndpoint {
	var prvEndpoints []PrivateEndpoint
	if pes == nil {
		return prvEndpoints
	}
	for ind, ePoint := range pes {
		pe := PrivateEndpoint{
			ConnectionString:    &pes[ind].ConnectionString,
			SRVConnectionString: &pes[ind].SRVConnectionString,
			Type:                &pes[ind].Type,
			Endpoints:           flattenEndpoints(ePoint.Endpoints),
		}
		prvEndpoints = append(prvEndpoints, pe)
	}
	return prvEndpoints
}

func flattenEndpoints(eps []mongodbatlas.Endpoint) []Endpoint {
	var endPoints []Endpoint
	for ind, _ := range eps {
		ep := Endpoint{
			EndpointID:   &eps[ind].EndpointID,
			ProviderName: &eps[ind].ProviderName,
			Region:       &eps[ind].Region,
		}
		endPoints = append(endPoints, ep)
	}
	return endPoints
}

func flattenLabels(clusterLabels []mongodbatlas.Label) []Labels {
	labels := make([]Labels, len(clusterLabels))
	for i, _ := range clusterLabels {
		labels[i] = Labels{
			Key:   &clusterLabels[i].Key,
			Value: &clusterLabels[i].Value,
		}
	}
	return labels
}

func formatMongoDBMajorVersion(val interface{}) string {
	if strings.Contains(val.(string), ".") {
		return val.(string)
	}
	return fmt.Sprintf("%.1f", cast.ToFloat32(val))
}

func isClusterInTargetState(client *mongodbatlas.Client, projectID, clusterName, targetState string) (bool, string, *mongodbatlas.AdvancedCluster, error) {
	cluster, resp, err := client.AdvancedClusters.Get(context.Background(), projectID, clusterName)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return "DELETED" == targetState, "DELETED", nil, nil
		}
		return false, "ERROR", nil, fmt.Errorf("error fetching cluster info (%s): %s", clusterName, err)
	}
	return cluster.StateName == targetState, cluster.StateName, cluster, nil
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

	return &args
}

func readCluster(client *mongodbatlas.Client, currentModel *Model) (*Model, *mongodbatlas.Response, error) {

	cluster, res, err := client.AdvancedClusters.Get(context.Background(), *currentModel.ProjectId, *currentModel.Name)

	currentModel.Id = &cluster.ID
	currentModel.BackupEnabled = cluster.BackupEnabled
	currentModel.BiConnector = flattenBiConnectorConfig(cluster.BiConnector)
	currentModel.ConnectionStrings = flattenConnectionStrings(cluster.ConnectionStrings)
	currentModel.ClusterType = &cluster.ClusterType
	currentModel.CreatedDate = &cluster.CreateDate

	currentModel.DiskSizeGB = cluster.DiskSizeGB
	currentModel.EncryptionAtRestProvider = &cluster.EncryptionAtRestProvider
	currentModel.Labels = flattenLabels(removeLabel(cluster.Labels, defaultLabel))
	currentModel.MongoDBMajorVersion = &cluster.MongoDBMajorVersion
	currentModel.MongoDBVersion = &cluster.MongoDBVersion

	currentModel.Paused = cluster.Paused
	currentModel.PitEnabled = cluster.PitEnabled
	currentModel.RootCertType = &cluster.RootCertType
	currentModel.ReplicationSpecs = flattenReplicationSpecs(cluster.ReplicationSpecs)

	currentModel.StateName = &cluster.StateName
	currentModel.VersionReleaseSystem = &cluster.VersionReleaseSystem
	return currentModel, res, err
}

func removeLabel(list []mongodbatlas.Label, item mongodbatlas.Label) []mongodbatlas.Label {
	var pos int
	for _, v := range list {
		if reflect.DeepEqual(v, item) {
			list = append(list[:pos], list[pos+1:]...)
			if pos > 0 {
				pos--
			}
			continue
		}
		pos++
	}
	return list
}

func updateAdvancedCluster(ctx context.Context, conn *mongodbatlas.Client, request *mongodbatlas.AdvancedCluster, projectID, name string) (*mongodbatlas.AdvancedCluster, *mongodbatlas.Response, error) {
	cluster, resp, err := conn.AdvancedClusters.Update(ctx, projectID, name, request)
	if err != nil {
		return nil, nil, err
	}

	return cluster, resp, nil
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
