#!/bin/bash

# Required variables
TFC_TOKEN=""
ORG_NAME="default"
API_URL="ace-blowfish.benjamin-lykins.sbx.hashidemos.io"

# Set headers
AUTH_HEADER="Authorization: Bearer $TFC_TOKEN"
CONTENT_HEADER="Content-Type: application/vnd.api+json"

# Step 1: Delete all workspaces
echo "Fetching all workspaces in org: $ORG_NAME..."
workspace_ids=$(curl \
  --header "Authorization: Bearer $TFC_TOKEN" \
  --header "Content-Type: application/vnd.api+json" \
    "https://$API_URL/api/v2/organizations/$ORG_NAME/workspaces?size=100" | \
jq -r '.data[].attributes.name')

if [ -z "$workspace_ids" ]; then
  echo "No workspaces found in organization: $ORG_NAME"

else
echo "$workspace_ids" | while IFS= read -r ws_id; do
  echo "Deleting workspace ID: $ws_id"
  curl \
    --header "Authorization: Bearer $TFC_TOKEN" \
    --header "Content-Type: application/vnd.api+json" \
    --request DELETE \
    "https://$API_URL/api/v2/organizations/$ORG_NAME/workspaces/$ws_id"
done
fi


# Step 2: Delete all projects
echo "Fetching all projects in org: $ORG_NAME..."
project_ids=$(curl \
  --header "Authorization: Bearer $TFC_TOKEN" \
  --header "Content-Type: application/vnd.api+json" \
  "https://$API_URL/api/v2/organizations/$ORG_NAME/projects" | \
  jq -r '.data[].id')

echo "$project_ids" | while IFS= read -r proj_id; do
  echo "Deleting project ID: $proj_id"
curl \
  --header "Authorization: Bearer $TFC_TOKEN" \
  --header "Content-Type: application/vnd.api+json" \
  --request DELETE \
  https://$API_URL/api/v2/projects/$proj_id \

done

echo "All workspaces and projects deleted."
