package resource

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util/validator"
	"github.com/openlyinc/pointy"
	log "github.com/sirupsen/logrus"
	matlasClient "go.mongodb.org/atlas/mongodbatlas"
)

const (
	publicKey  = "ApiKeys.PublicKey"
	privateKey = "ApiKeys.PrivateKey"
	projectId  = "ProjectId"
	userName   = "UserName"
)

var CreateRequiredFields = []string{publicKey, privateKey, projectId, userName}
var ReadRequiredFields = []string{publicKey, privateKey, projectId}
var DeleteRequiredFields = []string{publicKey, privateKey, projectId}
var ListRequiredFields = []string{publicKey, privateKey, projectId, userName}

const (
	errorX509AuthDBUsersCreate         = "error creating MongoDB X509 Authentication for DB User(%s) in the project(%s): %s"
	errorX509AuthDBUsersRead           = "error reading MongoDB X509 Authentication for DB Users(%s) in the project(%s): %s"
	errorAllX509AuthDBUsersRead        = "error reading all MongoDB X509 certificates for DB Users(%s) in the project(%s): %s"
	errorCustomerX509AuthDBUsersCreate = "error creating Customer X509 Authentication in the project(%s): %s"
	errorCustomerX509AuthDBUsersDelete = "error deleting Customer X509 Authentication in the project(%s): %s"
)

// Create handles the Create event from the Cloudformation service.
func Create(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup() //logger setup

	// Validate required fields in the request
	if modelValidation := validateModel(CreateRequiredFields, currentModel); modelValidation != nil {
		return *modelValidation, nil
	}
	log.Debugf("Create - creating MongoDB X509 Authentication for DB User:%+v", currentModel)

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
	log.Info("Creating MongoDB X509 Authentication for DB User starts")
	// progress callback setup
	if _, ok := req.CallbackContext["status"]; ok {
		sid := req.CallbackContext["ProjectId"].(string)
		currentModel.ProjectId = &sid
		return validateProgress(client, currentModel, "completed")
	}
	projectID := *currentModel.ProjectId
	username := *currentModel.UserName
	expirationMonths := *currentModel.MonthsUntilExpiration
	if expirationMonths > 0 {
		log.Info("Creating User Certificate")
		res, _, err := client.X509AuthDBUsers.CreateUserCertificate(context.Background(), projectID, username, expirationMonths)
		if err != nil {
			log.Errorf(errorX509AuthDBUsersCreate, username, projectID, err)
			return handler.ProgressEvent{
				OperationStatus:  handler.Failed,
				Message:          err.Error(),
				HandlerErrorCode: cloudformation.HandlerErrorCodeInvalidRequest}, nil

		}
		log.Infof("Creating User Certificate +%v", &res)
		log.Infof("Creating User Certificate +%v", res)
		log.Infof("Creating User Certificate %s", &res.Certificate)
		if res != nil {
			currentModel.CustomerX509 = &CustomerX509{
				Cas: pointy.String(res.Certificate),
			}
		}
	} else {
		log.Info("Save Custom Certificate DB User starts")
		customerX509Cas := *currentModel.CustomerX509.Cas
		_, _, err := client.X509AuthDBUsers.SaveConfiguration(context.Background(), projectID, &matlasClient.CustomerX509{Cas: customerX509Cas})
		if err != nil {
			log.Errorf(errorCustomerX509AuthDBUsersCreate, projectID, err)
			return handler.ProgressEvent{
				OperationStatus:  handler.Failed,
				Message:          err.Error(),
				HandlerErrorCode: cloudformation.HandlerErrorCodeInvalidRequest}, nil
		}

	}
	// track progress
	event := handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "Created  Certificate  for DB User ",
		ResourceModel:   currentModel,
	}
	log.Debugf("Create() return event:%+v", event)
	return event, nil

}

// Read handles the Read event from the Cloudformation service.
func Read(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup() //logger setup

	log.Debugf("Read - X509 certificates for Request() :%+v", currentModel)
	// Validate required fields in the request
	if modelValidation := validateModel(ReadRequiredFields, currentModel); modelValidation != nil {
		return *modelValidation, nil
	}

	log.Info("Read - X509 Certificates starts ")
	if isEnabled(currentModel) == false {
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          "config is not available",
			HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound}, nil
	}
	certificate, err := ReadUserX509Certificate(currentModel)
	log.Infof("Read - X509 Certificates starts : +%+v %s  %s ", certificate, *currentModel.UserName, *currentModel.ProjectId)
	if err != nil {
		log.Errorf(errorX509AuthDBUsersRead, *currentModel.UserName, *currentModel.ProjectId, err)
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          err.Error(),
			HandlerErrorCode: cloudformation.HandlerErrorCodeInvalidRequest}, nil

	}
	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "Read Complete",
		ResourceModel:   currentModel,
	}, nil
}

// Read handles the Read event from the Cloudformation service.
func ReadUserX509Certificate(currentModel *Model) (*Model, error) {
	setup() //logger setup

	log.Debugf("Read - X509 certificates for Request() :%+v", currentModel)
	// Create MongoDb Atlas Client using keys
	client, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		log.Errorf("Create - error: %+v", err)
		return nil, errors.New("unable to create mongo client")

	}
	// Create Atlas API Request Object
	log.Info("Read - X509 Certificates starts ")
	projectID := *currentModel.ProjectId
	username := *currentModel.UserName
	params := &matlasClient.ListOptions{
		PageNum:      0,
		ItemsPerPage: 100,
	}
	certificates, resp, err := client.X509AuthDBUsers.GetUserCertificates(context.Background(), projectID, username, params)
	if err != nil {
		log.Errorf("Create - error: %+v", err)
		return nil, fmt.Errorf(errorAllX509AuthDBUsersRead, *currentModel.UserName, projectID, err)

	}
	currentModel.Links = flattenLinks(resp.Links)
	flattenCertificates(certificates, currentModel)

	certificate, _, err := client.X509AuthDBUsers.GetCurrentX509Conf(context.Background(), projectID)
	log.Infof("Read - X509 Certificates starts : %+v ", certificate)
	if err != nil {
		log.Errorf(errorX509AuthDBUsersRead, *currentModel.UserName, projectID, err)
		return nil, fmt.Errorf(errorX509AuthDBUsersRead, *currentModel.UserName, projectID, err)
	} else if certificate != nil {
		currentModel.CustomerX509 = &CustomerX509{
			Cas: &certificate.Cas,
		}
	}

	return currentModel, nil
}

// Update handles the Update event from the Cloudformation service.
func Update(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	// Not implemented, return an empty handler.ProgressEvent
	// and an error
	return handler.ProgressEvent{}, errors.New("not implemented")
}

// Delete handles the Delete event from the Cloudformation service.
func Delete(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	setup() //logger setup

	log.Debugf("Delete - X509 Certificates  for Request() currentModel:%+v", currentModel)
	// Validate required fields in the request
	if modelValidation := validateModel(DeleteRequiredFields, currentModel); modelValidation != nil {
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
	if isEnabled(currentModel) == false {
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          "config is not available",
			HandlerErrorCode: cloudformation.HandlerErrorCodeNotFound}, nil
	}
	log.Info("Delete - X509 Certificates  starts ")
	projectID := *currentModel.ProjectId
	_, err = client.X509AuthDBUsers.DisableCustomerX509(context.Background(), projectID)
	if err != nil {
		log.Errorf(errorCustomerX509AuthDBUsersDelete, projectID, *currentModel.UserName)
		return handler.ProgressEvent{
			OperationStatus:  handler.Failed,
			Message:          "Unable to Delete",
			HandlerErrorCode: cloudformation.HandlerErrorCodeInternalFailure,
		}, nil
	}

	return handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "Delete Complete",
	}, nil

}

// List handles the List event from the Cloudformation service.
func List(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	// Not implemented, return an empty handler.ProgressEvent
	// and an error
	return handler.ProgressEvent{}, errors.New("not implemented")
}

// function to validate inputs to all actions
func validateModel(fields []string, model *Model) *handler.ProgressEvent {
	return validator.ValidateModel(fields, model)
}

func setup() {
	util.SetupLogger("mongodb-atlas-project")

}
func flattenLinks(linksResult []*matlasClient.Link) []Links {
	if linksResult != nil {
		links := make([]Links, 0)
		for _, link := range linksResult {
			var lin Links
			lin.Href = &link.Href
			lin.Rel = &link.Rel
			links = append(links, lin)
		}
		return links
	}
	return nil
}
func flattenCertificates(userCertificates []matlasClient.UserCertificate, currentModel *Model) *Model {
	if userCertificates != nil {
		certificates := make([]Certificate, len(userCertificates))
		for i, _ := range userCertificates {
			v := userCertificates[i]
			id := fmt.Sprintf("%v", &v.ID)
			role := Certificate{
				Id:        &id,
				CreatedAt: &v.CreatedAt,
				GroupId:   &v.GroupID,
				NotAfter:  &v.NotAfter,
				Subject:   &v.Subject,
			}

			certificates = append(certificates, role)
		}
		currentModel.Results = certificates
		currentModel.TotalCount = pointy.Int(len(userCertificates))
	}
	return currentModel
}

// function to track snapshot creation status
func validateProgress(client *matlasClient.Client, currentModel *Model, targetState string) (handler.ProgressEvent, error) {
	projectId := *currentModel.ProjectId
	isReady, state, err := certificateIsReady(client, projectId, targetState)
	if err != nil {
		return handler.ProgressEvent{}, err
	}

	if !isReady {
		p := handler.NewProgressEvent()
		p.ResourceModel = currentModel
		p.OperationStatus = handler.InProgress
		p.CallbackDelaySeconds = 10
		p.Message = "Pending"
		p.CallbackContext = map[string]interface{}{
			"status":    state,
			"ProjectId": *currentModel.ProjectId,
		}
		return p, nil
	}

	p := handler.NewProgressEvent()
	p.ResourceModel = currentModel
	p.OperationStatus = handler.Success
	p.Message = "Complete"
	return p, nil
}

// Read handles the Read event from the Cloudformation service.
func isEnabled(currentModel *Model) bool {
	setup() //logger setup

	log.Debugf("Read - X509 certificates for Request() :%+v", currentModel)
	// Create MongoDb Atlas Client using keys
	client, err := util.CreateMongoDBClient(*currentModel.ApiKeys.PublicKey, *currentModel.ApiKeys.PrivateKey)
	if err != nil {
		log.Errorf("Create - error: %+v", err)
		return false

	}
	projectID := *currentModel.ProjectId

	certificate, _, err := client.X509AuthDBUsers.GetCurrentX509Conf(context.Background(), projectID)
	log.Infof("Read - X509 Certificates starts : %+v ", certificate)
	if err != nil {
		log.Errorf(errorX509AuthDBUsersRead, *currentModel.UserName, projectID, err)
		return false
	} else if certificate != nil && certificate.Cas != "" {
		return true
	}

	return false
}

// function to check if snapshot already exist in atlas
func certificateIsReady(client *matlasClient.Client, projectId, targetState string) (bool, string, error) {

	certificate, resp, err := client.X509AuthDBUsers.GetCurrentX509Conf(context.Background(), projectId)
	if err != nil {
		if certificate == nil && resp == nil {
			return false, "", err
		}
		if resp != nil && resp.StatusCode == 404 {
			return true, "deleted", nil
		}
		return false, "", err
	}
	return resp.StatusCode == 200, "completed", nil
}
func intPtr(i int) *int {
	return &i
}
