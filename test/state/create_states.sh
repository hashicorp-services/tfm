#!/bin/bash
set -euo pipefail
#
# This script is used create multiple states for a TFE/TFC workspace for testing purposes. 
# An assumption is the workspace you are calling has auto-apply set and regardless will create apply for every run.
#
#  Usage: ./create_states.sh <Workspace name> <number of states to create>
#  example:  ./create_states.sh tfc-mig-state-test 10
#
# Currently the wait time between calling an API run is 8 seconds which is hard coded.
#



WORKSPACE_NAME=$1
ORG_NAME='tfm-testing-source'


echo "Look Up the Workspace ID"
WORKSPACE_ID=($(curl \
  --header "Authorization: Bearer $TF_TOKEN" \
  --header "Content-Type: application/vnd.api+json" \
  https://app.terraform.io/api/v2/organizations/$ORG_NAME/workspaces/$WORKSPACE_NAME \
  | jq -r '.data.id'))


echo "Creating Payload"

cat << EOF > create_state_payload.json
{
  "data": {
    "attributes": {
      "message": "Create new State"
    },
    "type":"runs",
    "relationships": {
      "workspace": {
        "data": {
          "type": "workspaces",
          "id": "$WORKSPACE_ID"
        }
      }
    }
  }
}

EOF


counter=1


while [ $counter -le $2 ]
do
    echo "\n Creating Run $counter for Workspace: $WORKSPACE_NAME :: ID: $WORKSPACE_ID"
    curl \
    --header "Authorization: Bearer $TF_TOKEN" \
    --header "Content-Type: application/vnd.api+json" \
    --request POST \
    --data @create_state_payload.json \
    https://app.terraform.io/api/v2/runs
    ((counter++))
    sleep 8
done;
