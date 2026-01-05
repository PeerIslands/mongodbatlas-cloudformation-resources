# Template Test Inferences for MongoDB::Atlas::OrgServiceAccount

This folder contains templates and scripts for testing the complete lifecycle of the `MongoDB::Atlas::OrgServiceAccount` CloudFormation resource after it has been submitted to the private registry.

## Purpose

This folder is for **developer testing** after the resource has been submitted to AWS CloudFormation Private Registry. It provides:
- Complex template with all configuration fields
- Lifecycle test script for CREATE → UPDATE → DELETE operations
- Validation with Atlas CLI

## Prerequisites

### Required
1. **Resource submitted to private registry**: The `MongoDB::Atlas::OrgServiceAccount` resource must be registered in AWS CloudFormation Private Registry
2. **AWS credentials configured**: 
   - AWS CLI configured with appropriate credentials
   - IAM permissions for CloudFormation operations (CreateStack, UpdateStack, DeleteStack, DescribeStacks)
3. **AWS Secrets Manager secret**: Secret containing Atlas API keys
   - Secret name format: `cfn/atlas/profile/{ProfileName}`
   - Secret structure: `{"PublicKey":"<key>","PrivateKey":"<secret>","BaseURL":"https://cloud.mongodb.com"}`
4. **Atlas CLI installed and configured**:
   ```bash
   # Install from: https://www.mongodb.com/docs/atlas/cli/stable/atlas-cli-install/
   atlas auth login
   ```
5. **MongoDB Atlas Organization**: You must have access to a MongoDB Atlas organization
   - Organization ID (24-character hexadecimal string)
   - API keys with appropriate permissions

### Environment Variables

The script uses the following environment variables (can be set or prompted):

- `MONGODB_ATLAS_ORG_ID` (required): Your MongoDB Atlas organization ID
- `AWS_REGION` (optional, default: `eu-west-1`): AWS region for CloudFormation operations
- `PROFILE` (optional, default: `default`): AWS Secrets Manager profile name

## Files in This Folder

- **`org-service-account-complex.json`**: Complex CloudFormation template with all configuration fields
- **`test-lifecycle.sh`**: Interactive lifecycle test script (CREATE → UPDATE → DELETE)
- **`README.md`**: This file
- **`output.json`**: Generated during test execution (contains resource information, deleted after DELETE)

## Usage

### Complete Lifecycle Testing (Interactive)

```bash
cd mongodbatlas-cloudformation-resources/cfn-resources/org-service-account/template-test-inferences

# Set environment variables (optional)
export MONGODB_ATLAS_ORG_ID=<your-org-id>
export AWS_REGION=eu-west-1
export PROFILE=default

# Run the lifecycle test
./test-lifecycle.sh
```

### Interactive Workflow

The script will prompt you for:

1. **OrgId**: MongoDB Atlas organization ID (defaults to `MONGODB_ATLAS_ORG_ID` env var)
2. **Service Account Name**: Name for the service account (defaults to auto-generated)
3. **Stack Name**: CloudFormation stack name (defaults to auto-generated)
4. **Description**: Description for the service account
5. **Roles**: Comma-separated list of roles (default: `ORG_MEMBER`)
6. **SecretExpiresAfterHours**: Secret expiration in hours (default: 720)

### Workflow Stages

```
┌─────────────────────────────────────────────────────────┐
│                    STAGE 1: CREATE                      │
├─────────────────────────────────────────────────────────┤
│ 1. Create CloudFormation stack                          │
│ 2. Wait for CREATE_COMPLETE                             │
│ 3. Validate service account in Atlas CLI                │
│ 4. Save resource information to output.json            │
│ 5. Display "CREATE COMPLETE"                            │
│ 6. Wait for user confirmation (type 'yes' to continue)   │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                    STAGE 2: UPDATE                       │
├─────────────────────────────────────────────────────────┤
│ 1. Update stack with modified values:                   │
│    - Name: <original> → <original>-updated                │
│    - Description: <original> → <original> - Updated     │
│    - Roles: <original> → ORG_MEMBER,ORG_GROUP_CREATOR    │
│ 2. Wait for UPDATE_COMPLETE                             │
│ 3. Validate updated configuration in Atlas CLI          │
│ 4. Update output.json                                    │
│ 5. Display "UPDATE COMPLETE"                            │
│ 6. Wait for user confirmation (type 'yes' to continue)   │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                    STAGE 3: DELETE                       │
├─────────────────────────────────────────────────────────┤
│ 1. Delete CloudFormation stack                          │
│ 2. Wait for DELETE_COMPLETE                             │
│ 3. Validate deletion in Atlas CLI                       │
│ 4. Clean up output.json                                  │
│ 5. Display "DELETE COMPLETE"                             │
└─────────────────────────────────────────────────────────┘
```

## Example Output

### Stage 1: CREATE

```
[INFO] ==========================================
[INFO] STAGE 1: CREATE
[INFO] ==========================================
[INFO] Creating CloudFormation stack: org-service-account-test-1234567890
[INFO] Waiting for stack creation to complete...
[INFO] Stack status: CREATE_IN_PROGRESS (waiting...)
[SUCCESS] Stack creation completed: CREATE_COMPLETE
[SUCCESS] Stack created successfully
[INFO] ClientId: mdb_sa_id_1234567890abcdef
[INFO] Validating resource in Atlas...
[SUCCESS] Service account found in Atlas:
{
  "clientId": "mdb_sa_id_1234567890abcdef",
  "name": "test-service-account-1234567890",
  "description": "Service account for lifecycle testing",
  "roles": ["ORG_MEMBER"],
  ...
}
[SUCCESS] CREATE COMPLETE
```

### Stage 2: UPDATE

```
[INFO] ==========================================
[INFO] STAGE 2: UPDATE
[INFO] ==========================================
[INFO] Updating CloudFormation stack: org-service-account-test-1234567890
[INFO] Updating Name: test-service-account-1234567890 → test-service-account-1234567890-updated
[INFO] Updating Description: Service account for lifecycle testing → Service account for lifecycle testing - Updated
[INFO] Updating Roles: ORG_MEMBER → ORG_MEMBER,ORG_GROUP_CREATOR
[INFO] Waiting for stack update to complete...
[SUCCESS] Stack update completed: UPDATE_COMPLETE
[SUCCESS] UPDATE COMPLETE
```

### Stage 3: DELETE

```
[INFO] ==========================================
[INFO] STAGE 3: DELETE
[INFO] ==========================================
[INFO] Deleting CloudFormation stack: org-service-account-test-1234567890
[INFO] Waiting for stack deletion to complete...
[SUCCESS] Stack deletion completed: DELETE_COMPLETE
[SUCCESS] Service account deleted from Atlas
[SUCCESS] DELETE COMPLETE
[SUCCESS] ==========================================
[SUCCESS] Lifecycle test completed successfully!
[SUCCESS] ==========================================
```

## Features

- ✅ **Automatic cleanup**: Failed stacks (ROLLBACK_COMPLETE, CREATE_FAILED, DELETE_FAILED) are automatically cleaned up
- ✅ **Real-time monitoring**: Progress updates with stack status
- ✅ **Atlas CLI validation**: Validates resource existence and configuration after each stage
- ✅ **Error reporting**: Detailed stack events on failure
- ✅ **Colored output**: Easy-to-read status messages
- ✅ **Graceful failures**: Proper error handling and cleanup

## Updatable Fields

The following fields can be updated in Stage 2:
- **Name**: Service account name
- **Description**: Service account description
- **Roles**: List of organization-level roles

**Note**: The following fields are **create-only** and cannot be updated:
- `OrgId`: Organization ID (create-only)
- `Profile`: AWS Secrets Manager profile (create-only)
- `SecretExpiresAfterHours`: Secret expiration time (create-only)

## Troubleshooting

### Stack Creation Fails

1. **Check CloudWatch Logs**:
   ```bash
   aws logs describe-log-groups \
     --log-group-name-prefix "/aws/lambda/mongodb-atlas-org-service-account" \
     --region eu-west-1
   ```

2. **Verify Credentials**:
   - Ensure AWS Secrets Manager secret exists: `cfn/atlas/profile/default`
   - Verify secret contains valid Atlas API keys
   - Check IAM execution role has `secretsmanager:GetSecretValue` permission

3. **Verify OrgId**:
   - Ensure OrgId is a valid 24-character hexadecimal string
   - Verify you have access to the organization

### Update Fails with AccessDenied

- Ensure IAM policy includes `cloudformation:UpdateStack` permission
- The minimal IAM policy (`cfn-test-and-submit-policy-minimal.json`) includes all required permissions

### Atlas CLI Validation Fails

- Ensure Atlas CLI is installed and configured: `atlas auth login`
- Verify API keys have appropriate permissions
- Check organization ID is correct

## Cleanup

The script automatically cleans up:
- Failed stacks (ROLLBACK_COMPLETE, CREATE_FAILED, DELETE_FAILED)
- Output files after successful DELETE

**Manual cleanup** (if needed):

```bash
# List stacks
aws cloudformation list-stacks \
  --region eu-west-1 \
  --stack-status-filter CREATE_COMPLETE UPDATE_COMPLETE \
  --query "StackSummaries[?contains(StackName, 'org-service-account')].StackName" \
  --output text

# Delete specific stack
aws cloudformation delete-stack \
  --stack-name <stack-name> \
  --region eu-west-1
```

## Additional Resources

- [MongoDB Atlas Service Accounts Documentation](https://www.mongodb.com/docs/atlas/security/service-accounts/)
- [Atlas API Documentation](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/group/endpoint-service-accounts)
- [CloudFormation Resource Documentation](../../docs/README.md)
- [Example Template](../../../examples/org-service-account/org-service-account.json)


