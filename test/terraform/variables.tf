variable "tfe_hostname" {
  description = "The TFE hostname"
  type        = string
  default     = "app.terraform.io"
}

variable "organization" {
  description = "The TFE Org"
  type        = string
  default     = "hc-implementation-services"
}

variable "gh_token" {
  description = "The Oauth Token for GitHub"
  type        = string
}

variable "source_tfe_token" {
  description = "The TFE token used for the source"
  type        = string
}

variable "destination_tfe_token" {
  description = "The TFE token used for the destination"
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
