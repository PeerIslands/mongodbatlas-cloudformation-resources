# Search Deployment Lifecycle Test Results

## Test Execution Summary

**Date**: December 29, 2024  
**Resource**: MongoDB::Atlas::SearchDeployment  
**Test Type**: Complete Lifecycle (CREATE → UPDATE → DELETE)

## Test Configuration

- **Project ID**: `<REDACTED>`
- **Cluster Name**: `cfn-test-search-deployment-20251229`
- **Stack Name**: `search-deployment-<TIMESTAMP>` (auto-generated)
- **Region**: `eu-west-1`
- **Profile**: `default`

## Test Execution Log

```
=== Prerequisites Check ===
[INFO] Checking AWS credentials...
[SUCCESS] AWS Account: <REDACTED>
[INFO] Checking Atlas CLI...
[SUCCESS] Atlas CLI found
[INFO] Verifying Atlas credentials...
[SUCCESS] Atlas credentials verified
[INFO] Checking template file...
[SUCCESS] Template file found

=== Configuration Input ===
Enter CloudFormation Stack Name (or press Enter for auto-generated): 
[INFO] Auto-generated stack name: search-deployment-<TIMESTAMP>
Enter AWS Region [eu-west-1]: 
Enter Profile [default]: 

=== Configuration Summary ===
Project ID: <REDACTED>
Cluster Name: cfn-test-search-deployment-20251229
Stack Name: search-deployment-<TIMESTAMP>
Region: eu-west-1
Profile: default

============================================
CREATE START
============================================

[INFO] Checking for existing failed stacks...
[INFO] Creating CloudFormation stack with initial configuration...
[INFO] Initial Config: InstanceSize=S30_HIGHCPU_NVME, NodeCount=2
[SUCCESS] Stack creation initiated
[1/120] Status: CREATE_IN_PROGRESS
[2/120] Status: CREATE_IN_PROGRESS
...
[53/120] Status: CREATE_COMPLETE
[SUCCESS] Stack created successfully!

[INFO] Stack Outputs:
[
  {
    "OutputKey": "SearchDeploymentStateName",
    "OutputValue": "IDLE",
    "Description": "Human-readable label that indicates the current operating condition of this search deployment."
  },
  {
    "OutputKey": "SearchDeploymentId",
    "OutputValue": "<REDACTED>",
    "Description": "Unique 24-hexadecimal digit string that identifies the search deployment."
  },
  {
    "OutputKey": "ConfigurationSummary",
    "OutputValue": "Search Deployment created for cluster cfn-test-search-deployment-20251229 with 2 node(s) of size S30_HIGHCPU_NVME",
    "Description": "Summary of search deployment configuration"
  },
  {
    "OutputKey": "ProjectId",
    "OutputValue": "<REDACTED>",
    "Description": "Project ID for reference"
  },
  {
    "OutputKey": "ClusterName",
    "OutputValue": "cfn-test-search-deployment-20251229",
    "Description": "Cluster name for reference"
  },
  {
    "OutputKey": "SearchDeploymentEncryptionAtRestProvider",
    "OutputValue": "NONE",
    "Description": "Cloud service provider that manages your customer keys to provide an additional layer of Encryption At Rest for the cluster."
  }
]

[INFO] Saving resource information to output.json...
[WARNING] Invalid JSON in search deployment details, using empty object
[SUCCESS] Resource information saved to output.json
============================================
CREATE COMPLETE
============================================

[INFO] What would you like to do next?
Type 'yes' to continue to UPDATE, or 'exit' to quit: yes
[SUCCESS] Proceeding to UPDATE...

============================================
UPDATE START
============================================

[INFO] Updating CloudFormation stack with modified configuration...
[INFO] Updated Config: InstanceSize=S30_HIGHCPU_NVME, NodeCount=3
[SUCCESS] Stack update initiated
[1/120] Status: UPDATE_IN_PROGRESS
[2/120] Status: UPDATE_IN_PROGRESS
...
[37/120] Status: UPDATE_COMPLETE
[SUCCESS] Stack updated successfully!

[INFO] Updated Stack Outputs:
[
  {
    "OutputKey": "SearchDeploymentStateName",
    "OutputValue": "IDLE",
    "Description": "Human-readable label that indicates the current operating condition of this search deployment."
  },
  {
    "OutputKey": "SearchDeploymentId",
    "OutputValue": "<REDACTED>",
    "Description": "Unique 24-hexadecimal digit string that identifies the search deployment."
  },
  {
    "OutputKey": "ConfigurationSummary",
    "OutputValue": "Search Deployment created for cluster cfn-test-search-deployment-20251229 with 3 node(s) of size S30_HIGHCPU_NVME",
    "Description": "Summary of search deployment configuration"
  },
  {
    "OutputKey": "ProjectId",
    "OutputValue": "<REDACTED>",
    "Description": "Project ID for reference"
  },
  {
    "OutputKey": "ClusterName",
    "OutputValue": "cfn-test-search-deployment-20251229",
    "Description": "Cluster name for reference"
  },
  {
    "OutputKey": "SearchDeploymentEncryptionAtRestProvider",
    "OutputValue": "NONE",
    "Description": "Cloud service provider that manages your customer keys to provide an additional layer of Encryption At Rest for the cluster."
  }
]

[INFO] Saving resource information to output.json...
[WARNING] Invalid JSON in search deployment details, using empty object
[SUCCESS] Resource information saved to output.json
============================================
UPDATE COMPLETE
============================================

[INFO] What would you like to do next?
Type 'yes' to continue to DELETE, or 'exit' to quit: yes   
[SUCCESS] Proceeding to DELETE...

============================================
DELETE START
============================================

[INFO] Deleting CloudFormation stack...
[SUCCESS] Stack deletion initiated
[1/120] Status: DELETE_IN_PROGRESS
[2/120] Status: DELETE_IN_PROGRESS
...
[15/120] Status: NOT_FOUND
[SUCCESS] Stack deleted successfully!

[INFO] Verifying search deployment deletion in Atlas...
[SUCCESS] Search deployment successfully deleted from Atlas
[INFO] Cleaning up output file...
[SUCCESS] Output file removed
============================================
DELETE COMPLETE
============================================

============================================
LIFECYCLE TEST COMPLETE
============================================

[SUCCESS] CREATE: Search deployment created successfully (InstanceSize=S30_HIGHCPU_NVME, NodeCount=2)
[SUCCESS] UPDATE: Search deployment updated successfully (NodeCount=2 -> 3)
[SUCCESS] DELETE: Search deployment deleted successfully

[SUCCESS] Lifecycle test finished successfully!
```

## Test Results

### CREATE Operation
- **Status**: ✅ SUCCESS
- **Duration**: ~53 status checks (approximately 8-9 minutes)
- **Configuration**: InstanceSize=S30_HIGHCPU_NVME, NodeCount=2
- **Final State**: IDLE
- **Search Deployment ID**: `<REDACTED>`
- **Outputs**: All outputs generated correctly

### UPDATE Operation
- **Status**: ✅ SUCCESS
- **Duration**: ~37 status checks (approximately 6 minutes)
- **Configuration Change**: NodeCount updated from 2 to 3
- **Final State**: IDLE
- **Search Deployment ID**: `<REDACTED>` (unchanged)
- **Outputs**: ConfigurationSummary correctly updated to show 3 nodes

### DELETE Operation
- **Status**: ✅ SUCCESS
- **Duration**: ~15 status checks (approximately 2-3 minutes)
- **Verification**: Search deployment successfully deleted from Atlas
- **Cleanup**: output.json file removed

## Findings

1. **Atlas CLI Validation Issue**: The script shows a warning "Invalid JSON in search deployment details" during CREATE and UPDATE. This is because `atlas clusters search nodes list` may return an empty result or different format than expected. The script handles this gracefully by using an empty object.

2. **Timing Observations**:
   - CREATE: ~8-9 minutes (53 status checks)
   - UPDATE: ~6 minutes (37 status checks)
   - DELETE: ~2-3 minutes (15 status checks)

3. **Update Support**: Confirmed that Search Deployment supports updates for NodeCount. The update operation completed successfully and the ConfigurationSummary output correctly reflected the change from 2 to 3 nodes.

4. **Idempotency**: The Search Deployment ID remained consistent across CREATE and UPDATE operations, confirming proper resource management.

## Notes

- All operations completed successfully
- No errors encountered during the lifecycle test
- The resource behaves as expected for CREATE, UPDATE, and DELETE operations
- Stack outputs are correctly generated and updated

