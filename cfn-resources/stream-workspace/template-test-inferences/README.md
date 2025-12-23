# Template Test Inferences

This folder contains scripts and templates for testing the `MongoDB::Atlas::StreamWorkspace` resource using CloudFormation templates after it's been deployed to the private registry.

**Note:** This CloudFormation resource is designed for AWS deployments. The CloudProvider parameter is constrained to "AWS" only in the template.

## Purpose

These scripts help developers:
- Test resource creation using CloudFormation templates
- Verify resource updates
- Verify resource deletion
- Understand how the resource works once deployed in the private registry

## Prerequisites

1. Resource must be submitted to AWS Private Registry:
   ```bash
   cd ../..
   ./cfn-submit-helper.sh stream-workspace
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

## Usage

### Complete Lifecycle Testing (Interactive)

```bash
./test-lifecycle.sh
```

Or with environment variables:

```bash
export MONGODB_ATLAS_PROJECT_ID=<your-project-id>
export WORKSPACE_NAME=<optional-workspace-name>
export STACK_NAME=<optional-stack-name>
export AWS_REGION=eu-west-1
export PROFILE=default
./test-lifecycle.sh
```

This interactive script will:

1. **Prompt for inputs** (if not provided via environment variables):
   - MongoDB Atlas Project ID (required)
   - Workspace Name (optional, auto-generated if not provided)
   - CloudFormation Stack Name (optional, auto-generated if not provided)
   - AWS Region (defaults to eu-west-1)
   - Profile (defaults to default)

2. **CREATE START** → Create workspace with initial config (SP30, VIRGINIA_USA)
   - Creates CloudFormation stack
   - Monitors stack creation progress
   - Validates workspace in Atlas CLI
   - Verifies tier (SP30) and region configuration
   - Saves resource information to `output.json`
   - **CREATE COMPLETE** → Notes that UPDATE is not supported
   - Waits for your choice (type 'delete' or 'exit')

3. **User Choice After CREATE:**
   - **NOTE: UPDATE is NOT supported** for StreamWorkspace
     - All main properties (WorkspaceName, ProjectId, Profile, StreamConfig) are create-only
     - CloudFormation cannot update resources when create-only properties require replacement
     - To change configuration, you must delete and recreate the workspace
   
   - **If 'delete' chosen** → **DELETE START**
   
   - **If 'exit' chosen** → Exits (stack remains for manual cleanup)

4. **DELETE START** → Delete workspace
   - Deletes CloudFormation stack
   - Monitors stack deletion progress
   - Validates deletion in Atlas CLI
   - Confirms workspace no longer exists
   - Cleans up `output.json` file
   - **DELETE COMPLETE**

### Features

- ✅ **Runs CREATE first**, then waits for user decision
- ✅ **User chooses** DELETE or EXIT after CREATE
- ⚠️ **UPDATE is NOT supported** - StreamConfig is create-only
- ✅ **Accepts inputs** via prompts or environment variables
- ✅ **Automatic cleanup** of failed stacks
- ✅ **Real-time progress** monitoring
- ✅ **Atlas CLI validation** after each step
- ✅ **Detailed error reporting**
- ✅ **Colored output** for better readability
- ✅ **Saves resource info** to `output.json` after create
- ✅ **Reads from `output.json`** for delete operation

## Files

- `stream-workspace-complex.json` - Complex CloudFormation template with all configuration fields
- `test-lifecycle.sh` - **Main lifecycle test script** (Create → Update → Delete with user confirmation)
- `test-results.md` - Test execution results and issues found
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
│  - Create workspace (SP30)         │
│  - Validate with Atlas CLI          │
│  - Save to output.json              │
│  CREATE COMPLETE                     │
│  ⚠️  UPDATE NOT SUPPORTED            │
│     (StreamConfig is create-only)    │
│  ⏸️  Wait for user choice            │
│     (delete / exit)                 │
└──────────────┬──────────────────────┘
               │
               ▼
        ┌──────────────┐
        │ DELETE START │
        │ (if chosen)   │
        │ - Delete      │
        │ - Validate    │
        │ - Cleanup     │
        │ DELETE        │
        │ COMPLETE      │
        └──────────────┘
```

## Example Output

```
=== Configuration Input ===
Enter MongoDB Atlas Project ID: <YOUR_PROJECT_ID>
Enter Workspace Name (or press Enter for auto-generated): 
Auto-generated workspace name: stream-workspace-1703001234
Enter CloudFormation Stack Name (or press Enter for auto-generated): 
Auto-generated stack name: stream-workspace-1703001234

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
CREATE START
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ Stack created successfully!
✅ Workspace 'stream-workspace-1703001234' exists in Atlas
✅ Tier matches: SP30

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
CREATE COMPLETE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

⚠️  NOTE: StreamWorkspace does not support updates
All main properties (WorkspaceName, ProjectId, Profile, StreamConfig) are create-only
To change configuration, you must delete and recreate the workspace

What would you like to do next?
Type 'delete' to delete the workspace, or 'exit' to quit: delete

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
DELETE START
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ Stack deleted successfully!
✅ Workspace 'stream-workspace-1703001234' successfully deleted from Atlas

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
DELETE COMPLETE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```
