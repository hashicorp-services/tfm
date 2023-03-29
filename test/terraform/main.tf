terraform {
  required_providers {
    tfe = {
      version = "~> 0.42.0"
    }
  }
}

terraform {
  cloud {
    organization = "hc-implementation-services"

    workspaces {
      name = "unit-test-baseline"
    }
  }
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
resource "tfe_agent_pool" "source" {
  provider = tfe.source

  name = "tfc-source"
}

resource "tfe_ssh_key" "source" {
  provider = tfe.source

  name = "tfm-mig"
  key  = <<EOF
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACDCgSWWMT5TOfXZnpHB4uAg1e6+gd06QJASShdH7wHbawAAAKiSKdb5kinW
+QAAAAtzc2gtZWQyNTUxOQAAACDCgSWWMT5TOfXZnpHB4uAg1e6+gd06QJASShdH7wHbaw
AAAEAVisUyUHfpsDucm4wBomapQslHyWUwOAnjJcJcGnP5isKBJZYxPlM59dmekcHi4CDV
7r6B3TpAkBJKF0fvAdtrAAAAIGptY2NvbGx1bUBqbWNjb2xsdW0tQzAyRjcwQVVNRDZSAQ
IDBAU=
-----END OPENSSH PRIVATE KEY-----
EOF
}

resource "tfe_team" "source" {
  provider = tfe.source

  name = "tfc-team"
}

resource "tfe_variable_set" "source" {
  provider = tfe.source

  name        = "source-varset"
  description = "varset description"
}

resource "tfe_variable" "source_varset" {
  provider = tfe.source


  key             = "variable_set_var"
  value           = "test"
  category        = "terraform"
  description     = "variable description"
  variable_set_id = tfe_variable_set.source.id
}

resource "tfe_oauth_client" "source" {
  provider = tfe.source

  name             = "github-hashicorp-services"
  api_url          = "https://api.github.com"
  http_url         = "https://github.com"
  oauth_token      = var.gh_token
  service_provider = "github"
}

module "workspacer_source" {

  source  = "app.terraform.io/hc-implementation-services/workspacer-tfm/tfe"
  version = "0.8.1"

  providers = {
    tfe = tfe.source
  }

  count = var.workspace_count

  organization   = var.source_tfe_organization
  workspace_name = "tfc-mig-vcs-${count.index}"
  workspace_desc = "Created by Terraform Workspacer module."
  workspace_tags = ["agent", "ssh", "vcs-driven"]
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
    identifier     = "hashicorp-services/tfm2"
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
    "${tfe_team.source.name}" = "admin"
  }

  depends_on = [
    tfe_team.source,
  ]
}

# Destination Resources
resource "tfe_agent_pool" "destination" {
  provider = tfe.destination

  name = "tfc-destination"
}

resource "tfe_ssh_key" "destination" {
  provider = tfe.destination

  name = "tfm-mig"
  key  = <<EOF
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACDCgSWWMT5TOfXZnpHB4uAg1e6+gd06QJASShdH7wHbawAAAKiSKdb5kinW
+QAAAAtzc2gtZWQyNTUxOQAAACDCgSWWMT5TOfXZnpHB4uAg1e6+gd06QJASShdH7wHbaw
AAAEAVisUyUHfpsDucm4wBomapQslHyWUwOAnjJcJcGnP5isKBJZYxPlM59dmekcHi4CDV
7r6B3TpAkBJKF0fvAdtrAAAAIGptY2NvbGx1bUBqbWNjb2xsdW0tQzAyRjcwQVVNRDZSAQ
IDBAU=
-----END OPENSSH PRIVATE KEY-----
EOF
}

resource "tfe_oauth_client" "destination" {
  provider = tfe.destination

  name             = "github-hashicorp-services"
  api_url          = "https://api.github.com"
  http_url         = "https://github.com"
  oauth_token      = var.gh_token
  service_provider = "github"
}