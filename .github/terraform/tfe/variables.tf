# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "organization" {
  type = string
}

variable "tfe_hostname" {
  description = "The TFE hostname"
  type        = string
  default     = "app.terraform.io"
}

variable "gh_token" {
  description = "The Oauth Token for GitHub"
  type        = string
}

variable "source_tfe_organization" {
  description = "The Source TFE Organization"
  type        = string
  default     = "tfm-testing-source"
}

variable "destination_tfe_organization" {
  description = "The Destination TFE Organization"
  type        = string
  default     = "tfm-testing-destination"
}

variable "workspace_count" {
  description = "How many workspaces to create"
  type        = number
  default     = 5
}