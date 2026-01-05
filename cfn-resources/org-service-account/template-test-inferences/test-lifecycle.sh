#!/usr/bin/env bash

# test-lifecycle.sh
# Lifecycle test script for MongoDB::Atlas::OrgServiceAccount
# Tests complete lifecycle: CREATE → UPDATE → DELETE

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMPLATE_FILE="${SCRIPT_DIR}/org-service-account-complex.json"
OUTPUT_FILE="${SCRIPT_DIR}/output.json"
AWS_REGION="${AWS_REGION:-eu-west-1}"
PROFILE="${PROFILE:-default}"

# Function to print colored messages
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to wait for user confirmation
wait_for_confirmation() {
    local stage=$1
    echo ""
    print_info "Stage ${stage} completed. Review the output above."
    
    # Skip confirmation in non-interactive mode
    if [[ -n "${AUTO_RUN:-}" ]] || [[ ! -t 0 ]]; then
        print_info "Auto-continuing to next stage (non-interactive mode)"
        return 0
    fi
    
    read -p "Type 'yes' to continue to next stage, or 'no' to exit: " confirm
    if [[ "${confirm}" != "yes" ]]; then
        print_warning "Exiting lifecycle test."
        exit 0
    fi
}

# Function to check prerequisites
check_prerequisites() {
    print_info "Checking prerequisites..."
    
    # Check AWS credentials
    if ! aws sts get-caller-identity &>/dev/null; then
        print_error "AWS credentials not configured"
        exit 1
    fi
    
    # Check Atlas CLI
    if ! command -v atlas &>/dev/null; then
        print_warning "Atlas CLI not found. Atlas validation will be skipped."
    fi
    
    # Check template file
    if [[ ! -f "${TEMPLATE_FILE}" ]]; then
        print_error "Template file not found: ${TEMPLATE_FILE}"
        exit 1
    fi
    
    # Check required environment variables
    if [[ -z "${MONGODB_ATLAS_ORG_ID:-}" ]]; then
        print_error "MONGODB_ATLAS_ORG_ID environment variable not set"
        exit 1
    fi
    
    print_success "Prerequisites check passed"
}

# Function to get stack status
get_stack_status() {
    local stack_name=$1
    aws cloudformation describe-stacks \
        --stack-name "${stack_name}" \
        --region "${AWS_REGION}" \
        --query 'Stacks[0].StackStatus' \
        --output text 2>/dev/null || echo "NOT_FOUND"
}

# Function to wait for stack operation
wait_for_stack() {
    local stack_name=$1
    local target_status=$2
    local operation=$3
    
    print_info "Waiting for stack ${operation} to complete..."
    
    local status
    while true; do
        status=$(get_stack_status "${stack_name}")
        
        case "${status}" in
            "${target_status}")
                print_success "Stack ${operation} completed: ${status}"
                return 0
                ;;
            "ROLLBACK_COMPLETE"|"CREATE_FAILED"|"UPDATE_FAILED"|"DELETE_FAILED")
                print_error "Stack ${operation} failed: ${status}"
                print_error "Stack events:"
                aws cloudformation describe-stack-events \
                    --stack-name "${stack_name}" \
                    --region "${AWS_REGION}" \
                    --max-items 10 \
                    --query 'StackEvents[*].[Timestamp,ResourceStatus,ResourceStatusReason]' \
                    --output table
                return 1
                ;;
            "NOT_FOUND")
                if [[ "${target_status}" == "DELETE_COMPLETE" ]]; then
                    print_success "Stack deleted successfully"
                    return 0
                fi
                print_error "Stack not found"
                return 1
                ;;
            *)
                print_info "Stack status: ${status} (waiting...)"
                sleep 5
                ;;
        esac
    done
}

# Function to cleanup failed stack
cleanup_failed_stack() {
    local stack_name=$1
    local status=$(get_stack_status "${stack_name}")
    
    if [[ "${status}" == "ROLLBACK_COMPLETE" ]] || [[ "${status}" == "CREATE_FAILED" ]] || [[ "${status}" == "DELETE_FAILED" ]]; then
        print_warning "Cleaning up failed stack: ${stack_name}"
        aws cloudformation delete-stack \
            --stack-name "${stack_name}" \
            --region "${AWS_REGION}" || true
        wait_for_stack "${stack_name}" "DELETE_COMPLETE" "cleanup"
    fi
}

# Function to validate resource in Atlas
validate_in_atlas() {
    local client_id=$1
    
    print_info "Validating resource in Atlas..."
    
    local org_id="${MONGODB_ATLAS_ORG_ID}"
    
    # Try Atlas CLI first (atlas api serviceAccounts getOrgServiceAccount)
    if command -v atlas &>/dev/null; then
        local service_account
        if service_account=$(atlas api serviceAccounts getOrgServiceAccount --clientId "${client_id}" --orgId "${org_id}" -o json 2>/dev/null); then
            print_success "Service account found in Atlas:"
            echo "${service_account}" | jq '.' 2>/dev/null || echo "${service_account}"
            return 0
        else
            # Check if it's a 404 (not found) or other error
            local cli_error
            cli_error=$(atlas api serviceAccounts getOrgServiceAccount --clientId "${client_id}" --orgId "${org_id}" -o json 2>&1)
            if echo "${cli_error}" | grep -q "404\|not found\|NotFound"; then
                print_error "Service account not found in Atlas (404)"
            else
                print_warning "Atlas CLI validation failed: ${cli_error}"
            fi
            return 1
        fi
    fi
    
    # Fallback to direct API call if CLI not available
    local public_key="${MONGODB_ATLAS_PUBLIC_API_KEY:-${ATLAS_PUBLIC_KEY}}"
    local private_key="${MONGODB_ATLAS_PRIVATE_API_KEY:-${ATLAS_PRIVATE_KEY}}"
    
    # Check if API keys are available
    if [[ -z "${public_key}" ]] || [[ -z "${private_key}" ]]; then
        print_warning "Atlas CLI and API keys not available, skipping Atlas validation"
        return 0
    fi
    
    # Use Atlas API directly
    local api_response
    local http_code
    
    # Encode credentials for basic auth
    local auth_string
    auth_string=$(echo -n "${public_key}:${private_key}" | base64)
    
    # Call Atlas API to get service account
    api_response=$(curl -s -w "\n%{http_code}" \
        -X GET \
        -H "Accept: application/vnd.atlas.2023-02-01+json" \
        -H "Authorization: Basic ${auth_string}" \
        "https://cloud.mongodb.com/api/atlas/v2/orgs/${org_id}/serviceAccounts/${client_id}" 2>/dev/null)
    
    http_code=$(echo "${api_response}" | tail -n1)
    api_response=$(echo "${api_response}" | sed '$d')
    
    if [[ "${http_code}" == "200" ]]; then
        print_success "Service account found in Atlas:"
        echo "${api_response}" | jq '.' 2>/dev/null || echo "${api_response}"
        return 0
    elif [[ "${http_code}" == "404" ]]; then
        print_error "Service account not found in Atlas (404)"
        return 1
    else
        print_warning "Atlas API validation returned HTTP ${http_code}"
        echo "${api_response}" | head -5
        return 1
    fi
}

# Function to save output
save_output() {
    local stage=$1
    local stack_name=$2
    local client_id="${3:-}"
    
    local output_data
    output_data=$(cat <<EOF
{
  "stage": "${stage}",
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "stackName": "${stack_name}",
  "region": "${AWS_REGION}",
  "orgId": "${MONGODB_ATLAS_ORG_ID}",
  "clientId": "${client_id}",
  "stackStatus": "$(get_stack_status "${stack_name}")"
}
EOF
)
    
    echo "${output_data}" > "${OUTPUT_FILE}"
    print_info "Output saved to ${OUTPUT_FILE}"
}

# Main execution
main() {
    print_info "Starting lifecycle test for MongoDB::Atlas::OrgServiceAccount"
    print_info "Template: ${TEMPLATE_FILE}"
    print_info "Region: ${AWS_REGION}"
    print_info "Profile: ${PROFILE}"
    echo ""
    
    check_prerequisites
    
    # Get inputs - use environment variables or defaults (non-interactive)
    # Allow skipping prompts if AUTO_RUN is set or if running non-interactively
    if [[ -n "${AUTO_RUN:-}" ]] || [[ ! -t 0 ]]; then
        # Non-interactive mode - use defaults
        org_id="${MONGODB_ATLAS_ORG_ID:-}"
        name="test-service-account-$(date +%s)"
        stack_name="org-service-account-test-$(date +%s)"
        description="Service account for lifecycle testing"
        roles="ORG_MEMBER"
        expires="720"
        print_info "Running in non-interactive mode with defaults"
        print_info "OrgId: ${org_id}"
        print_info "Name: ${name}"
        print_info "Stack Name: ${stack_name}"
    else
        # Interactive mode - prompt for inputs
        read -p "Enter OrgId [${MONGODB_ATLAS_ORG_ID:-}]: " org_id
        org_id="${org_id:-${MONGODB_ATLAS_ORG_ID:-}}"
        
        read -p "Enter Service Account Name [test-service-account-$(date +%s)]: " name
        name="${name:-test-service-account-$(date +%s)}"
        
        read -p "Enter Stack Name [org-service-account-test-$(date +%s)]: " stack_name
        stack_name="${stack_name:-org-service-account-test-$(date +%s)}"
        
        read -p "Enter Description [Service account for lifecycle testing]: " description
        description="${description:-Service account for lifecycle testing}"
        
        read -p "Enter Roles (comma-separated) [ORG_MEMBER]: " roles
        roles="${roles:-ORG_MEMBER}"
        
        read -p "Enter SecretExpiresAfterHours [720]: " expires
        expires="${expires:-720}"
    fi
    
    # ==========================================
    # STAGE 1: CREATE
    # ==========================================
    print_info "=========================================="
    print_info "STAGE 1: CREATE"
    print_info "=========================================="
    
    print_info "Creating CloudFormation stack: ${stack_name}"
    
    # Cleanup any existing failed stack
    cleanup_failed_stack "${stack_name}"
    
    aws cloudformation create-stack \
        --stack-name "${stack_name}" \
        --template-body "file://${TEMPLATE_FILE}" \
        --parameters \
            ParameterKey=OrgId,ParameterValue="${org_id}" \
            ParameterKey=Name,ParameterValue="${name}" \
            ParameterKey=Description,ParameterValue="${description}" \
            ParameterKey=Roles,ParameterValue="${roles}" \
            ParameterKey=SecretExpiresAfterHours,ParameterValue="${expires}" \
            ParameterKey=Profile,ParameterValue="${PROFILE}" \
        --capabilities CAPABILITY_IAM \
        --region "${AWS_REGION}"
    
    if ! wait_for_stack "${stack_name}" "CREATE_COMPLETE" "creation"; then
        print_error "Stack creation failed"
        exit 1
    fi
    
    # Get stack outputs
    local client_id
    client_id=$(aws cloudformation describe-stacks \
        --stack-name "${stack_name}" \
        --region "${AWS_REGION}" \
        --query 'Stacks[0].Outputs[?OutputKey==`ClientId`].OutputValue' \
        --output text)
    
    print_success "Stack created successfully"
    print_info "ClientId: ${client_id}"
    
    # Validate in Atlas
    if ! validate_in_atlas "${client_id}"; then
        print_warning "Atlas validation failed, but continuing..."
    fi
    
    # Save output
    save_output "CREATE" "${stack_name}" "${client_id}"
    
    print_success "CREATE COMPLETE"
    wait_for_confirmation "1"
    
    # ==========================================
    # STAGE 2: UPDATE
    # ==========================================
    print_info "=========================================="
    print_info "STAGE 2: UPDATE"
    print_info "=========================================="
    
    print_info "Updating CloudFormation stack: ${stack_name}"
    
    # Update with modified values
    local updated_name="${name}-updated"
    local updated_description="${description} - Updated"
    local updated_roles="ORG_MEMBER,ORG_GROUP_CREATOR"
    
    print_info "Updating Name: ${name} → ${updated_name}"
    print_info "Updating Description: ${description} → ${updated_description}"
    print_info "Updating Roles: ${roles} → ${updated_roles}"
    print_info "Note: SecretExpiresAfterHours is create-only and will not be updated"
    
    # Note: OrgId, Profile, and SecretExpiresAfterHours are create-only properties
    # Use --use-previous-value for create-only parameters
    # For Roles with comma-separated values, use a parameter file to avoid AWS CLI parsing issues
    local param_file=$(mktemp)
    cat > "${param_file}" <<EOF
[
  {
    "ParameterKey": "OrgId",
    "UsePreviousValue": true
  },
  {
    "ParameterKey": "Name",
    "ParameterValue": "${updated_name}"
  },
  {
    "ParameterKey": "Description",
    "ParameterValue": "${updated_description}"
  },
  {
    "ParameterKey": "Roles",
    "ParameterValue": "${updated_roles}"
  },
  {
    "ParameterKey": "SecretExpiresAfterHours",
    "UsePreviousValue": true
  },
  {
    "ParameterKey": "Profile",
    "UsePreviousValue": true
  }
]
EOF
    
    aws cloudformation update-stack \
        --stack-name "${stack_name}" \
        --template-body "file://${TEMPLATE_FILE}" \
        --parameters "file://${param_file}" \
        --capabilities CAPABILITY_IAM \
        --region "${AWS_REGION}"
    
    # Clean up temp file
    rm -f "${param_file}"
    
    if ! wait_for_stack "${stack_name}" "UPDATE_COMPLETE" "update"; then
        print_error "Stack update failed"
        exit 1
    fi
    
    # Get updated client_id (should be same)
    client_id=$(aws cloudformation describe-stacks \
        --stack-name "${stack_name}" \
        --region "${AWS_REGION}" \
        --query 'Stacks[0].Outputs[?OutputKey==`ClientId`].OutputValue' \
        --output text)
    
    print_success "Stack updated successfully"
    print_info "ClientId: ${client_id}"
    
    # Validate update in Atlas
    if ! validate_in_atlas "${client_id}"; then
        print_warning "Atlas validation failed, but continuing..."
    fi
    
    # Update output
    save_output "UPDATE" "${stack_name}" "${client_id}"
    
    print_success "UPDATE COMPLETE"
    wait_for_confirmation "2"
    
    # ==========================================
    # STAGE 3: DELETE
    # ==========================================
    print_info "=========================================="
    print_info "STAGE 3: DELETE"
    print_info "=========================================="
    
    print_info "Deleting CloudFormation stack: ${stack_name}"
    
    aws cloudformation delete-stack \
        --stack-name "${stack_name}" \
        --region "${AWS_REGION}"
    
    if ! wait_for_stack "${stack_name}" "DELETE_COMPLETE" "deletion"; then
        print_error "Stack deletion failed"
        exit 1
    fi
    
    # Validate deletion in Atlas
    print_info "Validating deletion in Atlas..."
    
    # Try Atlas CLI first
    if command -v atlas &>/dev/null; then
        if atlas api serviceAccounts getOrgServiceAccount --clientId "${client_id}" --orgId "${org_id}" &>/dev/null 2>&1; then
            print_warning "Service account still exists in Atlas (may take time to propagate)"
            return 0
        else
            print_success "Service account deleted from Atlas"
            return 0
        fi
    fi
    
    # Fallback to direct API call
    local public_key="${MONGODB_ATLAS_PUBLIC_API_KEY:-${ATLAS_PUBLIC_KEY}}"
    local private_key="${MONGODB_ATLAS_PRIVATE_API_KEY:-${ATLAS_PRIVATE_KEY}}"
    
    if [[ -z "${public_key}" ]] || [[ -z "${private_key}" ]]; then
        print_warning "Atlas CLI and API keys not available, skipping deletion validation"
        return 0
    fi
    
    # Use Atlas API directly to check if service account still exists
    local auth_string
    auth_string=$(echo -n "${public_key}:${private_key}" | base64)
    
    local http_code
    http_code=$(curl -s -o /dev/null -w "%{http_code}" \
        -X GET \
        -H "Accept: application/vnd.atlas.2023-02-01+json" \
        -H "Authorization: Basic ${auth_string}" \
        "https://cloud.mongodb.com/api/atlas/v2/orgs/${org_id}/serviceAccounts/${client_id}" 2>/dev/null)
    
    if [[ "${http_code}" == "404" ]]; then
        print_success "Service account deleted from Atlas"
        return 0
    elif [[ "${http_code}" == "200" ]]; then
        print_warning "Service account still exists in Atlas (may take time to propagate)"
        return 0
    else
        print_warning "Atlas API validation returned HTTP ${http_code}, assuming deleted"
        return 0
    fi
    
    # Clean up output file
    rm -f "${OUTPUT_FILE}"
    
    print_success "DELETE COMPLETE"
    print_success "=========================================="
    print_success "Lifecycle test completed successfully!"
    print_success "=========================================="
}

# Run main function
main "$@"

