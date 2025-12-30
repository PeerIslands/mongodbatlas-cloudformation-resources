// Copyright 2024 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resource

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	admin20250312010 "go.mongodb.org/atlas-sdk/v20250312010/admin"
	"go.mongodb.org/atlas-sdk/v20250312010/mockadmin"
)

// Test helpers
func createTestSearchDeploymentModel() *Model {
	projectID := "507f1f77bcf86cd799439011"
	clusterName := "test-cluster"
	profile := "default"
	instanceSize := "S20_HIGHCPU_NVME"
	nodeCount := 2

	return &Model{
		Profile:     &profile,
		ProjectId:   &projectID,
		ClusterName: &clusterName,
		Specs: []ApiSearchDeploymentSpec{
			{
				InstanceSize: &instanceSize,
				NodeCount:    &nodeCount,
			},
		},
	}
}

func createTestSearchDeploymentResponse() *admin20250312010.ApiSearchDeploymentResponse {
	id := "test-id-123"
	stateName := "IDLE"
	return &admin20250312010.ApiSearchDeploymentResponse{
		Id:        &id,
		StateName: &stateName,
		Specs: &[]admin20250312010.ApiSearchDeploymentSpec{
			{
				InstanceSize: "S20_HIGHCPU_NVME",
				NodeCount:    2,
			},
		},
	}
}

// Basic function tests
func TestSetup(t *testing.T) {
	assert.NotPanics(t, func() {
		setup()
	})
}

func TestList(t *testing.T) {
	req := handler.Request{}
	event, err := List(req, nil, nil)

	require.Error(t, err)
	assert.Equal(t, "not implemented: List", err.Error())
	assert.Equal(t, handler.ProgressEvent{}, event)
}

func TestInProgressEvent(t *testing.T) {
	model := createTestSearchDeploymentModel()
	apiResp := createTestSearchDeploymentResponse()
	event := inProgressEvent("Test Message", model, apiResp)

	assert.Equal(t, handler.InProgress, event.OperationStatus)
	assert.Equal(t, "Test Message", event.Message)
	assert.Equal(t, int64(callBackSeconds), event.CallbackDelaySeconds)
	assert.NotNil(t, event.CallbackContext)
	assert.Contains(t, event.CallbackContext, "callbackSearchDeployment")
	assert.NotNil(t, event.ResourceModel)
}

func TestHandleError(t *testing.T) {
	createSDKError := func(errorCode string, statusCode int) *admin20250312010.GenericOpenAPIError {
		apiErr := admin20250312010.ApiError{
			Error:     statusCode,
			ErrorCode: errorCode,
		}
		sdkErr := &admin20250312010.GenericOpenAPIError{}
		sdkErr.SetModel(apiErr)
		return sdkErr
	}

	testCases := map[string]struct {
		response          *http.Response
		err               error
		expectedStatus    handler.Status
		expectedErrorCode string
	}{
		"AlreadyExistsError": {
			response:          &http.Response{StatusCode: http.StatusBadRequest},
			err:               createSDKError(SearchDeploymentAlreadyExistsErrorAPI, http.StatusBadRequest),
			expectedStatus:    handler.Failed,
			expectedErrorCode: "AlreadyExists",
		},
		"DoesNotExistError": {
			response:          &http.Response{StatusCode: http.StatusBadRequest},
			err:               createSDKError(SearchDeploymentDoesNotExistsErrorAPI, http.StatusBadRequest),
			expectedStatus:    handler.Failed,
			expectedErrorCode: "NotFound",
		},
		"GenericError": {
			response:       &http.Response{StatusCode: http.StatusInternalServerError},
			err:            fmt.Errorf("internal server error"),
			expectedStatus: handler.Failed,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			event, err := handleError(tc.response, tc.err)

			require.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, event.OperationStatus)
			if tc.expectedErrorCode != "" {
				assert.Equal(t, tc.expectedErrorCode, event.HandlerErrorCode)
			}
		})
	}
}

func TestConstants(t *testing.T) {
	assert.Equal(t, 40, callBackSeconds)
	assert.Equal(t, "ATLAS_SEARCH_DEPLOYMENT_DOES_NOT_EXIST", SearchDeploymentDoesNotExistsErrorAPI)
	assert.Equal(t, "ATLAS_SEARCH_DEPLOYMENT_ALREADY_EXISTS", SearchDeploymentAlreadyExistsErrorAPI)
}

func TestRequiredFields(t *testing.T) {
	assert.Equal(t, []string{constants.ProjectID, constants.ClusterName, constants.Specs}, createRequiredFields)
	assert.Equal(t, []string{constants.ProjectID, constants.ClusterName}, readRequiredFields)
	assert.Equal(t, []string{constants.ProjectID, constants.ClusterName, constants.Specs}, updateRequiredFields)
	assert.Equal(t, []string{constants.ProjectID, constants.ClusterName}, deleteRequiredFields)
}

// Validation tests
func TestCreateValidationErrors(t *testing.T) {
	testCases := map[string]struct {
		currentModel *Model
		expectedMsg  string
	}{
		"missingProjectId":   {&Model{ClusterName: util.StringPtr("test-cluster"), Specs: []ApiSearchDeploymentSpec{{InstanceSize: util.StringPtr("S20_HIGHCPU_NVME"), NodeCount: util.IntPtr(2)}}}, "required"},
		"missingClusterName": {&Model{ProjectId: util.StringPtr("507f1f77bcf86cd799439011"), Specs: []ApiSearchDeploymentSpec{{InstanceSize: util.StringPtr("S20_HIGHCPU_NVME"), NodeCount: util.IntPtr(2)}}}, "required"},
		"missingSpecs":       {&Model{ProjectId: util.StringPtr("507f1f77bcf86cd799439011"), ClusterName: util.StringPtr("test-cluster")}, "required"},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			event, err := Create(handler.Request{}, nil, tc.currentModel)
			require.NoError(t, err)
			assert.Equal(t, handler.Failed, event.OperationStatus)
			assert.Contains(t, event.Message, tc.expectedMsg)
		})
	}
}

func TestReadValidationErrors(t *testing.T) {
	testCases := map[string]struct {
		currentModel *Model
		expectedMsg  string
	}{
		"missingProjectId":   {&Model{ClusterName: util.StringPtr("test-cluster")}, "required"},
		"missingClusterName": {&Model{ProjectId: util.StringPtr("507f1f77bcf86cd799439011")}, "required"},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			event, err := Read(handler.Request{}, nil, tc.currentModel)
			require.NoError(t, err)
			assert.Equal(t, handler.Failed, event.OperationStatus)
			assert.Contains(t, event.Message, tc.expectedMsg)
		})
	}
}

func TestUpdateValidationErrors(t *testing.T) {
	testCases := map[string]struct {
		currentModel *Model
		expectedMsg  string
	}{
		"missingProjectId":   {&Model{ClusterName: util.StringPtr("test-cluster"), Specs: []ApiSearchDeploymentSpec{{InstanceSize: util.StringPtr("S20_HIGHCPU_NVME"), NodeCount: util.IntPtr(2)}}}, "required"},
		"missingClusterName": {&Model{ProjectId: util.StringPtr("507f1f77bcf86cd799439011"), Specs: []ApiSearchDeploymentSpec{{InstanceSize: util.StringPtr("S20_HIGHCPU_NVME"), NodeCount: util.IntPtr(2)}}}, "required"},
		"missingSpecs":       {&Model{ProjectId: util.StringPtr("507f1f77bcf86cd799439011"), ClusterName: util.StringPtr("test-cluster")}, "required"},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			event, err := Update(handler.Request{}, nil, tc.currentModel)
			require.NoError(t, err)
			assert.Equal(t, handler.Failed, event.OperationStatus)
			assert.Contains(t, event.Message, tc.expectedMsg)
		})
	}
}

func TestDeleteValidationErrors(t *testing.T) {
	testCases := map[string]struct {
		currentModel *Model
		expectedMsg  string
	}{
		"missingProjectId":   {&Model{ClusterName: util.StringPtr("test-cluster")}, "required"},
		"missingClusterName": {&Model{ProjectId: util.StringPtr("507f1f77bcf86cd799439011")}, "required"},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			event, err := Delete(handler.Request{}, nil, tc.currentModel)
			require.NoError(t, err)
			assert.Equal(t, handler.Failed, event.OperationStatus)
			assert.Contains(t, event.Message, tc.expectedMsg)
		})
	}
}

// CRUD operation tests with mocks
func TestCreateWithMocks(t *testing.T) {
	originalInitEnv := initEnvWithClient
	defer func() { initEnvWithClient = originalInitEnv }()

	testCases := map[string]struct {
		req            handler.Request
		mockSetup      func(*mockadmin.AtlasSearchApi)
		expectedStatus handler.Status
		validateResult func(t *testing.T, event handler.ProgressEvent)
	}{
		"successfulCreate": {
			req: handler.Request{},
			mockSetup: func(m *mockadmin.AtlasSearchApi) {
				// Create returns response with IDLE state, so handler returns SUCCESS immediately
				idleResp := createTestSearchDeploymentResponse()
				stateName := "IDLE"
				idleResp.StateName = &stateName
				m.EXPECT().CreateClusterSearchDeployment(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(admin20250312010.CreateClusterSearchDeploymentApiRequest{ApiService: m})
				m.EXPECT().CreateClusterSearchDeploymentExecute(mock.Anything).
					Return(idleResp, &http.Response{StatusCode: 200}, nil)
			},
			expectedStatus: handler.Success,
			validateResult: func(t *testing.T, event handler.ProgressEvent) {
				assert.Equal(t, constants.Complete, event.Message)
				assert.NotNil(t, event.ResourceModel)
			},
		},
		"createWithCallback": {
			req: handler.Request{CallbackContext: map[string]interface{}{"callbackSearchDeployment": true}},
			mockSetup: func(m *mockadmin.AtlasSearchApi) {
				m.EXPECT().GetClusterSearchDeployment(mock.Anything, mock.Anything, mock.Anything).
					Return(admin20250312010.GetClusterSearchDeploymentApiRequest{ApiService: m})
				m.EXPECT().GetClusterSearchDeploymentExecute(mock.Anything).
					Return(createTestSearchDeploymentResponse(), &http.Response{StatusCode: 200}, nil)
			},
			expectedStatus: handler.Success,
		},
		"createWithError": {
			req: handler.Request{},
			mockSetup: func(m *mockadmin.AtlasSearchApi) {
				m.EXPECT().CreateClusterSearchDeployment(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(admin20250312010.CreateClusterSearchDeploymentApiRequest{ApiService: m})
				m.EXPECT().CreateClusterSearchDeploymentExecute(mock.Anything).
					Return(nil, &http.Response{StatusCode: 500}, fmt.Errorf("API error"))
			},
			expectedStatus: handler.Failed,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockSearchApi := mockadmin.NewAtlasSearchApi(t)
			tc.mockSetup(mockSearchApi)

			mockClient := &admin20250312010.APIClient{AtlasSearchApi: mockSearchApi}
			initEnvWithClient = func(req handler.Request, currentModel *Model, requiredFields []string) (*admin20250312010.APIClient, *handler.ProgressEvent) {
				return mockClient, nil
			}

			event, err := Create(tc.req, nil, createTestSearchDeploymentModel())
			require.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, event.OperationStatus)
			if tc.validateResult != nil {
				tc.validateResult(t, event)
			}
		})
	}
}

func TestReadWithMocks(t *testing.T) {
	originalInitEnv := initEnvWithClient
	defer func() { initEnvWithClient = originalInitEnv }()

	testCases := map[string]struct {
		mockSetup      func(*mockadmin.AtlasSearchApi)
		expectedStatus handler.Status
	}{
		"successfulRead": {
			mockSetup: func(m *mockadmin.AtlasSearchApi) {
				m.EXPECT().GetClusterSearchDeployment(mock.Anything, mock.Anything, mock.Anything).
					Return(admin20250312010.GetClusterSearchDeploymentApiRequest{ApiService: m})
				m.EXPECT().GetClusterSearchDeploymentExecute(mock.Anything).
					Return(createTestSearchDeploymentResponse(), &http.Response{StatusCode: 200}, nil)
			},
			expectedStatus: handler.Success,
		},
		"readNotFound": {
			mockSetup: func(m *mockadmin.AtlasSearchApi) {
				m.EXPECT().GetClusterSearchDeployment(mock.Anything, mock.Anything, mock.Anything).
					Return(admin20250312010.GetClusterSearchDeploymentApiRequest{ApiService: m})
				m.EXPECT().GetClusterSearchDeploymentExecute(mock.Anything).
					Return(nil, &http.Response{StatusCode: 404}, fmt.Errorf("not found"))
			},
			expectedStatus: handler.Failed,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockSearchApi := mockadmin.NewAtlasSearchApi(t)
			tc.mockSetup(mockSearchApi)

			mockClient := &admin20250312010.APIClient{AtlasSearchApi: mockSearchApi}
			initEnvWithClient = func(req handler.Request, currentModel *Model, requiredFields []string) (*admin20250312010.APIClient, *handler.ProgressEvent) {
				return mockClient, nil
			}

			event, err := Read(handler.Request{}, nil, createTestSearchDeploymentModel())
			require.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, event.OperationStatus)
		})
	}
}

func TestUpdateWithMocks(t *testing.T) {
	originalInitEnv := initEnvWithClient
	defer func() { initEnvWithClient = originalInitEnv }()

	testCases := map[string]struct {
		req            handler.Request
		mockSetup      func(*mockadmin.AtlasSearchApi)
		expectedStatus handler.Status
	}{
		"successfulUpdate": {
			req: handler.Request{},
			mockSetup: func(m *mockadmin.AtlasSearchApi) {
				m.EXPECT().UpdateClusterSearchDeployment(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(admin20250312010.UpdateClusterSearchDeploymentApiRequest{ApiService: m})
				m.EXPECT().UpdateClusterSearchDeploymentExecute(mock.Anything).
					Return(nil, &http.Response{StatusCode: 200}, nil)
				m.EXPECT().GetClusterSearchDeployment(mock.Anything, mock.Anything, mock.Anything).
					Return(admin20250312010.GetClusterSearchDeploymentApiRequest{ApiService: m})
				m.EXPECT().GetClusterSearchDeploymentExecute(mock.Anything).
					Return(createTestSearchDeploymentResponse(), &http.Response{StatusCode: 200}, nil)
			},
			expectedStatus: handler.InProgress,
		},
		"updateWithCallback": {
			req: handler.Request{CallbackContext: map[string]interface{}{"callbackSearchDeployment": true}},
			mockSetup: func(m *mockadmin.AtlasSearchApi) {
				m.EXPECT().GetClusterSearchDeployment(mock.Anything, mock.Anything, mock.Anything).
					Return(admin20250312010.GetClusterSearchDeploymentApiRequest{ApiService: m})
				m.EXPECT().GetClusterSearchDeploymentExecute(mock.Anything).
					Return(createTestSearchDeploymentResponse(), &http.Response{StatusCode: 200}, nil)
			},
			expectedStatus: handler.Success,
		},
		"updateWithError": {
			req: handler.Request{},
			mockSetup: func(m *mockadmin.AtlasSearchApi) {
				// Update handler now checks resource existence first
				m.EXPECT().GetClusterSearchDeployment(mock.Anything, mock.Anything, mock.Anything).
					Return(admin20250312010.GetClusterSearchDeploymentApiRequest{ApiService: m})
				m.EXPECT().GetClusterSearchDeploymentExecute(mock.Anything).
					Return(createTestSearchDeploymentResponse(), &http.Response{StatusCode: 200}, nil)
				m.EXPECT().UpdateClusterSearchDeployment(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(admin20250312010.UpdateClusterSearchDeploymentApiRequest{ApiService: m})
				m.EXPECT().UpdateClusterSearchDeploymentExecute(mock.Anything).
					Return(nil, &http.Response{StatusCode: 500}, fmt.Errorf("update failed"))
			},
			expectedStatus: handler.Failed,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockSearchApi := mockadmin.NewAtlasSearchApi(t)
			tc.mockSetup(mockSearchApi)

			mockClient := &admin20250312010.APIClient{AtlasSearchApi: mockSearchApi}
			initEnvWithClient = func(req handler.Request, currentModel *Model, requiredFields []string) (*admin20250312010.APIClient, *handler.ProgressEvent) {
				return mockClient, nil
			}

			event, err := Update(tc.req, nil, createTestSearchDeploymentModel())
			require.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, event.OperationStatus)
		})
	}
}

func TestDeleteWithMocks(t *testing.T) {
	originalInitEnv := initEnvWithClient
	defer func() { initEnvWithClient = originalInitEnv }()

	testCases := map[string]struct {
		req            handler.Request
		mockSetup      func(*mockadmin.AtlasSearchApi)
		expectedStatus handler.Status
	}{
		"successfulDelete": {
			req: handler.Request{},
			mockSetup: func(m *mockadmin.AtlasSearchApi) {
				m.EXPECT().DeleteClusterSearchDeployment(mock.Anything, mock.Anything, mock.Anything).
					Return(admin20250312010.DeleteClusterSearchDeploymentApiRequest{ApiService: m})
				m.EXPECT().DeleteClusterSearchDeploymentExecute(mock.Anything).
					Return(&http.Response{StatusCode: 200}, nil)
				// After delete, the handler checks if resource still exists
				m.EXPECT().GetClusterSearchDeployment(mock.Anything, mock.Anything, mock.Anything).
					Return(admin20250312010.GetClusterSearchDeploymentApiRequest{ApiService: m})
				m.EXPECT().GetClusterSearchDeploymentExecute(mock.Anything).
					Return(createTestSearchDeploymentResponse(), &http.Response{StatusCode: 200}, nil)
			},
			expectedStatus: handler.InProgress,
		},
		"deleteWithCallback": {
			req: handler.Request{CallbackContext: map[string]interface{}{"callbackSearchDeployment": true}},
			mockSetup: func(m *mockadmin.AtlasSearchApi) {
				m.EXPECT().GetClusterSearchDeployment(mock.Anything, mock.Anything, mock.Anything).
					Return(admin20250312010.GetClusterSearchDeploymentApiRequest{ApiService: m})
				m.EXPECT().GetClusterSearchDeploymentExecute(mock.Anything).
					Return(nil, &http.Response{StatusCode: 404}, fmt.Errorf("not found"))
			},
			expectedStatus: handler.Success,
		},
		"deleteWithError": {
			req: handler.Request{},
			mockSetup: func(m *mockadmin.AtlasSearchApi) {
				m.EXPECT().DeleteClusterSearchDeployment(mock.Anything, mock.Anything, mock.Anything).
					Return(admin20250312010.DeleteClusterSearchDeploymentApiRequest{ApiService: m})
				m.EXPECT().DeleteClusterSearchDeploymentExecute(mock.Anything).
					Return(&http.Response{StatusCode: 500}, fmt.Errorf("delete failed"))
			},
			expectedStatus: handler.Failed,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mockSearchApi := mockadmin.NewAtlasSearchApi(t)
			tc.mockSetup(mockSearchApi)

			mockClient := &admin20250312010.APIClient{AtlasSearchApi: mockSearchApi}
			initEnvWithClient = func(req handler.Request, currentModel *Model, requiredFields []string) (*admin20250312010.APIClient, *handler.ProgressEvent) {
				return mockClient, nil
			}

			event, err := Delete(tc.req, nil, createTestSearchDeploymentModel())
			require.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, event.OperationStatus)
		})
	}
}
