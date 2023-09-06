# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

resource "tfe_variable_set" "varset_aws" {
  provider = tfe.source

  name        = "tfm-ci-testing-varset-aws"
  description = "Variable set for aws credentials"
}

resource "tfe_variable_set" "varset_azure" {
  provider = tfe.source

  name        = "tfm-ci-testing-varset-azure"
  description = "Variable set for azure credentials"
}

resource "tfe_variable" "source_varset" {
  provider = tfe.source


  key             = "variable_set_var"
  value           = "test"
  category        = "terraform"
  description     = "variable description"
  variable_set_id = tfe_variable_set.source.id
}