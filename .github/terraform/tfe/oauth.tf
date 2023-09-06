# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

resource "tfe_oauth_client" "source" {
  provider = tfe.source

  name             = "github-hashicorp-services"
  api_url          = "https://api.github.com"
  http_url         = "https://github.com"
  oauth_token      = var.gh_token
  service_provider = "github"
}

resource "tfe_oauth_client" "destination" {
  provider = tfe.destination

  name             = "github-hashicorp-services"
  api_url          = "https://api.github.com"
  http_url         = "https://github.com"
  oauth_token      = var.gh_token
  service_provider = "github"
}