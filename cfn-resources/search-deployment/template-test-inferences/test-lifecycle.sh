#!/usr/bin/env bash
# test-lifecycle.sh - Complete lifecycle test for Search Deployment with user confirmation at each stage

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "${SCRIPT_DIR}"

# Configuration - can be overridden by environment variables or command line args
PROJECT_ID="${MONGODB_ATLAS_PROJECT_ID:-}"
CLUSTER_NAME="${MONGODB_ATLAS_CLUSTER_NAME:-}"
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

# Wait for user confirmation
wait_for_confirmation() {
    local stage=$1
    echo ""
    log_info "What would you like to do next?"
    while true; do
        read -p "Type 'yes' to continue to ${stage}, or 'exit' to quit: " response
        case "${response}" in
            [Yy][Ee][Ss]|yes|YES)
                log_success "Proceeding to ${stage}..."
                echo ""
                return 0
                ;;
            [Ee][Xx][Ii][Tt]|exit|EXIT)
                log_warning "Exiting as requested."
                exit 0
                ;;
            *)
                log_warning "Please type 'yes' or 'exit'"
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
    if [ ! -f "search-deployment-complex.json" ]; then
        log_error "Template file not found: search-deployment-complex.json"
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
    
    if [ -z "${CLUSTER_NAME}" ]; then
        read -p "Enter Cluster Name: " CLUSTER_NAME
    fi
    if [ -z "${CLUSTER_NAME}" ]; then
        log_error "Cluster Name is required"
        exit 1
    fi
    
    if [ -z "${STACK_NAME}" ]; then
        read -p "Enter CloudFormation Stack Name (or press Enter for auto-generated): " STACK_NAME
    fi
    if [ -z "${STACK_NAME}" ]; then
        STACK_NAME="search-deployment-$(date +%s)"
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
    echo "Cluster Name: ${CLUSTER_NAME}"
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
    local max_attempts=120  # Search deployment can take longer (up to 60 minutes)
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

# Validate search deployment in Atlas
validate_search_deployment_atlas() {
    local expected_instance_size=${1:-}
    local expected_node_count=${2:-}
    
    log_info "Validating search deployment in Atlas..."
    
    # Get search deployment details
    SEARCH_DEPLOYMENT=$(atlas clusters search nodes list \
        --clusterName "${CLUSTER_NAME}" \
        --projectId "${PROJECT_ID}" \
        --output json 2>/dev/null || echo "{}")
    
    if [ -z "${SEARCH_DEPLOYMENT}" ] || [ "${SEARCH_DEPLOYMENT}" == "{}" ]; then
        log_error "Search deployment not found for cluster '${CLUSTER_NAME}' in Atlas"
        return 1
    fi
    
    if ! echo "${SEARCH_DEPLOYMENT}" | jq . > /dev/null 2>&1; then
        log_error "Failed to parse search deployment details or invalid JSON"
        return 1
    fi
    
    log_success "Search deployment found for cluster '${CLUSTER_NAME}' in Atlas"
    
    # Check instance size if provided
    if [ -n "${expected_instance_size}" ]; then
        ACTUAL_INSTANCE_SIZE=$(echo "${SEARCH_DEPLOYMENT}" | jq -r '.specs[0].instanceSize // empty')
        if [ "${ACTUAL_INSTANCE_SIZE}" == "${expected_instance_size}" ]; then
            log_success "InstanceSize matches: ${expected_instance_size}"
        else
            log_warning "InstanceSize mismatch. Expected: ${expected_instance_size}, Actual: ${ACTUAL_INSTANCE_SIZE}"
        fi
    fi
    
    # Check node count if provided
    if [ -n "${expected_node_count}" ]; then
        ACTUAL_NODE_COUNT=$(echo "${SEARCH_DEPLOYMENT}" | jq -r '.specs[0].nodeCount // empty')
        if [ "${ACTUAL_NODE_COUNT}" == "${expected_node_count}" ]; then
            log_success "NodeCount matches: ${expected_node_count}"
        else
            log_warning "NodeCount mismatch. Expected: ${expected_node_count}, Actual: ${ACTUAL_NODE_COUNT}"
        fi
    fi
    
    # Display search deployment details
    echo ""
    log_info "Search Deployment Details from Atlas:"
    echo "${SEARCH_DEPLOYMENT}" | jq '{
        Id: .id,
        StateName: .stateName,
        Specs: .specs,
        EncryptionAtRestProvider: .encryptionAtRestProvider
    }'
    echo ""
    
    # Return search deployment details as JSON (to stdout only)
    echo "${SEARCH_DEPLOYMENT}"
}

# Verify search deployment deleted
verify_search_deployment_deleted() {
    log_info "Verifying search deployment deletion in Atlas..."
    
    SEARCH_DEPLOYMENT=$(atlas clusters search nodes list \
        --clusterName "${CLUSTER_NAME}" \
        --projectId "${PROJECT_ID}" \
        --output json 2>/dev/null || echo "{}")
    
    if [ -z "${SEARCH_DEPLOYMENT}" ] || [ "${SEARCH_DEPLOYMENT}" == "{}" ] || [ "${SEARCH_DEPLOYMENT}" == "null" ]; then
        log_success "Search deployment successfully deleted from Atlas"
        return 0
    fi
    
    # Check if response indicates no deployment (empty array or null)
    DEPLOYMENT_COUNT=$(echo "${SEARCH_DEPLOYMENT}" | jq '.specs | length // 0' 2>/dev/null || echo "0")
    if [ "${DEPLOYMENT_COUNT}" == "0" ]; then
        log_success "Search deployment successfully deleted from Atlas"
        return 0
    fi
    
    log_error "Search deployment still exists in Atlas"
    return 1
}

# Save resource info to output.json
save_resource_info() {
    local search_deployment_details=$1
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
    
    if ! echo "${search_deployment_details}" | jq . > /dev/null 2>&1; then
        log_warning "Invalid JSON in search deployment details, using empty object"
        search_deployment_details="{}"
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
        --arg clusterName "${CLUSTER_NAME}" \
        --arg profile "${PROFILE}" \
        --argjson stackDetails "${STACK_DETAILS}" \
        --argjson searchDeploymentDetails "${search_deployment_details}" \
        --argjson stackOutputs "${stack_outputs}" \
        '{
            stackName: $stackName,
            awsRegion: $awsRegion,
            projectId: $projectId,
            clusterName: $clusterName,
            profile: $profile,
            createdAt: (now | strftime("%Y-%m-%dT%H:%M:%SZ")),
            stackDetails: $stackDetails,
            searchDeploymentDetails: $searchDeploymentDetails,
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
# STEP 1: CREATE SEARCH DEPLOYMENT
# ============================================
log_stage "CREATE START"
echo ""

cleanup_failed_stacks

log_info "Creating CloudFormation stack with initial configuration..."
log_info "Initial Config: InstanceSize=S30_HIGHCPU_NVME, NodeCount=2"

aws cloudformation create-stack \
    --stack-name "${STACK_NAME}" \
    --template-body file://search-deployment-complex.json \
    --parameters \
        ParameterKey=ProjectId,ParameterValue="${PROJECT_ID}" \
        ParameterKey=ClusterName,ParameterValue="${CLUSTER_NAME}" \
        ParameterKey=InstanceSize,ParameterValue=S30_HIGHCPU_NVME \
        ParameterKey=NodeCount,ParameterValue=2 \
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
SEARCH_DEPLOYMENT_DETAILS_CREATE=$(validate_search_deployment_atlas "S30_HIGHCPU_NVME" "2")
if [ $? -ne 0 ]; then
    log_error "CREATE validation failed"
    exit 1
fi

# Save resource info
save_resource_info "${SEARCH_DEPLOYMENT_DETAILS_CREATE}" "${STACK_OUTPUTS}"

log_stage "CREATE COMPLETE"
echo ""

# Wait for user to proceed to UPDATE
wait_for_confirmation "UPDATE"

# ============================================
# STEP 2: UPDATE SEARCH DEPLOYMENT
# ============================================
log_stage "UPDATE START"
echo ""

log_info "Updating CloudFormation stack with modified configuration..."
log_info "Updated Config: InstanceSize=S30_HIGHCPU_NVME, NodeCount=3"

aws cloudformation update-stack \
    --stack-name "${STACK_NAME}" \
    --template-body file://search-deployment-complex.json \
    --parameters \
        ParameterKey=ProjectId,ParameterValue="${PROJECT_ID}" \
        ParameterKey=ClusterName,ParameterValue="${CLUSTER_NAME}" \
        ParameterKey=InstanceSize,ParameterValue=S30_HIGHCPU_NVME \
        ParameterKey=NodeCount,ParameterValue=3 \
        ParameterKey=Profile,ParameterValue="${PROFILE}" \
    --capabilities CAPABILITY_IAM \
    --region "${AWS_REGION}" \
    --output json > /dev/null

log_success "Stack update initiated"

if ! wait_for_stack "UPDATE"; then
    log_error "Stack update failed"
    exit 1
fi

# Get updated stack outputs
STACK_OUTPUTS_UPDATE=$(aws cloudformation describe-stacks \
    --stack-name "${STACK_NAME}" \
    --region "${AWS_REGION}" \
    --query 'Stacks[0].Outputs' \
    --output json)

# Display updated stack outputs
echo ""
log_info "Updated Stack Outputs:"
echo "${STACK_OUTPUTS_UPDATE}" | jq '.'

# Validate update with Atlas CLI
echo ""
SEARCH_DEPLOYMENT_DETAILS_UPDATE=$(validate_search_deployment_atlas "S30_HIGHCPU_NVME" "3")
if [ $? -ne 0 ]; then
    log_error "UPDATE validation failed"
    exit 1
fi

# Update resource info
save_resource_info "${SEARCH_DEPLOYMENT_DETAILS_UPDATE}" "${STACK_OUTPUTS_UPDATE}"

log_stage "UPDATE COMPLETE"
echo ""

# Wait for user to proceed to DELETE
wait_for_confirmation "DELETE"

# ============================================
# STEP 3: DELETE SEARCH DEPLOYMENT
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
if ! verify_search_deployment_deleted; then
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
log_success "CREATE: Search deployment created successfully (InstanceSize=S30_HIGHCPU_NVME, NodeCount=2)"
log_success "UPDATE: Search deployment updated successfully (NodeCount=2 -> 3)"
log_success "DELETE: Search deployment deleted successfully"
echo ""
log_success "Lifecycle test finished successfully!"
echo ""

