# StreamConnection: Terraform vs CloudFormation Gap Analysis

## Overview
This document identifies the gaps between the Terraform `mongodbatlas_stream_connection` resource and the CloudFormation `MongoDB::Atlas::StreamConnection` resource.

## Current Status
- **CloudFormation Resource**: Exists but incomplete
- **Terraform Resource**: Fully featured
- **Gap Level**: **HIGH** - Missing many critical features

---

## Missing Properties

### 1. WorkspaceName Support ❌
**Terraform**: Supports both `workspace_name` (new) and `instance_name` (deprecated)
**CloudFormation**: Only has `InstanceName` (deprecated)

**Impact**: HIGH - Should support `WorkspaceName` as primary identifier

### 2. ClusterProjectId ❌
**Terraform**: `cluster_project_id` - For cross-project cluster connections
**CloudFormation**: Missing

**Impact**: MEDIUM - Needed for cross-project connections

### 3. Enhanced Authentication (OAuth) ❌
**Terraform** supports:
- `authentication.method`
- `authentication.token_endpoint_url`
- `authentication.client_id`
- `authentication.client_secret` (sensitive)
- `authentication.scope`
- `authentication.sasl_oauthbearer_extensions`

**CloudFormation** only has:
- `Authentication.Mechanism`
- `Authentication.Username`
- `Authentication.Password`

**Impact**: HIGH - OAuth authentication is critical for many Kafka setups

### 4. Networking Block ❌
**Terraform**: 
```hcl
networking {
  access {
    type          = "PRIVATE_ENDPOINT"
    connection_id = "..."
  }
}
```

**CloudFormation**: Missing entirely

**Impact**: HIGH - Required for private endpoint connections

### 5. AWS Lambda Support ❌
**Terraform**: 
```hcl
aws {
  role_arn = "arn:aws:iam::..."
}
```

**CloudFormation**: Missing entirely

**Impact**: MEDIUM - Needed for AWS Lambda connection type

### 6. HTTPS Connection Support ❌
**Terraform**:
- `url` - HTTPS endpoint URL
- `headers` - Map of HTTP headers

**CloudFormation**: Missing entirely

**Impact**: MEDIUM - Needed for HTTPS connection type

### 7. Connection Types ❌
**Terraform** supports: `Kafka`, `Cluster`, `Sample`, `AWSLambda`, `Https`
**CloudFormation** supports: `Kafka`, `Cluster`, `Sample` only

**Impact**: HIGH - Missing two connection types

---

## Property Comparison Table

| Property | Terraform | CloudFormation | Status | Priority |
|----------|-----------|----------------|--------|----------|
| `project_id` / `ProjectId` | ✅ Required | ✅ Required | ✅ MATCH | - |
| `workspace_name` / `WorkspaceName` | ✅ Optional (new) | ❌ Missing | ❌ GAP | HIGH |
| `instance_name` / `InstanceName` | ⚠️ Optional (deprecated) | ✅ Required | ⚠️ PARTIAL | HIGH |
| `connection_name` / `ConnectionName` | ✅ Required | ✅ Required | ✅ MATCH | - |
| `type` / `Type` | ✅ Required | ✅ Required | ⚠️ PARTIAL | HIGH |
| `cluster_name` / `ClusterName` | ✅ Optional | ✅ Optional | ✅ MATCH | - |
| `cluster_project_id` / `ClusterProjectId` | ✅ Optional | ❌ Missing | ❌ GAP | MEDIUM |
| `db_role_to_execute` / `DbRoleToExecute` | ✅ Optional | ✅ Optional | ✅ MATCH | - |
| `authentication.mechanism` | ✅ Optional | ✅ Optional | ✅ MATCH | - |
| `authentication.username` | ✅ Optional | ✅ Optional | ✅ MATCH | - |
| `authentication.password` | ✅ Optional (sensitive) | ✅ Optional (sensitive) | ✅ MATCH | - |
| `authentication.method` | ✅ Optional | ❌ Missing | ❌ GAP | HIGH |
| `authentication.token_endpoint_url` | ✅ Optional | ❌ Missing | ❌ GAP | HIGH |
| `authentication.client_id` | ✅ Optional | ❌ Missing | ❌ GAP | HIGH |
| `authentication.client_secret` | ✅ Optional (sensitive) | ❌ Missing | ❌ GAP | HIGH |
| `authentication.scope` | ✅ Optional | ❌ Missing | ❌ GAP | HIGH |
| `authentication.sasl_oauthbearer_extensions` | ✅ Optional | ❌ Missing | ❌ GAP | HIGH |
| `bootstrap_servers` / `BootstrapServers` | ✅ Optional | ✅ Optional | ✅ MATCH | - |
| `config` / `Config` | ✅ Optional | ✅ Optional | ✅ MATCH | - |
| `security.broker_public_certificate` | ✅ Optional | ✅ Optional | ✅ MATCH | - |
| `security.protocol` | ✅ Optional | ✅ Optional | ✅ MATCH | - |
| `networking.access.type` | ✅ Optional | ❌ Missing | ❌ GAP | HIGH |
| `networking.access.connection_id` | ✅ Optional | ❌ Missing | ❌ GAP | HIGH |
| `aws.role_arn` | ✅ Optional | ❌ Missing | ❌ GAP | MEDIUM |
| `url` | ✅ Optional | ❌ Missing | ❌ GAP | MEDIUM |
| `headers` | ✅ Optional | ❌ Missing | ❌ GAP | MEDIUM |

---

## CRUD Operations Comparison

### CREATE Operation
**Terraform**: 
- Uses `workspace_name` or `instance_name` (backward compatibility)
- Supports all connection types
- Handles OAuth authentication
- Supports networking, AWS, HTTPS

**CloudFormation**: 
- Only uses `InstanceName`
- Limited connection types
- Basic authentication only
- Missing networking, AWS, HTTPS

**Status**: ⚠️ **PARTIAL** - Missing many features

### READ Operation
**Terraform**: Returns all fields including networking, AWS, headers
**CloudFormation**: Returns limited fields

**Status**: ⚠️ **PARTIAL** - Missing read support for new fields

### UPDATE Operation
**Terraform**: Can update all fields (except create-only)
**CloudFormation**: Can update basic fields only

**Status**: ⚠️ **PARTIAL** - Missing update support for new fields

### DELETE Operation
**Terraform**: Uses workspace_name or instance_name
**CloudFormation**: Uses instance_name only

**Status**: ⚠️ **PARTIAL** - Should support WorkspaceName

---

## Implementation Notes

### API Version
- **Terraform**: Uses Atlas SDK v20250312010
- **CloudFormation**: Uses Atlas SDK v20231115014 (older version)

**Action**: Update to latest SDK version to get new features

### Workspace vs Instance
- Terraform supports both `workspace_name` and `instance_name` for backward compatibility
- CloudFormation should add `WorkspaceName` support while keeping `InstanceName` for backward compatibility
- Both should map to the same API endpoint (workspace_name is the new name)

### Authentication Handling
- Password and ClientSecret are sensitive fields (write-only in CloudFormation)
- Need to preserve sensitive values during updates (don't overwrite if not provided)

### Connection Type Logic
- Different connection types require different fields:
  - **Cluster**: `cluster_name`, `cluster_project_id`, `db_role_to_execute`
  - **Kafka**: `bootstrap_servers`, `authentication`, `security`, `config`, `networking`
  - **AWSLambda**: `aws.role_arn`
  - **Https**: `url`, `headers`
  - **Sample**: Basic fields only

---

## Priority Implementation Order

### Phase 1: Critical (HIGH Priority)
1. ✅ Add `WorkspaceName` support (keep `InstanceName` for backward compatibility)
2. ✅ Add OAuth authentication fields
3. ✅ Add `Networking` block
4. ✅ Add `AWSLambda` and `Https` connection types

### Phase 2: Important (MEDIUM Priority)
5. ✅ Add `ClusterProjectId`
6. ✅ Add `AWS` block for Lambda connections
7. ✅ Add `URL` and `Headers` for HTTPS connections

### Phase 3: Enhancement
8. ✅ Update SDK version
9. ✅ Improve error handling
10. ✅ Add comprehensive tests

---

## Summary

**Total Gaps**: 15 missing properties/features
**Critical Gaps**: 7 (WorkspaceName, OAuth auth, Networking, Connection types)
**Medium Gaps**: 3 (ClusterProjectId, AWS, HTTPS)
**Low Gaps**: 0

**Estimated Effort**: 
- Phase 1: 2-3 days
- Phase 2: 1-2 days
- Phase 3: 1 day
- Testing: 2 days
- **Total**: 6-8 days

**Status**: ⚠️ **NEEDS SIGNIFICANT UPDATES** to achieve feature parity




