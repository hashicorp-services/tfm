# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

data "tfe_organization" "org" {
  name = var.organization
}

locals{
  projects = {
    project1 = "tfm-ci-test-0"
    project2 = "tfm-ci-test-1"
    project3 = "tfm-ci-test-3"
  }
}

resource "tfe_project" "source" {
  for_each     = local.projects
  organization = var.source_tfe_organization
  name         = each.value
}

resource "tfe_workspace" "ci-workspace-test" {
  name         = "ci-workspace-test"
  organization = data.tfe_organization.org.name
  tag_names    = ["test", "ci", "a"]
}

provider "tfe" {
  hostname     = var.tfe_hostname
  organization = var.organization
}

provider "tfe" {
  alias        = "source"
  hostname     = var.tfe_hostname
  organization = var.source_tfe_organization
  token        = var.source_tfe_token
}

provider "tfe" {
  alias        = "destination"
  hostname     = var.tfe_hostname
  organization = var.destination_tfe_organization
  token        = var.destination_tfe_token
}

# Source resources

#######################
### Test Scenario 1 ###
#######################

# Designed for testing the following:
# copy of vcs-driven workspaces
# Copy of state files
# copy of team access permissions
# copy of SSH Key mapping
# copy of variables
module "workspacer_vcs_driven" {

  source  = "app.terraform.io/tfm-testing-source/workspacer-tfm/tfe"
  version = "0.8.1"

  providers = {
    tfe = tfe.source
  }

  count = var.workspace_count

  organization   = var.source_tfe_organization
  workspace_name = "tfm-ci-test-vcs-${count.index}"
  workspace_desc = "Created by GitHub Actions CI e2e.workflow.yml. Created by Terraform Workspacer module."
  workspace_tags = ["agent", "ssh", "vcs-driven", "tfm"]
  force_delete   = true

  working_directory     = "/test/terraform/sample-resources/"
  auto_apply            = true
  file_triggers_enabled = true
  trigger_prefixes      = null 

  queue_all_runs = true
  #assessments_enabled   = true
  allow_destroy_plan = true
  #global_remote_state   = true

  #agent_pool_id  = tfe_agent_pool.source.id
  #execution_mode = "agent"

  ssh_key_id = tfe_ssh_key.source.id

  vcs_repo = {
    identifier     = "hashicorp-services/tfm"
    branch         = "main"
    oauth_token_id = tfe_oauth_client.source.oauth_token_id
    tags_regex     = null # conflicts with `trigger_prefixes` and `trigger_patterns`
  }

  tfvars = {
    teststring = "string"
    testlist   = ["1", "2", "3"]
    testmap    = { "a" = "1", "b" = "2", "c" = "3" }
  }

  team_access = {
    "${tfe_team.admin.name}" = "admin",
    "${tfe_team.appowner.name}" = "write",
    "${tfe_team.developer.name}" = "read"
  }

  depends_on = [
    tfe_team.admin,
    tfe_team.appowner,
    tfe_team.developer
  ]
}

#######################
### Test Scenario 2 ###
#######################

# Designed for testing the following:
# copy of vcs-driven workspace
# Copy of state files
# No copy of team access permissions
# No copy of SSH Key mapping
# No copy of variables
# No copy of agent pools mapping
module "workspacer_barebones" {
  source  = "app.terraform.io/tfm-testing-source/workspacer-tfm/tfe"
  version = "0.8.1"

  providers = {
    tfe = tfe.source
  }

  organization   = var.source_tfe_organization
  workspace_name = "tfm-ci-test-vcs-bare-bones"
  workspace_desc = "Created by GitHub Actions CI e2e.workflow.yml. Created by Terraform Workspacer module."
  workspace_tags = ["vcs-driven", "tfm"]
  force_delete   = true

  working_directory     = "/test/terraform/sample-resources/"
  auto_apply            = true
  file_triggers_enabled = true
  trigger_prefixes      = null 

  queue_all_runs = true
  allow_destroy_plan = true

  vcs_repo = {
    identifier     = "hashicorp-services/tfm"
    branch         = "main"
    oauth_token_id = tfe_oauth_client.source.oauth_token_id
    tags_regex     = null # conflicts with `trigger_prefixes` and `trigger_patterns`
  }
}

#######################
### Test Scenario 3 ###
#######################
# Designed to test the following:
# CLI-Driven workspace
# No state files present
module "workspacer_cli_driven" {
  source  = "app.terraform.io/tfm-testing-source/workspacer-tfm/tfe"
  version = "0.8.1"

  providers = {
    tfe = tfe.source
  }

  organization   = var.source_tfe_organization
  workspace_name = "tfm-ci-test-cli-nostate"
  workspace_desc = "Created by GitHub Actions CI e2e.workflow.yml. Created by Terraform Workspacer module."
  workspace_tags = ["cli-driven", "tfm"]
  force_delete   = true

  allow_destroy_plan = true

  ssh_key_id = tfe_ssh_key.source.id

  tfvars = {
    teststring = "string"
    testlist   = ["1", "2", "3"]
    testmap    = { "a" = "1", "b" = "2", "c" = "3" }
  }

  team_access = {
    "${tfe_team.appowner.name}" = "write",
    "${tfe_team.developer.name}" = "read"
  }

  depends_on = [
    tfe_team.appowner,
    tfe_team.developer
  ]
}

#######################
### Test Scenario 4 ###
#######################
# Designed for testing the following:
# Copy of agent pool mapping
module "workspacer_agent_execution" {

  source  = "app.terraform.io/tfm-testing-source/workspacer-tfm/tfe"
  version = "0.8.1"

  providers = {
    tfe = tfe.source
  }

  organization   = var.source_tfe_organization
  workspace_name = "tfm-ci-test-vcs-agent"
  workspace_desc = "Created by GitHub Actions CI e2e.workflow.yml. Created by Terraform Workspacer module."
  workspace_tags = ["agent", "ssh", "vcs-driven", "tfm"]
  force_delete   = true

  working_directory     = "/test/terraform/sample-resources/"
  auto_apply            = true
  file_triggers_enabled = true
  trigger_prefixes      = null 

  queue_all_runs = true
  #assessments_enabled   = true
  allow_destroy_plan = true
  #global_remote_state   = true

  agent_pool_id  = tfe_agent_pool.source.id
  execution_mode = "agent"

  vcs_repo = {
    identifier     = "hashicorp-services/tfm"
    branch         = "main"
    oauth_token_id = tfe_oauth_client.source.oauth_token_id
    tags_regex     = null # conflicts with `trigger_prefixes` and `trigger_patterns`
  }

  tfvars = {
    teststring = "string"
    testlist   = ["1", "2", "3"]
    testmap    = { "a" = "1", "b" = "2", "c" = "3" }
  }

  team_access = {
    "${tfe_team.developer.name}" = "read"
  }

  depends_on = [
    tfe_team.developer
  ]
}

# Destination Resources
resource "tfe_project" "migrated" {
  provider = tfe.destination
  name = "ci-test-migrated"
  organization = var.destination_tfe_organization
}