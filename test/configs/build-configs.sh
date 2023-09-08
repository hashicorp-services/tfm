#!/usr/bin/env bash
set -euo pipefail

# This script is used to build the ssh-map, vcs-map, and agents-map within the config files.
# Because these are created with new IDs each time e2e testing is run, we need a way to create these each time the test is run.
SOURCE_SSH_KEY_ID=$(curl --header "Authorization: Bearer $SRC_TFE_TOKEN" --request GET "https://app.terraform.io/api/v2/organizations/$SRC_TFE_ORG/ssh-keys" | jq '.data[] | select(.attributes.name == "tfm-ci-testing-src") | .id' | tr -d '"')
DESTINATION_SSH_KEY_ID=$(curl --header "Authorization: Bearer $DST_TFC_TOKEN" --request GET "https://app.terraform.io/api/v2/organizations/$DST_TFC_ORG/ssh-keys" | jq '.data[] | select(.attributes.name == "tfm-ci-testing-dest") | .id' | tr -d '"')

SOURCE_AGENTPOOL_ID=$(curl --header "Authorization: Bearer $SRC_TFE_TOKEN" --request GET "https://app.terraform.io/api/v2/organizations/$SRC_TFE_ORG/agent-pools" | jq '.data[] | select(.attributes.name == "tfm-ci-testing-src") | .id' | tr -d '"')
DESTINATION_AGEENTPOOL_ID=$(curl --header "Authorization: Bearer $DST_TFC_TOKEN" --request GET "https://app.terraform.io/api/v2/organizations/$DST_TFC_ORG/agent-pools" | jq '.data[] | select(.attributes.name == "tfm-ci-testing-dest") | .id' | tr -d '"')



SOURCE_OAUTH_CLIENT_ID=$(curl --header "Authorization: Bearer $SRC_TFE_TOKEN" --request GET "https://app.terraform.io/api/v2/organizations/$SRC_TFE_ORG/oauth-clients" | jq '.data[] | select(.attributes.name == "github-hashicorp-services-ci") | .id' | tr -d '"')
DESTINATION_OAUTH_CLIENT_ID=$(curl --header "Authorization: Bearer $DST_TFC_TOKEN" --request GET "https://app.terraform.io/api/v2/organizations/$DST_TFC_ORG/oauth-clients" | jq '.data[] | select(.attributes.name == "github-hashicorp-services-ci") | .id' | tr -d '"')

SOURCE_VCS_ID=$(curl --header "Authorization: Bearer $SRC_TFE_TOKEN" --request GET "https://app.terraform.io/api/v2/oauth-clients/$SOURCE_OAUTH_CLIENT_ID/oauth-tokens" | jq '.data[] | .id' | tr -d '"')
DESTINATION_VCS_ID=$(curl --header "Authorization: Bearer $DST_TFC_TOKEN" --request GET "https://app.terraform.io/api/v2/oauth-clients/$DESTINATION_OAUTH_CLIENT_ID/oauth-tokens" | jq '.data[] | .id' | tr -d '"')

cat > ./test/configs/.e2e-all-workspaces-test.hcl <<EOF
agents-map=[
  "$SOURCE_AGENTPOOL_ID=$DESTINATION_AGEENTPOOL_ID",
]

vcs-map=[
  "$SOURCE_VCS_ID=$DESTINATION_VCS_ID",
]

ssh-map=[
  "$SOURCE_SSH_KEY_ID=$DESTINATION_SSH_KEY_ID",
]
EOF

cat > ./test/configs/.e2e-workspace-map-test.hcl <<EOF
agents-map=[
  "$SOURCE_AGENTPOOL_ID=$DESTINATION_AGEENTPOOL_ID",
]

vcs-map=[
  "$SOURCE_VCS_ID=$DESTINATION_VCS_ID",
]

ssh-map=[
  "$SOURCE_SSH_KEY_ID=$DESTINATION_SSH_KEY_ID",
]

workspaces-map=[
  "tfm-ci-test-vcs-0=tfm-ci-test-vcs-0",
  "tfm-ci-test-vcs-1=tfm-ci-test-vcs-1",
  "tfm-ci-test-vcs-agent=tfm-ci-test-vcs-agent",
  "tfm-ci-test-vcs-4=new-tfm-ci-test-vcs-4",
  "tfm-ci-test-vcs-3=new-tfm-ci-test-vcs-3"
]
EOF

cat > ./test/configs/.e2e-workspaces-list-test.hcl <<EOF
agents-map=[
  "$SOURCE_AGENTPOOL_ID=$DESTINATION_AGEENTPOOL_ID",
]

vcs-map=[
  "$SOURCE_VCS_ID=$DESTINATION_VCS_ID",
]

ssh-map=[
  "$SOURCE_SSH_KEY_ID=$DESTINATION_SSH_KEY_ID",
]

"workspaces" = [
  "tfm-ci-test-vcs-0",
  "tfm-ci-test-vcs-1",
  "tfm-ci-test-cli-nostate",
  "tfm-ci-test-vcs-agent"
]
EOF

echo "[INFO] .e2e-all-workspaces-test.hcl"
cat ./test/configs/.e2e-all-workspaces-test.hcl

echo "[INFO] .e2e-workspace-map-test.hcl"
cat ./test/configs/.e2e-workspace-map-test.hcl

echo "[INFO] .e2e-workspaces-list-test.hcl"
cat ./test/configs/.e2e-workspaces-list-test.hcl