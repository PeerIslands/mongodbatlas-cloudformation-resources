#!/usr/bin/env bash
# cfn-test-create-inputs.sh
#
# This tool generates json files in the inputs/ for `cfn test`.
#

set -o errexit
set -o nounset
set -o pipefail
set -x

function usage {
    echo "usage:$0 <project_name>"
}

if [ "$#" -ne 1 ]; then usage; fi
if [[ "$*" == help ]]; then usage; fi

rm -rf inputs
mkdir inputs
name="${1}"
projectName="${1}"
if [ "$#" -ne 1 ]; then usage; fi
if [[ "$*" == help ]]; then usage; fi

projectId=$(mongocli iam projects list --output json | jq --arg NAME "${projectName}" -r '.results[] | select(.name==$NAME) | .id')
if [ -z "$projectId" ]; then
    projectId=$(mongocli iam projects create "${projectName}" --output=json | jq -r '.id')
    echo -e "Created project \"${projectName}\" with id: ${projectId}\n"
else
    echo -e "FOUND project \"${projectName}\" with id: ${projectId}\n"
fi

echo "Created project \"${projectName}\" with id: ${projectId}"
userName="testuser"
userName=$(mongocli atlas dbusers list --projectId ${projectId} --output json  |  jq --arg NAME "testuser6" -r '.[] | select(.username==$NAME) | .username')
echo "Created project "
if [[ "$userName" == "" ]];then
  echo "Created project dddddd ${projectId}"
    mongocli config set skip_update_check true
    mongocli atlas dbuser create atlasAdmin --username testuser6  --projectId ${projectId}  --x509Type MANAGED
  echo -e "Created user \"${userName}\" \n"
else
    echo -e "FOUND user \"${userName}\" \n"
fi

 echo "Created project ${projectId}"

jq --arg pubkey "$ATLAS_PUBLIC_KEY" \
   --arg pvtkey "$ATLAS_PRIVATE_KEY" \
   --arg org "$ATLAS_ORG_ID" \
   --arg group_id "$projectId" \
   --arg userName "$userName" \
   '.OrgId?|=$org | .Usernames?|=[$userName] |  .Name?|=$userName | .ApiKeys.PublicKey?|=$pubkey | .ApiKeys.PrivateKey?|=$pvtkey' \
   "$(dirname "$0")/inputs_1_create.template.json" > "inputs/inputs_1_create.json"
jq --arg pubkey "$ATLAS_PUBLIC_KEY" \
   --arg pvtkey "$ATLAS_PRIVATE_KEY" \
   --arg org "$ATLAS_ORG_ID" \
   --arg group_id "$projectId" \
   --arg userName "$userName" \
    '.OrgId?|=$org | .Usernames?|=[$userName] |  .Name?|=$userName | .ApiKeys.PublicKey?|=$pubkey | .ApiKeys.PrivateKey?|=$pvtkey' \
      "$(dirname "$0")/inputs_1_update.template.json" > "inputs/inputs_1_update.json"
name="${name}- more B@d chars !@(!(@====*** ;;::"
jq --arg pubkey "$ATLAS_PUBLIC_KEY" \
   --arg pvtkey "$ATLAS_PRIVATE_KEY" \
   --arg org "$ATLAS_ORG_ID" \
   --arg group_id "$projectId" \
   --arg userName "$userName" \
    '.OrgId?|=$org | .Usernames?|=[$userName] |  .Name?|=$userName | .ApiKeys.PublicKey?|=$pubkey | .ApiKeys.PrivateKey?|=$pvtkey' \
     "$(dirname "$0")/inputs_1_invalid.template.json" > "inputs/inputs_1_invalid.json"


ls -l inputs
