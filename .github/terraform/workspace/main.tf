# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

resource "null_resource" "workspace-a-resource" {}

variable "pet_count" {
  type        = number
  description = "Count of random_pet."
  default     = 10
}

variable "length" {
  type        = number
  description = "Length of random_pet."
  default     = 3
}

resource "random_pet" "main" {
  count = var.pet_count

  length    = var.length
  separator = "-"

  keepers = {
    always = timestamp()
  }
}