#!/usr/bin/env bash
set -euo pipefail

# This script is used to wipe out changes that tfm unit testing will do
# Eventually this will be variablized and smarter but for now it is using hardcoded names and values


    echo "Removing workspaces"

    curl --header "Authorization: Bearer $DST_TFC_TOKEN" --request DELETE "https://app.terraform.io/api/v2/organizations/$DST_TFC_ORG/workspaces/tfm-ci-test-vcs-0"
    curl --header "Authorization: Bearer $DST_TFC_TOKEN" --request DELETE "https://app.terraform.io/api/v2/organizations/$DST_TFC_ORG/workspaces/tfm-ci-test-vcs-1"
    curl --header "Authorization: Bearer $DST_TFC_TOKEN" --request DELETE "https://app.terraform.io/api/v2/organizations/$DST_TFC_ORG/workspaces/tfm-ci-test-vcs-2"
    curl --header "Authorization: Bearer $DST_TFC_TOKEN" --request DELETE "https://app.terraform.io/api/v2/organizations/$DST_TFC_ORG/workspaces/tfm-ci-test-vcs-3"
    curl --header "Authorization: Bearer $DST_TFC_TOKEN" --request DELETE "https://app.terraform.io/api/v2/organizations/$DST_TFC_ORG/workspaces/tfm-ci-test-vcs-4"
    curl --header "Authorization: Bearer $DST_TFC_TOKEN" --request DELETE "https://app.terraform.io/api/v2/organizations/$DST_TFC_ORG/workspaces/tfm-ci-test-vcs-bare-bones"
    curl --header "Authorization: Bearer $DST_TFC_TOKEN" --request DELETE "https://app.terraform.io/api/v2/organizations/$DST_TFC_ORG/workspaces/tfm-ci-test-cli-nostate"
    curl --header "Authorization: Bearer $DST_TFC_TOKEN" --request DELETE "https://app.terraform.io/api/v2/organizations/$DST_TFC_ORG/workspaces/tfm-ci-test-vcs-agent"
    curl --header "Authorization: Bearer $DST_TFC_TOKEN" --request DELETE "https://app.terraform.io/api/v2/organizations/$DST_TFC_ORG/workspaces/ci-workspace-test"

    echo "Removing Teams"

    ADMINTEAMID=$(curl --header "Authorization: Bearer $DST_TFC_TOKEN" --request GET "https://app.terraform.io/api/v2/organizations/$DST_TFC_ORG/teams" |  jq '.data[] | select(.attributes.name == "tfm-ci-testing-admins") | .id' | tr -d '"')
    APPOWNERTEAMID=$(curl --header "Authorization: Bearer $DST_TFC_TOKEN" --request GET "https://app.terraform.io/api/v2/organizations/$DST_TFC_ORG/teams" |  jq '.data[] | select(.attributes.name == "tfm-ci-testing-appowner") | .id' | tr -d '"')
    DEVELOPERTEAMID=$(curl --header "Authorization: Bearer $DST_TFC_TOKEN" --request GET "https://app.terraform.io/api/v2/organizations/$DST_TFC_ORG/teams" |  jq '.data[] | select(.attributes.name == "tfm-ci-testing-developer") | .id' | tr -d '"')

    curl --header "Authorization: Bearer $DST_TFC_TOKEN" --request DELETE "https://app.terraform.io/api/v2/teams/$ADMINTEAMID"
    curl --header "Authorization: Bearer $DST_TFC_TOKEN" --request DELETE "https://app.terraform.io/api/v2/teams/$APPOWNERTEAMID"
    curl --header "Authorization: Bearer $DST_TFC_TOKEN" --request DELETE "https://app.terraform.io/api/v2/teams/$DEVELOPERTEAMID"

    echo "Removing Varsets"

    AWSVARSETID=$(curl --header "Authorization: Bearer $DST_TFC_TOKEN" --request GET "https://app.terraform.io/api/v2/organizations/$DST_TFC_ORG/varsets" | jq '.data[] | select(.attributes.name == "tfm-ci-testing-varset-aws") | .id' | tr -d '"')
    AZUREVARSETID=$(curl --header "Authorization: Bearer $DST_TFC_TOKEN" --request GET "https://app.terraform.io/api/v2/organizations/$DST_TFC_ORG/varsets" | jq '.data[] | select(.attributes.name == "tfm-ci-testing-varset-azure") | .id' | tr -d '"')

    curl --header "Authorization: Bearer $DST_TFC_TOKEN" --request DELETE "https://app.terraform.io/api/v2/varsets/$AWSVARSETID"
    curl --header "Authorization: Bearer $DST_TFC_TOKEN" --request DELETE "https://app.terraform.io/api/v2/varsets/$AZUREVARSETID"

    echo "Target Nuked!"


