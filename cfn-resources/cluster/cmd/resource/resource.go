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
	"strconv"
	"strings"

	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util"
	"github.com/spf13/cast"
	"go.mongodb.org/atlas/mongodbatlas"
)

const (
	AWS        = "AWS"
	labelError = "you should not set `Infrastructure Tool` label, it is used for internal purposes"
	Id         = "Id"
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
	//TODO: add validation params
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

	// Callback
	if _, ok := req.CallbackContext["stateName"]; ok {
		if id, idExists := req.CallbackContext[Id]; idExists {
			idStr := fmt.Sprint(id)
			currentModel.Id = &idStr
			return createClusterCallback(req, client, currentModel, projectID)
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

	if currentModel.BiConnector != nil {
		clusterRequest.BiConnector = expandBiConnector(currentModel.BiConnector)
	}

	if currentModel.DiskSizeGB != nil {
		clusterRequest.DiskSizeGB = currentModel.DiskSizeGB
	}

	if len(currentModel.Labels) > 0 {
		clusterRequest.Labels = expandLabelSlice(currentModel.Labels)
		if containsLabelOrKey(clusterRequest.Labels, defaultLabel) {
			log.Infof("Create - error: %+v", err)
			return progress_events.GetFailedEventByCode(
				labelError,
				cloudformation.HandlerErrorCodeInvalidRequest), nil
		}
	}
	if currentModel.MongoDBMajorVersion != nil {
		clusterRequest.MongoDBMajorVersion = formatMongoDBMajorVersion(*currentModel.MongoDBMajorVersion)
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

	//TODO: advanced configuration oplog-size validation if exists

	//AWS
	// This is the intial call to Create, so inject a deployment
	// secret for this resource in order to lookup progress properly
	projectResID := &util.ResourceIdentifier{
		ResourceType: "Project",
		ResourceID:   projectID,
	}
	log.Debugf("Created projectResID:%s", projectResID)
	resourceID := util.NewResourceIdentifier("Cluster", *currentModel.Name, projectResID)
	log.Debugf("Created resourceID:%s", resourceID)
	resourceProps := map[string]string{
		"ClusterName": *currentModel.Name,
	}
	secretName, err := util.CreateDeploymentSecret(&req, resourceID, *currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey, &resourceProps)
	if err != nil {
		log.Infof("Create - CreateDeploymentSecret - error: %+v", err)
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          err.Error(),
			HandlerErrorCode: cloudformation.HandlerErrorCodeServiceInternalError}, nil
	}

	log.Infof("Created new deployment secret for cluster. Secert Name = Cluster Id:%s", *secretName)
	currentModel.Id = secretName

	//Create Cluster
	cluster, _, err := client.AdvancedClusters.Create(context.Background(), projectID, clusterRequest)
	if err != nil {
		log.Infof("Create - Cluster.Create() - error: %+v", err)
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          err.Error(),
			HandlerErrorCode: cloudformation.HandlerErrorCodeInvalidRequest}, nil
	}

	currentModel.StateName = &cluster.StateName
	log.Debugf("Inprogress model.Id : %s", *currentModel.Id)
	event := handler.ProgressEvent{
		OperationStatus:      handler.InProgress,
		Message:              fmt.Sprintf("Create Cluster `%s`", cluster.StateName),
		ResourceModel:        currentModel,
		CallbackDelaySeconds: 20,
		CallbackContext: map[string]interface{}{
			"stateName":   cluster.StateName,
			"projectId":   projectID,
			"clusterName": *currentModel.Name,
			Id:            *currentModel.Id,
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

	// Callback is not called - Create() and Update() get recalled on
	// long running operations
	callback := map[string]interface{}(req.CallbackContext)
	log.Debugf("Read -  callback: %v", callback)
	if currentModel.Id == nil {
		err := errors.New("No Id found in currentModel")
		log.Infof("Read - error: %+v", err)
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          err.Error(),
			HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound}, nil
	}
	secretName := *currentModel.Id
	log.Infof("Read for Cluster Id/SecretName:%s", secretName)
	key, err := util.GetApiKeyFromDeploymentSecret(&req, secretName)
	if err != nil {
		log.Infof("Read - error: %+v", err)
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          err.Error(),
			HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound}, nil
	}
	log.Debugf("key:%+v", key)

	//Create Client
	client, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		return progress_events.GetFailedEventByCode(fmt.Sprintf("Error creating mongoDB client : %s", err.Error()),
			cloudformation.HandlerErrorCodeInvalidRequest), nil
	}

	// currentModel is NOT populated on the Read after long-running Cluster create
	// need to parse pid and cluster name from Id (deployment secret name).

	//projectID := *currentModel.ProjectId
	//clusterName := *currentModel.Name

	// key.ResourceID should == *currentModel.Id
	id, err := util.ParseResourceIdentifier(*currentModel.Id)
	if err != nil {
		log.Infof("Read - error: %+v", err)
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          err.Error(),
			HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound}, nil
	}
	log.Debugf("Parsed resource identifier: id:%+v", id)

	projectID := id.Parent.ResourceID
	clusterName := id.ResourceID

	log.Debugf("Got projectID:%s, clusterName:%s, from id:%+v", projectID, clusterName, id)

	//Read call
	model, resp, err := readCluster(context.Background(), client, currentModel)
	if err != nil {
		log.Errorf("Cluster.Read() - error: %+v", err)
		if resp != nil && resp.StatusCode == 404 {
			log.Infof("error 404- err:%+v resp:%+v", err, resp)
			return handler.ProgressEvent{
				Message:          err.Error(),
				OperationStatus:  handler.Failed,
				HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound}, nil
		} else {
			log.Infof("error cluster get- err:%+v resp:%+v", err, resp)
			return handler.ProgressEvent{
				Message:          err.Error(),
				OperationStatus:  handler.Failed,
				HandlerErrorCode: cloudformation.HandlerErrorCodeServiceInternalError}, nil
		}
	}
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "Read Complete",
		ResourceModel:   model,
	}, nil
}

// Update handles the Update event from the Cloudformation service.
func Update(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup()
	log.Debugf("Update() currentModel:%+v", currentModel)

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

	//Update callback
	if _, ok := req.CallbackContext["stateName"]; ok {
		return updateClusterCallback(req, client, currentModel, *currentModel.ProjectId)
	}

	//Update Cluster
	model, resp, err := updateCluster(context.Background(), client, currentModel)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Infof("update 404 err: %+v", err)
			return handler.ProgressEvent{
				Message:          err.Error(),
				OperationStatus:  handler.Failed,
				HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound}, nil
		} else {
			log.Infof("update err: %+v", err)
			code := cloudformation.HandlerErrorCodeServiceInternalError
			if strings.Contains(err.Error(), "not exist") { // cfn test needs 404
				code = cloudformation.HandlerErrorCodeNotFound
			}
			if strings.Contains(err.Error(), "being deleted") {
				code = cloudformation.HandlerErrorCodeNotFound // cfn test needs 404
			}
			return handler.ProgressEvent{
				Message:          err.Error(),
				OperationStatus:  handler.Failed,
				HandlerErrorCode: code}, nil
		}
	}

	event := handler.ProgressEvent{
		OperationStatus:      handler.InProgress,
		Message:              fmt.Sprintf("Update Cluster `%s`", model.StateName),
		ResourceModel:        model,
		CallbackDelaySeconds: 65,
		CallbackContext: map[string]interface{}{
			"stateName":   model.StateName,
			"projectId":   currentModel.ProjectId,
			"clusterName": *currentModel.Name,
		},
	}
	log.Debugf("Update() return event:%+v", event)
	return event, nil
}

// Delete handles the Delete event from the Cloudformation service.
func Delete(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup()
	//log.Debugf("Delete() currentModel:%+v", currentModel)

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
	ctx := context.Background()

	if _, ok := req.CallbackContext["stateName"]; ok {
		return validateProgress(client, req, currentModel, "DELETED", "DELETING")
	}

	resp, err := client.AdvancedClusters.Delete(ctx, *currentModel.ProjectId, *currentModel.Name)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Infof("Delete 404 err: %+v", err)
			return handler.ProgressEvent{
				Message:          err.Error(),
				OperationStatus:  handler.Failed,
				HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound}, nil
		} else {
			log.Infof("Delete err: %+v", err)
			return handler.ProgressEvent{
				Message:          err.Error(),
				OperationStatus:  handler.Failed,
				HandlerErrorCode: cloudformation.HandlerErrorCodeServiceInternalError}, nil
		}
	}
	mm := fmt.Sprintf("%s-Deleting", *currentModel.ProjectId)
	currentModel.Id = &mm

	return handler.ProgressEvent{
		OperationStatus:      handler.InProgress,
		Message:              "Delete in progress",
		ResourceModel:        currentModel,
		CallbackDelaySeconds: 60,
		CallbackContext: map[string]interface{}{
			"stateName": "DELETING",
		}}, nil
}

// List handles the List event from the Cloudformation service.
func List(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	return handler.ProgressEvent{
		OperationStatus:  handler.Failed,
		Message:          "Not implemented",
		HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound}, nil
}

func createClusterCallback(req handler.Request, client *mongodbatlas.Client, currentModel *Model, projectID string) (handler.ProgressEvent, error) {
	progressEvent, err := validateProgress(client, req, currentModel, "IDLE", "CREATING")
	if err != nil {
		return progressEvent, nil
	}
	if progressEvent.Message == "Complete" {
		log.Debugf("Compelted creation:%s", *currentModel.Name)

		cluster, res, err := client.AdvancedClusters.Get(context.Background(), projectID, *currentModel.Name)
		if err != nil {
			return progress_events.GetFailedEventByResponse(fmt.Sprintf("Error creating resource : %s", err.Error()),
				res.Response), nil
		}

		log.Debugf("Updating cluster :%s", *currentModel.Name)

		return updateClusterSettings(currentModel, client, context.Background(), projectID, cluster)
	}
	return progressEvent, nil
}

func containsLabelOrKey(list []mongodbatlas.Label, item mongodbatlas.Label) bool {
	for _, v := range list {
		if reflect.DeepEqual(v, item) || v.Key == item.Key {
			return true
		}
	}

	return false
}

func expandBiConnector(biConnector *BiConnector) *mongodbatlas.BiConnector {
	if biConnector == nil {
		return nil
	}
	return &mongodbatlas.BiConnector{
		Enabled:        biConnector.Enabled,
		ReadPreference: cast.ToString(biConnector.ReadPreference),
	}
}

func expandReplicationSpecs(replicationSpecs []AdvancedReplicationSpec) []*mongodbatlas.AdvancedReplicationSpec {
	var rSpecs []*mongodbatlas.AdvancedReplicationSpec

	for _, s := range replicationSpecs {
		var numShards int

		rSpec := &mongodbatlas.AdvancedReplicationSpec{
			ID:            cast.ToString(s.ID),
			NumShards:     numShards,
			RegionConfigs: expandRegionsConfig(s.AdvancedRegionConfigs),
		}

		if s.NumShards != nil {
			rSpec.NumShards = *s.NumShards
		}
		if s.ZoneName != nil {
			rSpec.ZoneName = cast.ToString(s.ZoneName)
		}
		rSpecs = append(rSpecs, rSpec)
	}
	spew.Dump(rSpecs)
	fmt.Printf("specs: len %d %+v", len(replicationSpecs), rSpecs)
	return rSpecs
}

func expandAutoScaling(scaling *AdvancedAutoScaling) *mongodbatlas.AdvancedAutoScaling {

	advAutoScaling := &mongodbatlas.AdvancedAutoScaling{}
	if scaling.Compute != nil {
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

		advAutoScaling.Compute = &mongodbatlas.Compute{
			Enabled:          scaling.Compute.Enabled,
			ScaleDownEnabled: scaling.Compute.ScaleDownEnabled,
			MinInstanceSize:  minInstanceSize,
			MaxInstanceSize:  maxInstanceSize,
		}
	}
	if scaling.DiskGB != nil {
		advAutoScaling.DiskGB = &mongodbatlas.DiskGB{Enabled: scaling.DiskGB.Enabled}
	}
	return advAutoScaling
}

func expandRegionsConfig(regionConfigs []AdvancedRegionConfig) []*mongodbatlas.AdvancedRegionConfig {
	var regionsConfigs []*mongodbatlas.AdvancedRegionConfig
	for _, regionCfg := range regionConfigs {
		regionsConfigs = append(regionsConfigs, expandRegionConfig(regionCfg))
	}
	return regionsConfigs
}

func expandRegionConfig(regionCfg AdvancedRegionConfig) *mongodbatlas.AdvancedRegionConfig {
	var region string
	if regionCfg.RegionName != nil {
		region = *regionCfg.RegionName
	}
	advRegionConfig := &mongodbatlas.AdvancedRegionConfig{
		ProviderName: AWS,
		RegionName:   region,
		Priority:     regionCfg.Priority,
	}

	if regionCfg.AutoScaling != nil {
		advRegionConfig.AutoScaling = expandAutoScaling(regionCfg.AutoScaling)
	}
	if regionCfg.AnalyticsSpecs != nil {
		advRegionConfig.AnalyticsSpecs = expandRegionConfigSpec(regionCfg.AnalyticsSpecs)
	}
	if regionCfg.ElectableSpecs != nil {
		advRegionConfig.ElectableSpecs = expandRegionConfigSpec(regionCfg.ElectableSpecs)
	}
	if regionCfg.ReadOnlySpecs != nil {
		advRegionConfig.ReadOnlySpecs = expandRegionConfigSpec(regionCfg.ReadOnlySpecs)
	}
	return advRegionConfig
}

func expandRegionConfigSpec(spec *Specs) *mongodbatlas.Specs {
	spew.Dump(spec)
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
	var val int64
	if spec.DiskIOPS != nil {
		v, err := strconv.ParseInt(*spec.DiskIOPS, 10, 64)
		if err == nil {
			val = v
		}
		log.Debugf("set diskIops %d", val)
	}
	return &mongodbatlas.Specs{
		DiskIOPS:      &val,
		EbsVolumeType: ebsVolumeType,
		InstanceSize:  instanceSize,
		NodeCount:     spec.NodeCount,
	}
}

func expandLabelSlice(labels []Labels) []mongodbatlas.Label {

	res := make([]mongodbatlas.Label, len(labels))

	for i, _ := range labels {
		var key string
		if labels[i].Key != nil {
			key = *labels[i].Key
		}
		var value string
		if labels[i].Key != nil {
			value = *labels[i].Value
		}
		res[i] = mongodbatlas.Label{
			Key:   key,
			Value: value,
		}
	}
	return res
}

func flattenAutoScaling(scaling *mongodbatlas.AdvancedAutoScaling) *AdvancedAutoScaling {

	if scaling == nil {
		return nil
	}
	advAutoScaling := &AdvancedAutoScaling{}
	if scaling.DiskGB != nil {
		advAutoScaling.DiskGB = &DiskGB{Enabled: scaling.DiskGB.Enabled}
	}
	if scaling.Compute != nil {
		compute := &Compute{}
		if scaling.Compute.Enabled != nil {
			compute.Enabled = scaling.Compute.Enabled
		}
		if scaling.Compute.ScaleDownEnabled != nil {
			compute.ScaleDownEnabled = scaling.Compute.ScaleDownEnabled
		}
		if scaling.Compute.MinInstanceSize != "" {
			compute.MinInstanceSize = &scaling.Compute.MinInstanceSize
		}
		if scaling.Compute.MaxInstanceSize != "" {
			compute.MaxInstanceSize = &scaling.Compute.MaxInstanceSize
		}

		advAutoScaling.Compute = compute
	}
	return advAutoScaling
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
	advRegConfig := AdvancedRegionConfig{
		AutoScaling: flattenAutoScaling(regionCfg.AutoScaling),
		RegionName:  &regionCfg.RegionName,
		Priority:    regionCfg.Priority,
	}
	if regionCfg.AnalyticsSpecs != nil {
		advRegConfig.AnalyticsSpecs = flattenRegionConfigSpec(regionCfg.AnalyticsSpecs)
	}
	if regionCfg.ElectableSpecs != nil {
		advRegConfig.ElectableSpecs = flattenRegionConfigSpec(regionCfg.ElectableSpecs)
	}

	if regionCfg.ReadOnlySpecs != nil {
		advRegConfig.ReadOnlySpecs = flattenRegionConfigSpec(regionCfg.ReadOnlySpecs)
	}

	return advRegConfig
}

func flattenRegionConfigSpec(spec *mongodbatlas.Specs) *Specs {
	if spec == nil {
		return nil
	}
	var diskIops string
	if spec.DiskIOPS != nil {
		diskIops = strconv.FormatInt(*spec.DiskIOPS, 10)
		log.Debugf("get diskIops %s", diskIops)

	}
	return &Specs{
		DiskIOPS:      &diskIops,
		EbsVolumeType: &spec.EbsVolumeType,
		InstanceSize:  &spec.InstanceSize,
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

func flattenProcessArgs(p *mongodbatlas.ProcessArgs) *ProcessArgs {
	spew.Dump(p)
	return &ProcessArgs{
		DefaultReadConcern:               &p.DefaultReadConcern,
		DefaultWriteConcern:              &p.DefaultWriteConcern,
		FailIndexKeyTooLong:              p.FailIndexKeyTooLong,
		JavascriptEnabled:                p.JavascriptEnabled,
		MinimumEnabledTLSProtocol:        &p.MinimumEnabledTLSProtocol,
		NoTableScan:                      p.NoTableScan,
		OplogSizeMB:                      castNO64(p.OplogSizeMB),
		SampleSizeBIConnector:            castNO64(p.SampleSizeBIConnector),
		SampleRefreshIntervalBIConnector: castNO64(p.SampleRefreshIntervalBIConnector),
	}
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

func readCluster(ctx context.Context, client *mongodbatlas.Client, currentModel *Model) (*Model, *mongodbatlas.Response, error) {

	cluster, res, err := client.AdvancedClusters.Get(ctx, *currentModel.ProjectId, *currentModel.Name)

	if err != nil || res.StatusCode != 200 {
		return currentModel, res, err
	}
	currentModel.ProjectId = &cluster.GroupID
	currentModel.Name = &cluster.Name

	if currentModel.BackupEnabled != nil {
		currentModel.BackupEnabled = cluster.BackupEnabled
	}
	if currentModel.BiConnector != nil {
		currentModel.BiConnector = flattenBiConnectorConfig(cluster.BiConnector)
	}
	//Readonly
	currentModel.ConnectionStrings = flattenConnectionStrings(cluster.ConnectionStrings)
	if currentModel.ClusterType != nil {
		currentModel.ClusterType = &cluster.ClusterType
	}
	//Readonly
	currentModel.CreatedDate = &cluster.CreateDate
	if currentModel.DiskSizeGB != nil {
		currentModel.DiskSizeGB = cluster.DiskSizeGB
	}
	if currentModel.EncryptionAtRestProvider != nil {
		currentModel.EncryptionAtRestProvider = &cluster.EncryptionAtRestProvider
	}
	if currentModel.Labels != nil {
		currentModel.Labels = flattenLabels(removeLabel(cluster.Labels, defaultLabel))
	}
	if currentModel.MongoDBMajorVersion != nil {
		currentModel.MongoDBMajorVersion = &cluster.MongoDBMajorVersion
	}
	//Readonly
	currentModel.MongoDBVersion = &cluster.MongoDBVersion

	if currentModel.Paused != nil {
		currentModel.Paused = cluster.Paused
	}
	if currentModel.PitEnabled != nil {
		currentModel.PitEnabled = cluster.PitEnabled
	}
	if currentModel.RootCertType != nil {
		currentModel.RootCertType = &cluster.RootCertType
	}
	if currentModel.ReplicationSpecs != nil {
		currentModel.ReplicationSpecs = flattenReplicationSpecs(cluster.ReplicationSpecs)
	}
	//Readonly
	currentModel.StateName = &cluster.StateName
	if currentModel.VersionReleaseSystem != nil {
		currentModel.VersionReleaseSystem = &cluster.VersionReleaseSystem
	}

	if currentModel.AdvancedSettings != nil {
		processArgs, res, err := client.Clusters.GetProcessArgs(ctx, *currentModel.ProjectId, *currentModel.Name)
		if err != nil || res.StatusCode != 200 {
			return currentModel, res, err
		}
		currentModel.AdvancedSettings = flattenProcessArgs(processArgs)
	}
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

func updateCluster(ctx context.Context, client *mongodbatlas.Client, currentModel *Model) (*Model, *mongodbatlas.Response, error) {

	clusterRequest := &mongodbatlas.AdvancedCluster{}

	if currentModel.BackupEnabled != nil {
		clusterRequest.BackupEnabled = currentModel.BackupEnabled
	}

	if currentModel.BiConnector != nil {
		clusterRequest.BiConnector = expandBiConnector(currentModel.BiConnector)
	}

	if currentModel.ClusterType != nil {
		clusterRequest.ClusterType = *currentModel.ClusterType
	}

	if currentModel.DiskSizeGB != nil {
		clusterRequest.DiskSizeGB = currentModel.DiskSizeGB
	}

	if currentModel.EncryptionAtRestProvider != nil {
		clusterRequest.EncryptionAtRestProvider = cast.ToString(currentModel.EncryptionAtRestProvider)
	}

	if len(currentModel.Labels) > 0 {
		clusterRequest.Labels = expandLabelSlice(currentModel.Labels)
		if containsLabelOrKey(clusterRequest.Labels, defaultLabel) {
			log.Infof("Update - error :%s", labelError)
			return nil, nil, errors.New(labelError)
		}
	}

	if currentModel.MongoDBMajorVersion != nil {
		clusterRequest.MongoDBMajorVersion = formatMongoDBMajorVersion(*currentModel.MongoDBMajorVersion)
	}

	if currentModel.PitEnabled != nil {
		clusterRequest.PitEnabled = currentModel.PitEnabled
	}

	if currentModel.ReplicationSpecs != nil {
		clusterRequest.ReplicationSpecs = expandReplicationSpecs(currentModel.ReplicationSpecs)
	}

	if currentModel.RootCertType != nil {
		clusterRequest.RootCertType = *currentModel.RootCertType
	}

	if currentModel.VersionReleaseSystem != nil {
		clusterRequest.VersionReleaseSystem = *currentModel.VersionReleaseSystem
	}

	_, resp, err := client.AdvancedClusters.Update(ctx, *currentModel.ProjectId, *currentModel.Name, clusterRequest)
	return currentModel, resp, err
}

func updateAdvancedCluster(ctx context.Context, conn *mongodbatlas.Client, request *mongodbatlas.AdvancedCluster, projectID, name string) (*mongodbatlas.AdvancedCluster, *mongodbatlas.Response, error) {
	cluster, resp, err := conn.AdvancedClusters.Update(ctx, projectID, name, request)
	if err != nil {
		return nil, nil, err
	}

	return cluster, resp, nil
}

func updateClusterCallback(req handler.Request, client *mongodbatlas.Client, currentModel *Model, projectID string) (handler.ProgressEvent, error) {
	progressEvent, err := validateProgress(client, req, currentModel, "IDLE", "UPDATING")
	if err != nil {
		return progressEvent, nil
	}
	if progressEvent.Message == "Complete" {
		log.Debugf("Compelted creation:%s", *currentModel.Name)

		cluster, res, err := client.AdvancedClusters.Get(context.Background(), projectID, *currentModel.Name)
		if err != nil {
			return progress_events.GetFailedEventByResponse(fmt.Sprintf("Error in Get Cluster : %s", err.Error()),
				res.Response), nil
		}

		log.Debugf("Updating cluster :%s", *currentModel.Name)

		return updateClusterSettings(currentModel, client, context.Background(), projectID, cluster)
	}
	return progressEvent, nil
}

func updateClusterSettings(currentModel *Model, client *mongodbatlas.Client, ctx context.Context, projectID string, cluster *mongodbatlas.AdvancedCluster) (handler.ProgressEvent, error) {
	//Update advanced configuration
	if currentModel.AdvancedSettings != nil {
		log.Debugf("AdvancedSettings: %+v", *currentModel.AdvancedSettings)
		spew.Dump(*currentModel.AdvancedSettings)
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

	jsonStr, _ := json.Marshal(currentModel)
	log.Infof("Cluster Response --- value: %s ", jsonStr)
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         fmt.Sprintf("Create Cluster `%s`", cluster.StateName),
		ResourceModel:   currentModel,
	}, nil
}

func validateProgress(client *mongodbatlas.Client, req handler.Request, currentModel *Model, targetState string, pendingState string) (handler.ProgressEvent, error) {
	log.Debugf(" Cluster validateProgress() currentModel:%+v", currentModel)
	isReady, state, _, err := isClusterInTargetState(client, *currentModel.ProjectId, *currentModel.Name, targetState)
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
			Id:          currentModel.Id,
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
