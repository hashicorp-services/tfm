# ADR #0003: Feature - Terraform Community Edition to TFC Migration

Date: 2024-02-26

## Responsible Architect

Joshua Tracy

## Author

Joshua Tracy

## Contributors

* Joshua Tracy
* Jeff McCollum
* Alex Basista

## Lifecycle

Proof of Concept

## Status

Proposed

## Context

Users need a tool to assist in the migration from terraform community edition managed configurations to Terraform Cloud workspace managed configurations.

### Migration Workflow

Today migration can be accomplished following the following steps:

- Get the state file
- Create a TFC Workspace
- Push the state file to the TFC workspace
- Link the code to the workspace containing its code somehow

The above are the 4 high-level steps required for a migration to happen. Each step can be broken down into more steps and complexity can increase quickly based on the terraform configuration setup.

### TFM Capabilities

TFM must be able to do the following:

- Automate speedup the process of getting the state file.
- Determine the landscape of each terraform configuration.
- Automate and speedup the process of creating TFC workspace.
- Automate and speedup the process of uploading state files to worksapces.
- Automate and speedup tieing together the configuration code with the workspaces containing the respective state.
- Automate and speedup the post migration process of modifying code.

### What Kind of Terraform Configurations Exist In the Wild?

User typically have the following configurations:

- A single terraform configuration with a configured backend stored in the root of a VCS repository.
- Multiple terraform configurations each with a configured backend stored in multiple directory paths within a VCS repository.
- A single terraform configuration using multiple terraform workspaces with a single backend configuration.
- Multiple terraform configurations each terraform ce workspaces and a single backend configured for each configuration.

## Decision

- A function to assist users in retreiving the terraform configurations will be implemented in the form of `tfm core clone` to allow users to download VCS configurations. This will allow the process of getting the state file to be automated at scale.
- A function to build a metadata file that contains the information about how each cloned VCS repository is configured will be created. The metadata file will contain information about each path within the VCS repository that contains terraform code.
- A function to determine if a VCS repository contains multiple configurations will be created.
- A function to determine if a VCS repository contains configurations using terraform ce workspaces.
- A function will be created to assist users in downloading the state file for every identified terraform configuration.
- A function will be created to assist users in creating TFC workspaces.
- A function will be created to assist users in uploading state files to the created workspaces.
- A function will be created to assist in connecting the VCS repositories to the TFC workspaces.
- A function will be created to assist users in removing the old backend configuration if desired.
- The function to assist users in removing backend configurations will be optional.
- The function to assist users in removing the backend will create a branch and push it to the VCS.
- The function to assist users in removing the backend will NOT create a PR.
- A function to help users cleanup the working environment will be created.
- The tfm command that creates workspaces will create them with an invisible tfm tag to allow `tfm nuke workspaces` to delete all of the workspaces created by tfm.

## Consequences

- Users will be able to quickly migrate multiple terraform configurations within minutes VS days.
- Users will be able to quickly remove unwanted backend code from their configurations post migration.
- Users will not have the ability to choose their workspace names.
- No consideration for CLI driven or API driven workspaces.
- No consideration for configurations using terraform variable files.
- Only GitHub VCS will be supported initially.