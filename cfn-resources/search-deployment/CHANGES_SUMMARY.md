# Search Deployment Resource - Changes Summary

## Overview
This document summarizes all changes made to fix CloudFormation contract test failures and prepare the resource for submission.

## Test Status

### Final Contract Test Results
- **Status**: PASSED ✅
- **Passed**: 8 tests
- **Failed**: 0 tests
- **Skipped**: 4 tests (expected)
- **Duration**: ~57 minutes

## Git Diff Summary

### Files Changed (excluding cfn-testing-helper.sh)
```
13 files changed, 371 insertions(+), 126 deletions(-)
```

### Key Changes

#### 1. `cmd/resource/resource.go` (259 lines changed)

**A. Error Constants Updated:**
```diff
 const (
-	callBackSeconds                    = 40
-	SearchDeploymentDoesNotExistsError = "ATLAS_FTS_DEPLOYMENT_DOES_NOT_EXIST"
-	SearchDeploymentAlreadyExistsError = "ATLAS_FTS_DEPLOYMENT_ALREADY_EXISTS"
+	callBackSeconds                       = 40
+	SearchDeploymentAlreadyExistsErrorAPI = "ATLAS_SEARCH_DEPLOYMENT_ALREADY_EXISTS"
+	SearchDeploymentDoesNotExistsErrorAPI = "ATLAS_SEARCH_DEPLOYMENT_DOES_NOT_EXIST"
 )
```

**B. Added Callback Context:**
```diff
+var callbackContext = map[string]any{"callbackSearchDeployment": true}
+
+func IsCallback(req *handler.Request) bool {
+	_, found := req.CallbackContext["callbackSearchDeployment"]
+	return found
+}
```

**C. Update Handler - Fixed JSON Response Error (Lines 192-205):**
```diff
-	// Check if resource exists before updating - required by contract tests
-	_, checkResp, err := connV2.AtlasSearchApi.GetClusterSearchDeployment(context.Background(), projectID, clusterName).Execute()
+	// Check if resource exists before updating (required by contract tests)
+	checkResp, checkHTTPResp, err := connV2.AtlasSearchApi.GetClusterSearchDeployment(context.Background(), projectID, clusterName).Execute()
 	if err != nil {
-		return handleError(checkResp, err)
+		// If resource doesn't exist, return NotFound (required by contract tests)
+		if checkHTTPResp != nil && checkHTTPResp.StatusCode == http.StatusNotFound {
+			return progressevent.GetFailedEventByResponse("Search deployment not found", checkHTTPResp), nil
+		}
+		return handleError(checkHTTPResp, err)
+	}
+	if checkResp == nil || checkResp.Id == nil {
+		// Resource doesn't exist - return NotFound with proper HTTP response
+		notFoundResp := &http.Response{StatusCode: http.StatusNotFound}
+		return progressevent.GetFailedEventByResponse("Search deployment not found", notFoundResp), nil
 	}
```

**Key Changes:**
- Added explicit check for 404 StatusNotFound before calling handleError
- Changed manual ProgressEvent construction to use `progressevent.GetFailedEventByResponse()` for consistency
- Added proper handling for nil response case
- Aligned with flex-cluster coding standards

**D. Create Handler Improvements:**
```diff
-	// handling of subsequent retry calls
-	if _, ok := req.CallbackContext[constants.ID]; ok {
-		return HandleStateTransition(*connV2, currentModel, constants.IdleState), nil
+	if IsCallback(&req) {
+		return ValidateProgress(*connV2, currentModel, false), nil
 	}
```

- Improved idempotency handling for already-exists case
- Better state checking (returns SUCCESS when deployment is IDLE)
- Enhanced callback context handling

**E. Read Handler - Added NotFound Check:**
```diff
 	if err != nil {
 		return handleError(resp, err)
 	}
+
+	if apiResp == nil || apiResp.Id == nil {
+		return handler.ProgressEvent{
+			OperationStatus:  handler.Failed,
+			Message:          "Search deployment not found",
+			HandlerErrorCode: string(types.HandlerErrorCodeNotFound),
+		}, nil
+	}
```

#### 2. `cmd/resource/resource_test.go` (45 lines changed)

**A. Updated Error Constants in Tests:**
```diff
-			err:               createSDKError(SearchDeploymentAlreadyExistsError, http.StatusBadRequest),
+			err:               createSDKError(SearchDeploymentAlreadyExistsErrorAPI, http.StatusBadRequest),
```

**B. Delete Handler Test - Added Missing Mock Expectation:**
```diff
 		"successfulDelete": {
 			req: handler.Request{},
 			mockSetup: func(m *mockadmin.AtlasSearchApi) {
 				m.EXPECT().DeleteClusterSearchDeployment(mock.Anything, mock.Anything, mock.Anything).
 					Return(admin20250312010.DeleteClusterSearchDeploymentApiRequest{ApiService: m})
 				m.EXPECT().DeleteClusterSearchDeploymentExecute(mock.Anything).
 					Return(&http.Response{StatusCode: 200}, nil)
+				// After delete, the handler checks if resource still exists
+				m.EXPECT().GetClusterSearchDeployment(mock.Anything, mock.Anything, mock.Anything).
+					Return(admin20250312010.GetClusterSearchDeploymentApiRequest{ApiService: m})
+				m.EXPECT().GetClusterSearchDeploymentExecute(mock.Anything).
+					Return(createTestSearchDeploymentResponse(), &http.Response{StatusCode: 200}, nil)
 			},
 			expectedStatus: handler.InProgress,
 		},
```

**C. Updated Callback Context in Tests:**
```diff
-			req: handler.Request{CallbackContext: map[string]interface{}{constants.ID: "test-id"}},
+			req: handler.Request{CallbackContext: map[string]interface{}{"callbackSearchDeployment": true}},
```

**D. Create Test - Updated to Expect SUCCESS for IDLE State:**
```diff
+				// Create returns response with IDLE state, so handler returns SUCCESS immediately
+				idleResp := createTestSearchDeploymentResponse()
+				stateName := "IDLE"
+				idleResp.StateName = &stateName
 				m.EXPECT().CreateClusterSearchDeployment(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
 					Return(admin20250312010.CreateClusterSearchDeploymentApiRequest{ApiService: m})
 				m.EXPECT().CreateClusterSearchDeploymentExecute(mock.Anything).
-					Return(nil, &http.Response{StatusCode: 200}, nil)
-				m.EXPECT().GetClusterSearchDeployment(mock.Anything, mock.Anything, mock.Anything).
-					Return(admin20250312010.GetClusterSearchDeploymentApiRequest{ApiService: m})
-				m.EXPECT().GetClusterSearchDeploymentExecute(mock.Anything).
-					Return(createTestSearchDeploymentResponse(), &http.Response{StatusCode: 200}, nil)
+					Return(idleResp, &http.Response{StatusCode: 200}, nil)
 			},
-			expectedStatus: handler.InProgress,
+			expectedStatus: handler.Success,
```

**E. Update Test - Added Existence Check Mock:**
```diff
 		"updateWithError": {
 			req: handler.Request{},
 			mockSetup: func(m *mockadmin.AtlasSearchApi) {
+				// Update handler now checks resource existence first
+				m.EXPECT().GetClusterSearchDeployment(mock.Anything, mock.Anything, mock.Anything).
+					Return(admin20250312010.GetClusterSearchDeploymentApiRequest{ApiService: m})
+				m.EXPECT().GetClusterSearchDeploymentExecute(mock.Anything).
+					Return(createTestSearchDeploymentResponse(), &http.Response{StatusCode: 200}, nil)
 				m.EXPECT().UpdateClusterSearchDeployment(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
```

#### 3. `examples/search-deployment/search-deployment.json` (1 line changed)

```diff
-  "Description": "This template creates a Search Index on the MongoDB Atlas API, this will be billed to your Atlas account.",
+  "Description": "This template creates a Search Deployment (dedicated search nodes) for a MongoDB Atlas cluster. This will be billed to your Atlas account.",
```

### Other Modified Files

- `cmd/resource/mappings.go` - Updated mappings
- `cmd/resource/state_transition.go` - Enhanced state transition handling
- `cmd/resource/state_transition_test.go` - Updated tests
- `test/cfn-test-create-inputs.sh` - Updated test input generation
- `test/cfn-test-delete-inputs.sh` - Updated test cleanup
- `test/inputs_*.template.json` - Updated test inputs
- `template.yml` - Updated timeout settings
- `resource-role.yaml` - Updated execution role
- `Makefile` - Updated build configuration

### New Files Created

- `CHANGES_SUMMARY.md` - This file
- `template-test-inferences/search-deployment-complex.json` - Complex template
- `template-test-inferences/test-lifecycle.sh` - Lifecycle test script
- `template-test-inferences/README.md` - Template test inferences documentation

## Impact

### Contract Tests
- ✅ Fixed 2 failing tests (`contract_update_read`, `contract_update_tag_updatable`)
- ✅ All 8 contract tests now passing
- ✅ Improved error handling consistency

### Unit Tests
- ✅ Fixed 1 failing test (`TestDeleteWithMocks/successfulDelete`)
- ✅ Updated tests to match new handler behavior
- ✅ All unit tests passing

### Code Quality
- ✅ Aligned with flex-cluster coding standards
- ✅ Consistent error handling using `progressevent.GetFailedEventByResponse()`
- ✅ Better comments explaining contract test requirements
- ✅ Improved callback context handling

## Verification

All changes verified:
- ✅ Unit tests pass: `go test ./cmd/resource/... -v`
- ✅ Contract tests pass: `./cfn-testing-helper.sh search-deployment`
- ✅ Resource submitted to private registry successfully
- ✅ Example templates created and documented
- ✅ Template test inferences created

## Submit Status

✅ **Resource Successfully Submitted**

### Submit Command Used

```bash
export AWS_DEFAULT_REGION=eu-west-1
export AWS_REGION=eu-west-1
source /Users/home/repos/PeerIslands/Mongo-TF-CFN-Converter/CONVERSION_PROMPTS/setup-credentials.sh /Users/home/repos/PeerIslands/Mongo-TF-CFN-Converter/CONVERSION_PROMPTS/credPersonalCfnDev.properties
export MONGODB_ATLAS_CLUSTER_NAME='cfn-test-search-deployment-20251229'
cd /Users/home/repos/PeerIslands/Mongo-TF-CFN-Converter/mongodbatlas-cloudformation-resources/cfn-resources
LOG_FILE="search-deployment/cfn-submit-$(date +%Y%m%d-%H%M%S).log"
script -q "$LOG_FILE" bash -c './cfn-submit-helper.sh search-deployment'
```

### Submit Results

- TypeArn: `arn:aws:cloudformation:eu-west-1:<ACCOUNT_ID>:type/resource/MongoDB-Atlas-SearchDeployment`
- TypeVersionArn: `arn:aws:cloudformation:eu-west-1:<ACCOUNT_ID>:type/resource/MongoDB-Atlas-SearchDeployment/00000001`
- Execution Role Stack: `mongodb-atlas-searchdeployment-role-stack`
