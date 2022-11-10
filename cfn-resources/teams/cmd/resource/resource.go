package resource

import (
	"context"
	"fmt"
	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util/validator"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	matlasClient "go.mongodb.org/atlas/mongodbatlas"
)

var CreateRequiredFields = []string{"ApiKeys.PublicKey", "ApiKeys.PrivateKey", "OrgId"}
var ReadRequiredFields = []string{"ApiKeys.PublicKey", "ApiKeys.PrivateKey", "OrgId", "TeamId"}
var UpdateRequiredFields = []string{"ApiKeys.PublicKey", "ApiKeys.PrivateKey", "OrgId", "TeamId"}
var DeleteRequiredFields = []string{"ApiKeys.PublicKey", "ApiKeys.PrivateKey", "OrgId", "TeamId"}
var ListRequiredFields = []string{"ApiKeys.PublicKey", "ApiKeys.PrivateKey", "OrgId"}

const (
	errorTeamCreate = "error creating Team information: %s"
	errorTeamRead   = "error getting Team information: %s"
	errorTeamUpdate = "error updating Team information: %s"
	errorTeamDelete = "error deleting Team (%s): %s"
)

func Create(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup() //logger setup

	log.Debugf("Delete encryption for Request() currentModel:%+v", currentModel)
	// Validate required fields in the request
	modelValidation := validateModel(CreateRequiredFields, currentModel)
	if modelValidation != nil {
		return *modelValidation, nil
	}
	// Create MongoDb Atlas Client using keys
	client, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		log.Errorf("Create - error: %+v", err)
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          err.Error(),
			HandlerErrorCode: cloudformation.HandlerErrorCodeInvalidRequest}, nil

	}
	// Create Atlas API Request Object
	teamRequest := &matlasClient.Team{
		Name:      *currentModel.Name,
		Usernames: currentModel.Usernames,
	}
	team, teamsResp, err := client.Teams.Create(context.Background(), *currentModel.OrgId,
		teamRequest)
	if err != nil && err.Error() != "<nil>" {
		if err != nil {
			return handler.ProgressEvent{}, fmt.Errorf(errorTeamCreate, err)
		}
	}
	currentModel.TeamId = &team.ID
	log.Info("Created Successfully - (%s)", teamsResp.Body)

	event := handler.ProgressEvent{
		OperationStatus: handler.Success,
		ResourceModel:   currentModel,
	}
	log.Infof("Create() return event:%+v", event)
	return event, nil
}
func Read(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup() //logger setup

	log.Debugf("Delete encryption for Request() currentModel:%+v", currentModel)
	// Validate required fields in the request
	modelValidation := validateModel(ReadRequiredFields, currentModel)
	if modelValidation != nil {
		return *modelValidation, nil
	}
	// Create MongoDb Atlas Client using keys
	client, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		log.Errorf("Create - error: %+v", err)
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          err.Error(),
			HandlerErrorCode: cloudformation.HandlerErrorCodeInvalidRequest}, nil

	}

	isExist := isExist(currentModel)
	// Check if snapshot already exist due to this issue https://github.com/mongodb/go-client-mongodb-atlas/issues/315
	if !isExist {
		log.Infof(errorTeamRead, *currentModel.TeamId)
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          "Resource Not Found",
			HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound}, nil
	}

	// API call to read snapshot
	team, _, err := client.Teams.Get(context.Background(), *currentModel.OrgId, *currentModel.TeamId)
	if err != nil {
		log.Infof(errorTeamRead, err)
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          "Resource Not Found",
			HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound}, nil
	}

	users, _, err := client.Teams.GetTeamUsersAssigned(context.Background(), *currentModel.OrgId, *currentModel.TeamId)
	if err != nil {
		log.Infof(errorTeamRead, err)

	}
	if users != nil {
		var newUsers []string
		for i := 0; i < len(users); i++ {
			newUsers = append(newUsers, users[i].ID)

		}
		currentModel.Usernames = newUsers
	}

	currentModel.TeamId = &team.ID
	currentModel.Name = &team.Name

	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "Read Complete",
		ResourceModel:   currentModel,
	}, nil

}
func Update(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup() //logger setup

	log.Debugf("Delete encryption for Request() currentModel:%+v", currentModel)
	// Validate required fields in the request
	modelValidation := validateModel(UpdateRequiredFields, currentModel)
	if modelValidation != nil {
		return *modelValidation, nil
	}
	// Create MongoDb Atlas Client using keys
	client, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		log.Errorf("Create - error: %+v", err)
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          err.Error(),
			HandlerErrorCode: cloudformation.HandlerErrorCodeInvalidRequest}, nil

	}

	isExist := isExist(currentModel)
	// Check if snapshot already exist due to this issue https://github.com/mongodb/go-client-mongodb-atlas/issues/315
	if !isExist {
		log.Infof(errorTeamRead, *currentModel.TeamId)
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          "Resource Not Found",
			HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound}, nil
	}

	// API call to read snapshot
	team, res, err := client.Teams.Get(context.Background(), *currentModel.OrgId, *currentModel.TeamId)
	if err != nil {
		log.Infof(errorTeamRead, err)
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          "Resource Not Found",
			HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound}, nil
	}
	log.Infof("Read -reading snapshot status (%d)", res.StatusCode)
	if team.Name != *currentModel.Name {
		_, _, err := client.Teams.Rename(context.Background(), *currentModel.OrgId, *currentModel.TeamId, *currentModel.Name)
		if err != nil {
			log.Infof(errorTeamUpdate, err)
		}
	}
	isEqualValue := isEqual(team.Usernames, currentModel.Usernames)
	if !isEqualValue {
		// Get the current team's users
		users, _, err := client.Teams.GetTeamUsersAssigned(context.Background(), *currentModel.OrgId, *currentModel.TeamId)

		if err != nil {
			log.Infof("Read -reading snapshot status (%d)", res.StatusCode)
		}
		usernames := currentModel.Usernames
		var newUsers []string
		for i := 0; i < len(usernames); i++ {
			currentUser, isExistingUser := isUserExist(users, usernames[i])
			if isExistingUser {
				_, err := client.Teams.RemoveUserToTeam(context.Background(), *currentModel.OrgId, *currentModel.TeamId, currentUser.ID)
				if err != nil {
					log.Infof("Read -reading snapshot status (%d)", res.StatusCode)
				}
			} else {
				user, _, err := client.AtlasUsers.GetByName(context.Background(), usernames[i])
				if err != nil {
					log.Infof("Read -reading snapshot status (%d)", res.StatusCode)
				}
				// if the user exists, we will storage its ID
				newUsers = append(newUsers, user.ID)

			}

		}
		_, _, err = client.Teams.AddUsersToTeam(context.Background(), *currentModel.OrgId, *currentModel.TeamId, newUsers)
		if err != nil {
			log.Infof("Read -reading snapshot status (%d)", res.StatusCode)
		}

	}
	rolenames := currentModel.RoleNames
	if len(rolenames) > 0 {
		teamRequest := &matlasClient.TeamUpdateRoles{RoleNames: rolenames}
		_, _, err = client.Teams.UpdateTeamRoles(context.Background(), *currentModel.OrgId, *currentModel.TeamId, teamRequest)
		if err != nil {
			log.Infof("Read -reading snapshot status (%d)", res.StatusCode)
		}
	}
	event := handler.ProgressEvent{
		OperationStatus: handler.Success,
		ResourceModel:   currentModel,
	}
	return event, nil
}
func List(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup() //logger setup

	log.Debugf("Delete encryption for Request() currentModel:%+v", currentModel)
	// Validate required fields in the request
	modelValidation := validateModel(ListRequiredFields, currentModel)
	if modelValidation != nil {
		return *modelValidation, nil
	}
	// Create MongoDb Atlas Client using keys
	client, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		log.Errorf("Create - error: %+v", err)
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          err.Error(),
			HandlerErrorCode: cloudformation.HandlerErrorCodeInvalidRequest}, nil

	}

	// Create Atlas API Request Object
	orgId := *currentModel.OrgId
	params := &matlasClient.ListOptions{
		PageNum:      0,
		ItemsPerPage: 100,
	}
	// API call to read snapshot
	teams, _, err := client.Teams.List(context.Background(), orgId, params)
	if err != nil {
		return handler.ProgressEvent{}, fmt.Errorf("error reading teamlist with  id(Organization: %s): %s", orgId, err)
	}
	var models []interface{}

	for i := 0; i < len(teams); i++ {
		var model Model
		model.TeamId = &teams[i].ID
		model.Name = &teams[i].Name
		model.Usernames = teams[i].Usernames
		fmt.Println(teams[i])
		models = append(models, model)
	}
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "List Complete",
		ResourceModels:  models,
	}, nil
}
func Delete(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup() //logger setup

	log.Debugf("Delete encryption for Request() currentModel:%+v", currentModel)
	// Validate required fields in the request
	modelValidation := validateModel(DeleteRequiredFields, currentModel)
	if modelValidation != nil {
		return *modelValidation, nil
	}
	// Create MongoDb Atlas Client using keys
	client, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		log.Errorf("Create - error: %+v", err)
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          err.Error(),
			HandlerErrorCode: cloudformation.HandlerErrorCodeInvalidRequest}, nil

	}

	isExist := isExist(currentModel)
	// Check if snapshot already exist due to this issue https://github.com/mongodb/go-client-mongodb-atlas/issues/315
	if !isExist {
		log.Infof(errorTeamDelete, currentModel.TeamId)
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          "Resource Not Found",
			HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound,
		}, nil
	}
	_, err = client.Teams.RemoveTeamFromOrganization(context.Background(), *currentModel.OrgId, *currentModel.TeamId)
	if err != nil {
		var target *matlasClient.ErrorResponse
		if errors.As(err, &target) && target.ErrorCode == "CANNOT_DELETE_TEAM_ASSIGNED_TO_PROJECT" {
			projectID, err := getProjectIDByTeamID(context.Background(), client, *currentModel.TeamId)
			if err != nil {
				log.Infof(errorTeamDelete, currentModel.TeamId)
				return handler.ProgressEvent{
					OperationStatus:  handler.Failed,
					Message:          "Unable to Delete",
					HandlerErrorCode: cloudformation.HandlerErrorCodeInternalFailure,
				}, nil
			}

			_, err = client.Teams.RemoveTeamFromProject(context.Background(), projectID, *currentModel.TeamId)
			if err != nil {
				log.Infof(errorTeamDelete, currentModel.TeamId)
				return handler.ProgressEvent{
					OperationStatus:  handler.Failed,
					Message:          "Unable to Delete",
					HandlerErrorCode: cloudformation.HandlerErrorCodeInternalFailure,
				}, nil
			}

		}

	}
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "Delete Complete",
	}, nil

}
func setup() {
	util.SetupLogger("mongodb-atlas-project")

}
func validateProgress(client *matlasClient.Client, currentModel *Model, targetState string) (handler.ProgressEvent, error) {
	isReady, state, err := teamIsReady(client, currentModel, targetState)
	if err != nil {
		return handler.ProgressEvent{}, err
	}

	if !isReady {
		p := handler.NewProgressEvent()
		p.ResourceModel = currentModel
		p.OperationStatus = handler.InProgress
		p.CallbackDelaySeconds = 35
		p.Message = "Pending"
		p.CallbackContext = map[string]interface{}{
			"status":      state,
			"snapshot_id": *currentModel.Name,
		}
		return p, nil
	}

	p := handler.NewProgressEvent()
	p.ResourceModel = currentModel
	p.OperationStatus = handler.Success
	p.Message = "Complete"
	return p, nil
}
func isExist(currentModel *Model) bool {

	client, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		return false
	}
	team, _, err := client.Teams.Get(context.Background(), *currentModel.OrgId, *currentModel.TeamId)
	if err != nil {
		return false
	}
	if team != nil {
		return true
	}

	return false
}
func isUserExist(users []matlasClient.AtlasUser, username string) (matlasClient.AtlasUser, bool) {

	for _, user := range users {
		log.Infof("Read - errors reading snapshot with id %s", username)
		if user.Username == username {
			return user, true
		}

	}

	return matlasClient.AtlasUser{}, false
}
func teamIsReady(client *matlasClient.Client, currentModel *Model, targetState string) (bool, string, error) {
	snapshotRequest := &matlasClient.SnapshotReqPathParameters{
		GroupID:     "",
		SnapshotID:  "",
		ClusterName: "",
	}

	snapshot, resp, err := client.CloudProviderSnapshots.GetOneCloudProviderSnapshot(context.Background(), snapshotRequest)
	if err != nil {
		if snapshot == nil && resp == nil {
			return false, "", err
		}
		if resp != nil && resp.StatusCode == 404 {
			return true, "deleted", nil
		}
		return false, "", err
	}
	return snapshot.Status == targetState, snapshot.Status, nil
}
func getProjectIDByTeamID(ctx context.Context, conn *matlasClient.Client, teamID string) (string, error) {
	options := &matlasClient.ListOptions{}
	projects, _, err := conn.Projects.GetAllProjects(ctx, options)
	if err != nil {
		return "", fmt.Errorf("error getting projects information: %s", err)
	}

	for _, project := range projects.Results {
		teams, _, err := conn.Projects.GetProjectTeamsAssigned(ctx, project.ID)
		if err != nil {
			return "", fmt.Errorf("error getting teams from project information: %s", err)
		}

		for _, team := range teams.Results {
			if team.TeamID == teamID {
				return project.ID, nil
			}
		}
	}

	return "", nil
}

// function to validate inputs to all actions
// function to validate inputs to all actions
func validateModel(fields []string, model *Model) *handler.ProgressEvent {
	return validator.ValidateModel(fields, model)
}
func isEqual(first, second []string) bool {
	if len(first) != len(second) {
		return false
	}
	exists := make(map[string]bool)
	for _, value := range first {
		exists[value] = true
	}
	for _, value := range second {
		if !exists[value] {
			return false
		}
	}
	return true
}
