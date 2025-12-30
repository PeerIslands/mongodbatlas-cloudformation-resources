# Template Test Inferences

This folder contains scripts and templates for testing the `MongoDB::Atlas::SearchDeployment` resource using CloudFormation templates after it's been deployed to the private registry.

**Note:** This CloudFormation resource is designed for AWS deployments.

## Purpose

These scripts help developers:
- Test resource creation using CloudFormation templates
- Verify resource updates (NodeCount and InstanceSize can be updated)
- Verify resource deletion
- Understand how the resource works once deployed in the private registry

## Prerequisites

1. Resource must be submitted to AWS Private Registry:
   ```bash
   export AWS_DEFAULT_REGION=eu-west-1
   export AWS_REGION=eu-west-1
   source /Users/home/repos/PeerIslands/Mongo-TF-CFN-Converter/CONVERSION_PROMPTS/setup-credentials.sh /Users/home/repos/PeerIslands/Mongo-TF-CFN-Converter/CONVERSION_PROMPTS/credPersonalCfnDev.properties
   export MONGODB_ATLAS_CLUSTER_NAME='cfn-test-search-deployment-20251229'
   cd /Users/home/repos/PeerIslands/Mongo-TF-CFN-Converter/mongodbatlas-cloudformation-resources/cfn-resources
   LOG_FILE="search-deployment/cfn-submit-$(date +%Y%m%d-%H%M%S).log"
   script -q "$LOG_FILE" bash -c './cfn-submit-helper.sh search-deployment'
   ```

2. **Credentials Configuration** (assume credentials are available as environment variables):

   **Option A: Use Environment Variables**
   ```bash
   export AWS_ACCESS_KEY_ID=your-aws-access-key
   export AWS_SECRET_ACCESS_KEY=your-aws-secret-key
   export AWS_DEFAULT_REGION=eu-west-1
   export MONGODB_ATLAS_PUBLIC_API_KEY=your-atlas-public-key
   export MONGODB_ATLAS_PRIVATE_API_KEY=your-atlas-private-key
   export MONGODB_ATLAS_PROJECT_ID=your-project-id
   export MONGODB_ATLAS_CLUSTER_NAME=your-cluster-name
   ./test-lifecycle.sh
   ```

   **Option B: Use AWS/Atlas CLI Configuration**
   ```bash
   # Configure AWS CLI
   aws configure
   
   # Configure Atlas CLI
   atlas config set public_api_key your-key
   atlas config set private_api_key your-secret
   
   ./test-lifecycle.sh
   ```

3. AWS Secrets Manager secret configured (for CloudFormation resource):
   ```bash
   aws secretsmanager create-secret \
     --name "cfn/atlas/profile/default" \
     --secret-string '{"PublicKey":"your-key","PrivateKey":"your-secret"}'
   ```
   Note: This is for the CloudFormation resource itself. The script uses Atlas CLI directly.

4. IP Access List configured (for testing):
   - Add `0.0.0.0/1` and `128.0.0.0/1` to your Atlas project's IP access list
   - See [README.md](../../../README.md) for details

5. Atlas CLI installed:
   ```bash
   # Install from: https://www.mongodb.com/docs/atlas/cli/stable/atlas-cli-install/
   ```

6. **Existing Cluster Required**:
   - The cluster must already exist in Atlas
   - Search deployment is created for an existing cluster
   - Only one search deployment can exist per cluster

## Usage

### Complete Lifecycle Testing (Interactive)

```bash
./test-lifecycle.sh
```

Or with environment variables:

```bash
export MONGODB_ATLAS_PROJECT_ID=<your-project-id>
export MONGODB_ATLAS_CLUSTER_NAME=<your-cluster-name>
export STACK_NAME=<optional-stack-name>
export AWS_REGION=eu-west-1
export PROFILE=default
./test-lifecycle.sh
```

This interactive script will:

1. **Prompt for inputs** (if not provided via environment variables):
   - MongoDB Atlas Project ID (required)
   - Cluster Name (required - must be an existing cluster)
   - CloudFormation Stack Name (optional, auto-generated if not provided)
   - AWS Region (defaults to eu-west-1)
   - Profile (defaults to default)

2. **CREATE START** → Create search deployment with initial config (InstanceSize=S30_HIGHCPU_NVME, NodeCount=2)
   - Creates CloudFormation stack
   - Monitors stack creation progress (can take 5-15 minutes for search deployment)
   - Validates search deployment in Atlas CLI
   - Verifies InstanceSize and NodeCount configuration
   - Saves resource information to `output.json`
   - **CREATE COMPLETE** → Waits for user confirmation (type 'yes' to continue)

3. **UPDATE START** → Update search deployment (NodeCount: 2 → 3)
   - Updates CloudFormation stack
   - Monitors stack update progress (can take 5-15 minutes)
   - Validates updated configuration in Atlas CLI
   - Verifies NodeCount was updated
   - Updates `output.json` with latest information
   - **UPDATE COMPLETE** → Waits for user confirmation (type 'yes' to continue)

4. **DELETE START** → Delete search deployment
   - Deletes CloudFormation stack
   - Monitors stack deletion progress (can take 5-15 minutes)
   - Validates deletion in Atlas CLI
   - Confirms search deployment no longer exists
   - Cleans up `output.json` file
   - **DELETE COMPLETE**

### Features

- ✅ **Runs CREATE first**, then waits for user confirmation
- ✅ **UPDATE is supported** - NodeCount and InstanceSize can be updated
- ✅ **Runs UPDATE**, then waits for user confirmation
- ✅ **Runs DELETE**, then completes
- ✅ **Accepts inputs** via prompts or environment variables
- ✅ **Automatic cleanup** of failed stacks (ROLLBACK_COMPLETE, CREATE_FAILED, DELETE_FAILED)
- ✅ **Real-time progress** monitoring with status updates
- ✅ **Atlas CLI validation** after each step
- ✅ **Detailed error reporting** with stack events
- ✅ **Colored output** for better readability
- ✅ **Saves resource info** to `output.json` after create/update
- ✅ **Reads from `output.json`** for reference during lifecycle

## Files

- `search-deployment-complex.json` - Complex CloudFormation template with all configuration fields
- `test-lifecycle.sh` - **Main lifecycle test script** (Create → Update → Delete with user confirmation)
- `README.md` - This file

## Workflow

```
┌─────────────────────────────────────┐
│  ./test-lifecycle.sh                │
│  (Prompts for inputs)               │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  CREATE START                       │
│  - Create search deployment         │
│    (S30_HIGHCPU_NVME, NodeCount=2)  │
│  - Validate with Atlas CLI          │
│  - Save to output.json              │
│  CREATE COMPLETE                     │
│  ⏸️  Wait for user confirmation      │
│     (yes / exit)                    │
└──────────────┬──────────────────────┘
               │
               ▼
        ┌──────────────┐
        │ User types   │
        │ 'yes'        │
        └──────┬───────┘
               │
               ▼
┌─────────────────────────────────────┐
│  UPDATE START                       │
│  - Update search deployment         │
│    (NodeCount: 2 -> 3)              │
│  - Validate with Atlas CLI          │
│  - Update output.json               │
│  UPDATE COMPLETE                     │
│  ⏸️  Wait for user confirmation      │
│     (yes / exit)                    │
└──────────────┬──────────────────────┘
               │
               ▼
        ┌──────────────┐
        │ User types   │
        │ 'yes'        │
        └──────┬───────┘
               │
               ▼
┌─────────────────────────────────────┐
│  DELETE START                       │
│  - Delete search deployment         │
│  - Validate deletion in Atlas CLI   │
│  - Clean up output.json             │
│  DELETE COMPLETE                     │
└─────────────────────────────────────┘
```

## Atlas CLI Commands for Verification

### List Search Deployments
```bash
atlas clusters search nodes list \
  --clusterName <CLUSTER_NAME> \
  --projectId <PROJECT_ID> \
  --output json
```

### Expected Output Structure
```json
{
  "id": "search-deployment-id",
  "stateName": "IDLE",
  "specs": [
    {
      "instanceSize": "S30_HIGHCPU_NVME",
      "nodeCount": 2
    }
  ],
  "encryptionAtRestProvider": "AWS"
}
```

## Important Notes

1. **Update Support**: Unlike some resources, SearchDeployment supports updates:
   - ✅ `NodeCount` can be updated
   - ✅ `InstanceSize` can be updated
   - ❌ `ProjectId`, `ClusterName`, `Profile` are create-only (require replacement)

2. **Timing**: Search deployment operations can take 5-15 minutes:
   - CREATE: Typically 5-10 minutes
   - UPDATE: Typically 5-15 minutes (depending on changes)
   - DELETE: Typically 5-10 minutes

3. **One Per Cluster**: Only one search deployment can exist per cluster. If you try to create a second one, it will fail.

4. **Cluster Must Exist**: The cluster must already exist before creating a search deployment.

## Example Output

```
============================================
CREATE START
============================================

[INFO] Creating CloudFormation stack with initial configuration...
[INFO] Initial Config: InstanceSize=S30_HIGHCPU_NVME, NodeCount=2
[SUCCESS] Stack creation initiated
[1/120] Status: CREATE_IN_PROGRESS
[2/120] Status: CREATE_IN_PROGRESS
...
[15/120] Status: CREATE_COMPLETE
[SUCCESS] Stack created successfully!

[INFO] Stack Outputs:
{
  "SearchDeploymentId": "...",
  "SearchDeploymentStateName": "IDLE",
  ...
}

[INFO] Validating search deployment in Atlas...
[SUCCESS] Search deployment found for cluster 'my-cluster' in Atlas
[SUCCESS] InstanceSize matches: S30_HIGHCPU_NVME
[SUCCESS] NodeCount matches: 2

============================================
CREATE COMPLETE
============================================
```

## Troubleshooting

### Stack Creation Fails
- Check CloudWatch logs for the handler Lambda function
- Verify AWS Secrets Manager secret exists: `cfn/atlas/profile/default`
- Verify cluster exists in Atlas
- Check IP Access List allows your IP

### Atlas CLI Validation Fails
- Verify Atlas CLI credentials are configured
- Check project ID is correct
- Verify cluster name is correct
- Ensure search deployment was actually created (check Atlas UI)

### Update Fails
- Verify the resource supports updates (NodeCount and InstanceSize can be updated)
- Check that you're not trying to update create-only properties (ProjectId, ClusterName, Profile)

