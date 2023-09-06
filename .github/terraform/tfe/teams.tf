# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

resource "tfe_team" "admin" {
  provider = tfe.source

  name = "tfm-ci-testing-admins"
}

resource "tfe_team" "appowner" {
  provider = tfe.source

  name = "tfm-ci-testing-appowner"
}

resource "tfe_team" "developer" {
  provider = tfe.source

  name = "tfm-ci-testing-developer"
}