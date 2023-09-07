# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

resource "tfe_variable_set" "aws" {
  provider = tfe.source

  name        = "tfm-ci-testing-varset-aws"
  description = "Variable set for aws credentials"
}

resource "tfe_variable_set" "azure" {
  provider = tfe.source

  name        = "tfm-ci-testing-varset-azure"
  description = "Variable set for azure credentials"
}

resource "tfe_variable" "variable1" {
  provider = tfe.source


  key             = "variable_set_var"
  value           = "test"
  category        = "terraform"
  description     = "variable description"
  variable_set_id = tfe_variable_set.aws.id
}

resource "tfe_variable" "variable2" {
  provider = tfe.source


  key             = "env_var"
  value           = "test"
  category        = "environment"
  description     = "variable description"
  variable_set_id = tfe_variable_set.aws.id
}

resource "tfe_variable" "variable3" {
  provider = tfe.source


  key             = "variable_set_var"
  value           = "test"
  category        = "terraform"
  description     = "variable description"
  variable_set_id = tfe_variable_set.azure.id
}

resource "tfe_variable" "variable4" {
  provider = tfe.source


  key             = "env_var"
  value           = "test"
  category        = "environment"
  description     = "variable description"
  variable_set_id = tfe_variable_set.azure.id
}