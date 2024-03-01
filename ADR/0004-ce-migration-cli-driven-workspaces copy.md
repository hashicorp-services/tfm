# ADR #0004: How to Handle CE/OSS to TFC Migration for CLI Driven Workspaces 

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

Not all users that migrate from Terraform communtiy edition managed configurations to Terraform Cloud managed worksapces want to use VCS driven workspaces. In situations where users want to migrate to CLI driven workspaces we should assist those users in configuring a `cloud {}` block automatically. This will help drive migration at scale.

CLI Driven workspaces require the following configuration at a minimum:

```
terraform {
  cloud {
    organization = "org-name"

    workspaces {
      name = "workspace-name"
    }
  }
}
```

## Decision

- A command will be implemented to assit users in adding the `cloud {}` configuration to terraform code.
- The `cloud {}` block will be automagically populated with the organization name taken from the `.tfm` config file setting `dst_tfc_org` and workspace name will be taken from the workspace name constructed from the metadata file created using `tfm core init-repos`.
- Any instances of `backend {}` will be commented out in favor of the `cloud{}` block.
- A VCS branch will be created, the change commmited, and pushed, but no PR will be created.

## Consequences

Users will be able configure terraform configurations for use with TFC/TFE CLI Driven workspaces at scale.