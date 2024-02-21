# Testing Workflows for TFM

## Directory Structure

- actions contains different tests to be run using the tfm tool. Each test should have its own directory with the naming convention `test-<thing to test>`. Within that directoy an `action.yml` should exist containing the test actions. These are referenced by the `e2e.workflow.yml` workflow.
- terraform contains different terraform configurations for creating workspaces, teams, variables, etc within the `tfm-testing-source` workspace.
- workflows contains the different workflows related to TFM.

## End to End Testing ( e2e.workflow.yml )

`e2e.workflow.yml` is the entry point for end-to-end testing of TFM. It is scheduled to run daily.

### Setup

Information about the e2e workflow in the event the maintainer needs to issue new API tokens or troubleshoot.

- A github service account `svc-tfm` exists to facilitate creation of github tokens to create VCS connections using the tfe provider during testing.
- The password for the `svc-tfm` user is stored in Onepassword.
- A long lived workspace named `gh-actions-ci-master-workspace` exists within the TFC organization `tfm-testing-source`. This workspace is a CLI driven workspace that executes terraform apply during the workflow to create all other resources within the org using the TFE provider, and later a terraform destroy to remove them.
- A team exists within the TFC organization `tfm-testing-source` named `owners`. Generate an API Token for this team.
- A team exists within the TFC organization `tfm-testing-destination` named `owners`. Generate an API Token for this team.
- A module is being sourced from `app.terraform.io/tfm-testing-source/workspacer-tfm/tfe`
- `svc-tfm` has permissions to the repositories `hashicorp-services/tfm-oss-migration-rivendell`, `hashicorp-services/tfm-oss-migration-mordor2`, and `hashicorp-services/tfm-oss-migration-isengard` because these repos contain configurations using an s3 backend to store state for testing.

### API Tokens

- Update the `gh-actions-ci-master-workspace` workspace variable `TFE_TOKEN` with the SOURCE org `owner` team API token each time you issue a new one.
- Update the `gh-actions-ci-master-workspace` workspace variable `tfm_source_token` with the SOURCE orgs `owner` team API token each time you issue a new one.
- Update the `tfm` repo github actions secret `SOURCETOKEN` with the SOURCE orgs `owner` team API token each time you issue a new one.
- Update the `gh-actions-ci-master-workspace` workspace variable `tfm_destination_token` with the DESTINATION orgs `owner` team API token each time you issue a new one.
- Update the `tfm` repo github actions secret `DESTINATIONTOKEN` with the DESTINATION orgs `owner` team API token each time you issue a new one.
- Update the `gh-actions-ci-master-workspace` workspace variable `gh_token` each time you issue a new one. The token string you were given by your VCS provider, e.g. ghp_xxxxxxxxxxxxxxx for a GitHub personal access token. For more information on how to generate this token string for your VCS provider, see the Create an OAuth Client documentation. This token is used for creating VCS connections to TFE/C Orgs.
  - - The `gh_token` needs to be issued by the `svc-tfm`  github service account.
- Update the `tfm` repo github actions secret `GHORGANIZATION` with the hashicorp services github org.
- Update the `tfm` repo github actions secret `GHUSERNAME` with a github username with access to the organization.
- Update the `tfm` repo github actions secret `GHTOKEN` with a GitHub token. This token is used to clone repos to the GitHub actions pipeline for testing.
  - - The `GHTOKEN` needs to be issued by the `svc-tfm` gtihub service account with read-write contents permissions to the repositories being cloned.
- TO DO - FIGURE OUT HOW TO AUTH THE AN AWS ACCOUNT WITH GITHUB REPO STATE FILES STORED IN S3 BACKEND.

### Whats Happening

1. The e2e workflow checks out the code in the `.github/workflows/terraform/tfe` folder and runs a CLI-driven run against the `gh-actions-ci-master-workspace` workspace within the TFC organization `tfm-testing-source`.
2. That creates workspaces, teams, variables, etc. Everything TFM will need to test tfm commands.
   - Note: This creates resources in both the `tfm-testing-source` and `tfm-testing-destination` orgs.
3. All of the actions in `.github/actions` are run to test the various functions of TFM.
4. All of the resources that were copied using TFM to the TFC org `tfm-testing-destination` are deleted.
5. Terraform destroy is run from `gh-actions-ci-master-workspace` to clean up `tfm-testing-source` and `tfm-testing-destination` orgs.

## Jira (jira-issues.yml)

This action will use the org provided Jira service account to open issues in the ASE Service Portfolio Jira Board. This only runs when issues are created in the `tfm` repo and is not synced back from Jira to GitHub. Next `tfm` will support more functionality.