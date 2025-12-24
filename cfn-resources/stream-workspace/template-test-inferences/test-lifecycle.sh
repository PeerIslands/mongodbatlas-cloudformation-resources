#!/usr/bin/env bash
# test-lifecycle.sh - Complete lifecycle test with user confirmation at each stage

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "${SCRIPT_DIR}"

# Configuration - can be overridden by environment variables or command line args
PROJECT_ID="${MONGODB_ATLAS_PROJECT_ID:-}"
WORKSPACE_NAME="${WORKSPACE_NAME:-}"
STACK_NAME="${STACK_NAME:-}"
AWS_REGION="${AWS_REGION:-eu-west-1}"
PROFILE="${PROFILE:-default}"
OUTPUT_FILE="${OUTPUT_FILE:-output.json}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO] $1${NC}"
}

log_success() {
    echo -e "${GREEN}[SUCCESS] $1${NC}"
}

log_error() {
    echo -e "${RED}[ERROR] $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}[WARNING] $1${NC}"
}

log_stage() {
    echo -e "${CYAN}============================================${NC}"
    echo -e "${CYAN}$1${NC}"
    echo -e "${CYAN}============================================${NC}"
}

# Wait for user confirmation to proceed to DELETE
# NOTE: UPDATE is not supported for StreamWorkspace as StreamConfig is create-only
wait_for_delete() {
    echo ""
    log_warning "NOTE: StreamWorkspace does not support updates (StreamConfig is create-only)"
    log_info "What would you like to do next?"
    while true; do
        read -p "Type 'delete' to delete the workspace, or 'exit' to quit: " response
        case "${response}" in
            [Dd][Ee][Ll][Ee][Tt][Ee]|delete|DELETE)
                log_success "Proceeding to DELETE..."
                echo ""
                return 0
                ;;
            [Ee][Xx][Ii][Tt]|exit|EXIT)
                log_warning "Exiting as requested."
                exit 0
                ;;
            *)
                log_warning "Please type 'delete' or 'exit'"
                ;;
        esac
    done
}

# Check prerequisites
check_prerequisites() {
    echo "=== Prerequisites Check ==="
    
    log_info "Checking AWS credentials..."
    if ! aws sts get-caller-identity > /dev/null 2>&1; then
        log_error "AWS credentials not configured"
        log_error "Please configure AWS credentials using:"
        log_error "  - aws configure"
        log_error "  - Environment variables: AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY"
        exit 1
    fi
    AWS_ACCOUNT=$(aws sts get-caller-identity --query Account --output text)
    log_success "AWS Account: ${AWS_ACCOUNT}"

    log_info "Checking Atlas CLI..."
    if ! command -v atlas &> /dev/null; then
        log_error "Atlas CLI not found. Please install it first."
        log_error "Install: https://www.mongodb.com/docs/atlas/cli/stable/atlas-cli-install/"
        exit 1
    fi
    log_success "Atlas CLI found"
    
    # Verify Atlas credentials
    log_info "Verifying Atlas credentials..."
    if atlas projects list --output json > /dev/null 2>&1; then
        log_success "Atlas credentials verified"
    else
        log_warning "Atlas CLI credentials not configured or invalid"
        log_warning "Please configure using:"
        log_warning "  - Environment variables: MONGODB_ATLAS_PUBLIC_API_KEY, MONGODB_ATLAS_PRIVATE_API_KEY"
        log_warning "  - atlas config (atlas config set public_api_key, private_api_key)"
        log_warning "Continuing anyway, but Atlas validation may fail..."
    fi

    log_info "Checking template file..."
    if [ ! -f "stream-workspace-complex.json" ]; then
        log_error "Template file not found: stream-workspace-complex.json"
        exit 1
    fi
    log_success "Template file found"
    echo ""
}

# Get user inputs
get_user_inputs() {
    echo "=== Configuration Input ==="
    
    if [ -z "${PROJECT_ID}" ]; then
        read -p "Enter MongoDB Atlas Project ID: " PROJECT_ID
    fi
    if [ -z "${PROJECT_ID}" ]; then
        log_error "Project ID is required"
        exit 1
    fi
    
    if [ -z "${WORKSPACE_NAME}" ]; then
        read -p "Enter Workspace Name (or press Enter for auto-generated): " WORKSPACE_NAME
    fi
    if [ -z "${WORKSPACE_NAME}" ]; then
        WORKSPACE_NAME="stream-workspace-$(date +%s)"
        log_info "Auto-generated workspace name: ${WORKSPACE_NAME}"
    fi
    
    if [ -z "${STACK_NAME}" ]; then
        read -p "Enter CloudFormation Stack Name (or press Enter for auto-generated): " STACK_NAME
    fi
    if [ -z "${STACK_NAME}" ]; then
        STACK_NAME="stream-workspace-$(date +%s)"
        log_info "Auto-generated stack name: ${STACK_NAME}"
    fi
    
    read -p "Enter AWS Region [${AWS_REGION}]: " input_region
    if [ -n "${input_region}" ]; then
        AWS_REGION="${input_region}"
    fi
    
    read -p "Enter Profile [${PROFILE}]: " input_profile
    if [ -n "${input_profile}" ]; then
        PROFILE="${input_profile}"
    fi
    
    echo ""
    echo "=== Configuration Summary ==="
    echo "Project ID: ${PROJECT_ID}"
    echo "Workspace Name: ${WORKSPACE_NAME}"
    echo "Stack Name: ${STACK_NAME}"
    echo "Region: ${AWS_REGION}"
    echo "Profile: ${PROFILE}"
    echo ""
}

# Cleanup function
cleanup_failed_stacks() {
    log_info "Checking for existing failed stacks..."
    EXISTING_STACKS=$(aws cloudformation list-stacks \
        --region "${AWS_REGION}" \
        --stack-status-filter ROLLBACK_COMPLETE CREATE_FAILED DELETE_FAILED \
        --query "StackSummaries[?contains(StackName, '${STACK_NAME}')].StackName" \
        --output text 2>/dev/null || echo "")
    
    if [ -n "${EXISTING_STACKS}" ]; then
        log_warning "Found existing failed stacks. Cleaning up..."
        for stack in ${EXISTING_STACKS}; do
            log_info "  Deleting stack: ${stack}"
            aws cloudformation delete-stack --stack-name "${stack}" --region "${AWS_REGION}" > /dev/null 2>&1 || true
        done
        sleep 10
    fi
}

# Wait for stack operation
wait_for_stack() {
    local operation=$1  # CREATE, UPDATE, or DELETE
    local max_attempts=60
    local attempt=1
    
    while [ ${attempt} -le ${max_attempts} ]; do
        sleep 10
        STATUS=$(aws cloudformation describe-stacks \
            --stack-name "${STACK_NAME}" \
            --region "${AWS_REGION}" \
            --query 'Stacks[0].StackStatus' \
            --output text 2>/dev/null || echo "NOT_FOUND")
        
        echo "[${attempt}/${max_attempts}] Status: ${STATUS}"
        
        case "${operation}" in
            CREATE)
                if [[ "${STATUS}" == "CREATE_COMPLETE" ]]; then
                    log_success "Stack created successfully!"
                    return 0
                elif [[ "${STATUS}" == "ROLLBACK_COMPLETE" ]] || [[ "${STATUS}" == "CREATE_FAILED" ]]; then
                    log_error "Stack creation failed. Status: ${STATUS}"
                    aws cloudformation describe-stack-events \
                        --stack-name "${STACK_NAME}" \
                        --region "${AWS_REGION}" \
                        --max-items 5 \
                        --query 'StackEvents[?ResourceStatus==`CREATE_FAILED`]' \
                        --output table
                    return 1
                fi
                ;;
            UPDATE)
                if [[ "${STATUS}" == "UPDATE_COMPLETE" ]]; then
                    log_success "Stack updated successfully!"
                    return 0
                elif [[ "${STATUS}" == "UPDATE_ROLLBACK_COMPLETE" ]] || [[ "${STATUS}" == "UPDATE_FAILED" ]]; then
                    log_error "Stack update failed. Status: ${STATUS}"
                    aws cloudformation describe-stack-events \
                        --stack-name "${STACK_NAME}" \
                        --region "${AWS_REGION}" \
                        --max-items 5 \
                        --query 'StackEvents[?ResourceStatus==`UPDATE_FAILED`]' \
                        --output table
                    return 1
                fi
                ;;
            DELETE)
                if [[ "${STATUS}" == "DELETE_COMPLETE" ]] || [[ "${STATUS}" == "NOT_FOUND" ]]; then
                    log_success "Stack deleted successfully!"
                    return 0
                elif [[ "${STATUS}" == "DELETE_FAILED" ]]; then
                    log_error "Stack deletion failed. Status: ${STATUS}"
                    return 1
                fi
                ;;
        esac
        
        attempt=$((attempt + 1))
    done
    
    log_error "Operation timed out"
    return 1
}

# Validate workspace in Atlas
validate_workspace_atlas() {
    local expected_name=$1
    local expected_tier=${2:-}
    local expected_max_tier=${3:-}
    local expected_region=${4:-}
    
    log_info "Validating workspace in Atlas..." >&2
    
    # Check if workspace exists
    if ! atlas streams instances describe "${expected_name}" --projectId "${PROJECT_ID}" > /dev/null 2>&1; then
        log_error "Workspace '${expected_name}' not found in Atlas" >&2
        return 1
    fi
    log_success "Workspace '${expected_name}' exists in Atlas" >&2
    
    # Get workspace details
    WORKSPACE_DETAILS=$(atlas streams instances describe "${expected_name}" --projectId "${PROJECT_ID}" --output json 2>/dev/null)
    
    if [ -z "${WORKSPACE_DETAILS}" ] || ! echo "${WORKSPACE_DETAILS}" | jq . > /dev/null 2>&1; then
        log_error "Failed to get workspace details or invalid JSON" >&2
        return 1
    fi
    
    if [ -n "${expected_tier}" ]; then
        ACTUAL_TIER=$(echo "${WORKSPACE_DETAILS}" | jq -r '.streamConfig.tier // empty')
        if [ "${ACTUAL_TIER}" == "${expected_tier}" ]; then
            log_success "Tier matches: ${expected_tier}" >&2
        else
            log_error "Tier mismatch. Expected: ${expected_tier}, Actual: ${ACTUAL_TIER}" >&2
            return 1
        fi
    fi
    
    if [ -n "${expected_max_tier}" ]; then
        ACTUAL_MAX_TIER=$(echo "${WORKSPACE_DETAILS}" | jq -r '.streamConfig.maxTierSize // empty')
        if [ "${ACTUAL_MAX_TIER}" == "${expected_max_tier}" ]; then
            log_success "MaxTierSize matches: ${expected_max_tier}" >&2
        else
            log_warning "MaxTierSize mismatch. Expected: ${expected_max_tier}, Actual: ${ACTUAL_MAX_TIER}" >&2
        fi
    fi
    
    if [ -n "${expected_region}" ]; then
        ACTUAL_REGION=$(echo "${WORKSPACE_DETAILS}" | jq -r '.dataProcessRegion.region // empty')
        if [ "${ACTUAL_REGION}" == "${expected_region}" ]; then
            log_success "Region matches: ${expected_region}" >&2
        else
            log_warning "Region mismatch. Expected: ${expected_region}, Actual: ${ACTUAL_REGION}" >&2
        fi
    fi
    
    # Display workspace details
    echo "" >&2
    log_info "Workspace Details from Atlas:" >&2
    echo "${WORKSPACE_DETAILS}" | jq '{
        Name: .name,
        Id: ._id,
        DataProcessRegion: .dataProcessRegion,
        StreamConfig: .streamConfig,
        Hostnames: .hostnames
    }' >&2
    echo "" >&2
    
    # Return workspace details as JSON (to stdout only)
    echo "${WORKSPACE_DETAILS}"
}

# Verify workspace deleted
verify_workspace_deleted() {
    local workspace_name=$1
    
    log_info "Verifying workspace deletion in Atlas..."
    
    if atlas streams instances describe "${workspace_name}" --projectId "${PROJECT_ID}" > /dev/null 2>&1; then
        log_error "Workspace '${workspace_name}' still exists in Atlas"
        return 1
    fi
    
    log_success "Workspace '${workspace_name}' successfully deleted from Atlas"
    return 0
}

# Save resource info to output.json
save_resource_info() {
    local workspace_details=$1
    local stack_outputs=$2
    
    log_info "Saving resource information to ${OUTPUT_FILE}..."
    
    # Get stack details
    STACK_DETAILS=$(aws cloudformation describe-stacks \
        --stack-name "${STACK_NAME}" \
        --region "${AWS_REGION}" \
        --output json)
    
    # Validate JSON inputs
    if ! echo "${STACK_DETAILS}" | jq . > /dev/null 2>&1; then
        log_error "Invalid JSON in stack details"
        return 1
    fi
    
    if ! echo "${workspace_details}" | jq . > /dev/null 2>&1; then
        log_warning "Invalid JSON in workspace details, using empty object"
        workspace_details="{}"
    fi
    
    if ! echo "${stack_outputs}" | jq . > /dev/null 2>&1; then
        log_warning "Invalid JSON in stack outputs, using empty array"
        stack_outputs="[]"
    fi
    
    # Combine all information
    OUTPUT_DATA=$(jq -n \
        --arg stackName "${STACK_NAME}" \
        --arg awsRegion "${AWS_REGION}" \
        --arg projectId "${PROJECT_ID}" \
        --arg workspaceName "${WORKSPACE_NAME}" \
        --arg profile "${PROFILE}" \
        --argjson stackDetails "${STACK_DETAILS}" \
        --argjson workspaceDetails "${workspace_details}" \
        --argjson stackOutputs "${stack_outputs}" \
        '{
            stackName: $stackName,
            awsRegion: $awsRegion,
            projectId: $projectId,
            workspaceName: $workspaceName,
            profile: $profile,
            createdAt: (now | strftime("%Y-%m-%dT%H:%M:%SZ")),
            stackDetails: $stackDetails,
            workspaceDetails: $workspaceDetails,
            stackOutputs: $stackOutputs
        }')
    
    echo "${OUTPUT_DATA}" > "${OUTPUT_FILE}"
    log_success "Resource information saved to ${OUTPUT_FILE}"
}

# ============================================
# MAIN SCRIPT
# ============================================

check_prerequisites
get_user_inputs

# ============================================
# STEP 1: CREATE WORKSPACE
# ============================================
log_stage "CREATE START"
echo ""

cleanup_failed_stacks

log_info "Creating CloudFormation stack with initial configuration..."
log_info "Initial Config: Tier=SP2, MaxTierSize=SP50, Region=VIRGINIA_USA"
log_info "Testing MaxTierSize >= Tier validation: SP2 <= SP50 (valid)"

aws cloudformation create-stack \
    --stack-name "${STACK_NAME}" \
    --template-body file://stream-workspace-complex.json \
    --parameters \
        ParameterKey=ProjectId,ParameterValue="${PROJECT_ID}" \
        ParameterKey=WorkspaceName,ParameterValue="${WORKSPACE_NAME}" \
        ParameterKey=CloudProvider,ParameterValue=AWS \
        ParameterKey=Region,ParameterValue=VIRGINIA_USA \
        ParameterKey=Tier,ParameterValue=SP2 \
        ParameterKey=MaxTierSize,ParameterValue=SP50 \
        ParameterKey=Profile,ParameterValue="${PROFILE}" \
    --capabilities CAPABILITY_IAM \
    --region "${AWS_REGION}" \
    --output json > /dev/null

log_success "Stack creation initiated"

if ! wait_for_stack "CREATE"; then
    log_error "Stack creation failed"
    exit 1
fi

# Get stack outputs
STACK_OUTPUTS=$(aws cloudformation describe-stacks \
    --stack-name "${STACK_NAME}" \
    --region "${AWS_REGION}" \
    --query 'Stacks[0].Outputs' \
    --output json)

# Display stack outputs
echo ""
log_info "Stack Outputs:"
echo "${STACK_OUTPUTS}" | jq '.'

# Validate with Atlas CLI
echo ""
WORKSPACE_DETAILS_CREATE=$(validate_workspace_atlas "${WORKSPACE_NAME}" "SP2" "SP50" "VIRGINIA_USA")
if [ $? -ne 0 ]; then
    log_error "CREATE validation failed"
    exit 1
fi

# Save resource info
save_resource_info "${WORKSPACE_DETAILS_CREATE}" "${STACK_OUTPUTS}"

log_stage "CREATE COMPLETE"
echo ""

# NOTE: UPDATE is not supported for StreamWorkspace
# StreamConfig is create-only, so updates require resource replacement
# which is not possible with custom-named resources in CloudFormation
log_warning "NOTE: StreamWorkspace does not support updates"
log_info "All main properties (WorkspaceName, ProjectId, Profile, StreamConfig) are create-only"
log_info "To change configuration, you must delete and recreate the workspace"
echo ""

# Wait for user to proceed to DELETE
wait_for_delete

# ============================================
# STEP 2: DELETE WORKSPACE
# ============================================
log_stage "DELETE START"
echo ""

log_info "Deleting CloudFormation stack..."

aws cloudformation delete-stack \
    --stack-name "${STACK_NAME}" \
    --region "${AWS_REGION}"

log_success "Stack deletion initiated"

if ! wait_for_stack "DELETE"; then
    log_error "Stack deletion failed"
    exit 1
fi

# Verify deletion with Atlas CLI
echo ""
if ! verify_workspace_deleted "${WORKSPACE_NAME}"; then
    log_error "DELETE validation failed"
    exit 1
fi

# Clean up output file
log_info "Cleaning up output file..."
rm -f "${OUTPUT_FILE}"
log_success "Output file removed"

log_stage "DELETE COMPLETE"
echo ""

# ============================================
# FINAL SUMMARY
# ============================================
echo "============================================"
echo "LIFECYCLE TEST COMPLETE"
echo "============================================"
echo ""
log_success "CREATE: Workspace created successfully (Tier=SP2, MaxTierSize=SP50)"
log_info "UPDATE: Not supported (StreamConfig is create-only)"
log_success "DELETE: Workspace deleted successfully"
echo ""
log_success "Lifecycle test finished successfully!"
echo ""
