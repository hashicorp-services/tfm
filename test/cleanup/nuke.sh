#!/usr/bin/env bash
set -euo pipefail

# This script is used to wipe out changes that tfm unit testing will do
# Eventually this will be variablized and smarter but for now it is using hardcoded names and values

if $RUNNUKE = "true"
then

    echo "Removing workspaces"

    curl --header "Authorization: Bearer $DESTINATIONTOKEN" --request DELETE "https://app.terraform.io/api/v2/organizations/$DESTINATIONORGANIZATION/workspaces/tfc-mig-vcs-0"
    curl --header "Authorization: Bearer $DESTINATIONTOKEN" --request DELETE "https://app.terraform.io/api/v2/organizations/$DESTINATIONORGANIZATION/workspaces/tfc-mig-vcs-1"
    curl --header "Authorization: Bearer $DESTINATIONTOKEN" --request DELETE "https://app.terraform.io/api/v2/organizations/$DESTINATIONORGANIZATION/workspaces/tfc-mig-vcs-2"
    curl --header "Authorization: Bearer $DESTINATIONTOKEN" --request DELETE "https://app.terraform.io/api/v2/organizations/$DESTINATIONORGANIZATION/workspaces/tfc-mig-vcs-30"
    curl --header "Authorization: Bearer $DESTINATIONTOKEN" --request DELETE "https://app.terraform.io/api/v2/organizations/$DESTINATIONORGANIZATION/workspaces/tfc-mig-vcs-40"

    echo "Removing Team"

    TEAMID=$(curl --header "Authorization: Bearer $DESTINATIONTOKEN" --request GET "https://app.terraform.io/api/v2/organizations/$DESTINATIONORGANIZATION/teams" |  jq '.data[] | select(.attributes.name == "tfc-team") | .id' | tr -d '"')

    curl --header "Authorization: Bearer $DESTINATIONTOKEN" --request DELETE "https://app.terraform.io/api/v2/teams/$TEAMID"

    echo "Removing Varset"

    VARSETID=$(curl --header "Authorization: Bearer $DESTINATIONTOKEN" --request GET "https://app.terraform.io/api/v2/organizations/$DESTINATIONORGANIZATION/varsets" | jq '.data[] | select(.attributes.name == "source-varset") | .id' | tr -d '"')

    curl --header "Authorization: Bearer $DESTINATIONTOKEN" --request DELETE "https://app.terraform.io/api/v2/varsets/$VARSETID"

    echo "Target Nuked!"
else
    echo "Not running Nuke"
fi


