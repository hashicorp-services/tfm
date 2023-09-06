# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

resource "tfe_agent_pool" "source" {
  provider = tfe.source

  name = "tfm-ci-testing-src"
}

resource "tfe_agent_pool" "destination" {
  provider = tfe.destination

  name = "tfm-ci-testing-dest"
}