# ADR #0006: How to Handle CE/OSS to TFC Migration for Variables

Date: 2024-02-29

## Responsible Architect
Joshua Tracy

## Author

Joshua Tracy

## Contributors

* Joshua Tracy
* Jeff McCollum
* Alex Basista

## Lifecycle

Pilot

## Status

Accepted

## Context

Part of migrating from Terraform Community Edition to Terraform Cloud or Enterprise is also migrating terraform inputs / variables in some fashion. A feature to assist users in making existing `.tfvars` files compatible with TFC/TFE runs or creating variables to be managed into TFC/TFE is required. 

Per [This Link](https://developer.hashicorp.com/terraform/cloud-docs/workspaces/variables) Although Terraform Cloud uses variables from terraform.tfvars, Terraform Enterprise currently ignores this file. A feature to locate all `*.tfvars` files in a given `config_path` could be created. The feature could rename the file to `terraform.auto.tfvars` for consistency across the board between TFC and TFE.

An additional feature in the form of a command or flad to an existing command could be used to allow users to specify the path to an existing vars file and the go-tfe sdk could be used to create those variables in the TFC/TFE workspace.

## Decision

- A command will be created for converting all `*.tfvars` files to `terraform.auto.tfvars` files.
- The `tfm core init-repos` commmand will be updated to identify and add the `tfvars` files in each `config_path`.
- A flag simliar to `-variables-path /path/to/terraform.tfvars` will be created to convert variables to workspaces managed variables.

## Consequences

Users will be able to modify or create terraform variables for TFC/TFE support at scale.