# Testing Workflows for TFM

## Directory Structure

- actions contains different tests to be run using the tfm tool. Each test should have its own directory with the naming convention `test-<thing to test>`. Within that directoy an `action.yml` should exist containing the test actions. These are referenced by the `e2e.workflow.yml` workflow.
- terraform contains different terraform configurations for creating workspaces, teams, variables, etc within the `tfm-testing-source` workspace.
- workflows contains the different workflows related to TFM.

## End to End Testing ( e2e.workflow.yml )
`e2e.workflow.yml` is the entry point for end-to-end testing of TFM. It is scheduled to run daily.

### Setup
Information about the e2e workflow in the event the maintainer needs to issue new API tokens or troubleshoot.

- A long lived workspace named `gh-actions-ci-master-workspace` exists within the TFC organization `tfm-testing-source`. This workspace is a CLI driven workspace that executes terraform apply during the workflow to create all other resources within the org using the TFE provider, and later a terraform destroy to remove them.
- A team exists within the TFC organization `tfm-testing-source` named `github-actions-tfm`. Generate an API Token for this team.
- A team exists within the TFC organization `tfm-testing-destination` named `github-actions-tfm`. Generate an API Token for this team.
- A module is being sourced from `app.terraform.io/tfm-testing-source/workspacer-tfm/tfe`

### API Tokens
- Update the `gh-actions-ci-master-workspace` workspace variable `TFE_TOKEN` with the SOURCE org `github-actions-tfm` team API token each time you issue a new one.
- Update the `gh-actions-ci-master-workspace` workspace variable `tfm_source_token` with the SOURCE orgs `github-actions-tfm` team API token each time you issue a new one.
- Update the `tfm` repo github actions secret `SOURCETOKEN` with the SOURCE orgs `github-actions-tfm` team API token each time you issue a new one.
- Update the `gh-actions-ci-master-workspace` workspace variable `tfm_destination_token` with the DESTINATION orgs `github-actions-tfm` team API token each time you issue a new one.
- Update the `tfm` repo github actions secret `DESTINATIONTOKEN` with the DESTINATION orgs `github-actions-tfm` team API token each time you issue a new one.
- Update the `gh-actions-ci-master-workspace` workspace variable `gh_token` each time you issue a new one. The token string you were given by your VCS provider, e.g. ghp_xxxxxxxxxxxxxxx for a GitHub personal access token. For more information on how to generate this token string for your VCS provider, see the Create an OAuth Client documentation. This token is used for creating VCS connections to TFE/C Orgs.

### Whats Happening
1. The e2e workflow checks out the code in the `.github/workflows/terraform/tfe` folder and runs a CLI-driven run against the `gh-actions-ci-master-workspace` workspace within the TFC organization `tfm-testing-source`.
2. That creates workspaces, teams, variables, etc. Everything TFM will need to test tfm commands.
- - This creates resources in both the `tfm-testing-source` and `tfm-testing-destination` orgs.
3. All of the actions in `.github/actions` are run to test the various functions of TFM.
4. All of the resources that were copied using TFM to the TFC org `tfm-testing-destination` are deleted.
5. Terraform destroy is run from `gh-actions-ci-master-workspace` to clean up `tfm-testing-source` and `tfm-testing-destination` orgs.