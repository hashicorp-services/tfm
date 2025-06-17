# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.14.0](https://github.com/hashicorp-services/tfm/compare/v0.13.0...v0.14.0) (2025-05-16)

### Features

- Implements [#305](https://github.com/hashicorp-services/tfm/issues/305)
- When using the remote state copy command, support workspace names changing between source and destination.

### Chore

- Fix links on tfm docs site

## [0.13.0](https://github.com/hashicorp-services/tfm/compare/v0.12.1...v0.13.0) (2025-05-07)

### Features

- Implements [#285] (https://github.com/hashicorp-services/tfm/issues/285)
- When migration workspace VCS connections, the source and destination VCS type don't need to match.

### Bug Fixes

- Fixes issue [#286] (https://github.com/hashicorp-services/tfm/issues/286)
- Removes the unused flags from the copy workspaces command

### Chore

- Update Codeowners: Thank you @carljavier & @Josh-Tracy for all your `tfm` contributions!
- Remove pr-title action
- Update go-tfe version. Unable to update viper due to breaking change.

## [0.12.1](https://github.com/hashicorp-services/tfm/compare/v0.12.0...v0.12.1) (2025-02-27)

### Bug Fixes

- Fixes issue [#283](https://github.com/hashicorp-services/tfm/issues/283)
- Remove the nuke command due to the above issue

## [0.12.0](https://github.com/hashicorp-services/tfm/compare/v0.11.3...v0.12.0) (2025-02-14)

### Features

- Add Run Trigger support to copy a workspaces Run Trigger setting [[#281](https://github.com/hashicorp-services/tfm/issues/281)]

### Chore

- Update documentation
- Update e2e tests

## [0.11.3](https://github.com/hashicorp-services/tfm/compare/v0.11.2...v0.11.3) (2025-01-29)

### Bug Fixes

- Checks if the OAuthID is nil before output. [[#273](https://github.com/hashicorp-services/tfm/issues/273)]

## [0.11.2](https://github.com/hashicorp-services/tfm/compare/v0.11.1...v0.11.2) (2025-01-28)

### Bug Fixes

- Moves the GHA VCS list into it's own command due to token limitations. [[#272](https://github.com/hashicorp-services/tfm/issues/272)]

## [0.11.1](https://github.com/hashicorp-services/tfm/compare/v0.11.0...v0.11.1) (2025-01-24)

### Features

- Update the `list vcs` command and the `copy workspaces --vcs`  flag to support GitHub App connections instead  of only support OAuth VCS connections. [[#268](https://github.com/hashicorp-services/tfm/issues/268)]

## [0.11.0](https://github.com/hashicorp-services/tfm/compare/v0.10.0...v0.11.0) (2025-01-10)

### Features

- Add flag to copy the workspace state sharing settings [[#254](https://github.com/hashicorp-services/tfm/issues/254)]

## [0.10.0](https://github.com/hashicorp-services/tfm/compare/v0.9.3...v0.10.0) (2024-11-20)

### Features

- Add flag to sensitive variables from being copied [[#235](https://github.com/hashicorp-services/tfm/issues/235)]

### Chore

- Updated Dependencies
- Replaces gox for tfm builds
- Updated GitHub Actions

## [0.9.3](https://github.com/hashicorp-services/tfm/compare/v0.9.2...v0.9.3) (2024-03-27)

### Bug Fixes

- **TFE 202303-1 and earlier Support** Fix list workspaces command for TFE deployments prior to the introduction of projects.

## [0.9.2](https://github.com/hashicorp-services/tfm/compare/v0.9.1...v0.9.2) (2024-03-27)

### Bug Fixes

- **TFE 202401-2 Support** Upgrade go-tfe to 1.48 to support TFE 202401-2 and later versions of TFE.

## [0.9.1](https://github.com/hashicorp-services/tfm/compare/v0.9.0...v0.9.1) (2024-03-11)

### Bug Fixes

- **State Copy Workspace Erorr Log** Fix error log name to support Windows

## [0.9.0](https://github.com/hashicorp-services/tfm/compare/v0.8.0...v0.9.0) (2024-03-07)

### Features

- **Terraform CE to TFC/TFE Migration** Support for GitLab VCS.
- **Terraform CE to TFC/TFE Migration** New required configuration file setting `vcs_type` introduced for defining a supported vcs type
- **Terraform CE to TFC/TFE Migration** `tfm core remove-backend --comment` new `--comment` flag added to optionally comment out `backend {}` configurations instead of deleting them.

## [0.8.0](https://github.com/hashicorp-services/tfm/compare/v0.7.0...v0.8.0) (2024-03-01)

### Features

- **Terraform CE to TFC/TFE Migration** TFM now supports discovery and migration of VCS repos configured as monorepos ( repos with more than 1 terraform configuration ) and terraform community edition workspaces. See docs and `tfm core -h` for more info.
- **Terraform CE to TFC/TFE Migration:** `tfm core init-repos` command added to assist in migrations. Builds a metadata file in the tfm working directory with information about how each VCS repo is configured.
- **tfm core migrate:** `tfm core migrate` command now executes the `tfm core init-repos` command.

## [0.7.0](https://github.com/hashicorp-services/tfm/compare/v0.6.0...v0.7.0) (2024-02-21)

### Features

- **Terraform CE to TFC/TFE Migration** Support for assisting with terraform community edition migrations to TFC/TFE managed workspaces. See docs and `tfm core -h` for more info.
- **Copy Projects:** `tfm copy projects` feature added

## [0.6.0] (November 08, 2023)

- Added command to lock and unlock workspaces. [[#155](https://github.com/hashicorp-services/tfm/issues/155)]
- Updated Dependencies
  - `go-tfe` from 1.35.0 to 1.39.0
  - `cobra` from 1.7.0 to 1.8.0
  - `color` from 1.15.0 to 1.16.0

## [0.5.0] (October 09, 2023)

- Added auto approve flag support to the `nuke` command
- doc updates
- Added more outputs to the workspaces `json` [#150](https://github.com/hashicorp-services/tfm/issues/150)
- Added starter Jira GitHub Action to open issues when a new GH issue is opened
- Add the ability to set all destination workspaces to use a agent pool even if the source doesn't [#126](https://github.com/hashicorp-services/tfm/issues/126)
- Updated dependencies
  - `go-tfe` from 1.32.1 to 1.35.0
  - `viper` from 1.16.0 to 1.17.0
  
## [0.4.4] (September 25, 2023)

- When setting the new state serial in the destination state, read the source state file's serial instead of using a computed serial number.

## [0.4.3] (September 19, 2023)

- Error logs created during state copy failures are now unique with a timestamp instead of being overwritten.

## [0.4.2] (September 18, 2023)

- TFM now retries connections in the event network issues or API rate limiting prevents an operation from taking place. **issue** #143
- BUG FIX: An earlier introduction of the `tfm copy workspaces --state --last x` command prevented pagination when copying state files. Implemented a fix. TFM will now copy more than 20 state files per worksapce. **issue** #144
- TFM will now handle errors encountered when migrating a state file for a given workspace. In the event an error is encountered TFM will stop migrating states for the given workspace and proceed to the next workspace. It will output a `workspace_error_log.txt` file in the working directory.

## [0.4.1] (September 14, 2023)

- BUG FIX: Applies after state migration would not upload the generated state file and output a 409 conflict error. TFM was not applying lineage to the workspace, but was copying the state file that contained a lineage. TFM now sets the lineage for the state files at the workspace level to match the one in the migrated state files.  **issue** #139

## [0.1.0] (June 8, 2023)

- TFM is now in beta and Open Source! :fireworks:
- Added warning when no workspaces are configured
- Rename workspace-map to workspaces-map
- Fix VCS list function to use destination side when specified
- Added Delete cmd to delete a single workspace by it's name or unique ID
- Updated Repo per legal and compliance standards
- Added `CODEOWNERS`
- Added delete vcs connection to remove a VCS connection from a workspace
- Added json output to selected commands
- Updated GO versions and SDKs
- Updated docs for beta release

## [0.0.5-pre-alpha] (April 27, 2023)

- Fixed migrated State lineage optional attribute
- Added tfm nuke workspaces (hidden command)
- Add tfm vcs list for configured org

## [0.0.4-pre-alpha] (April 14, 2023)

- Added tfm list workspaces
- Updated GitHub Pages for public viewing

## [0.0.3-pre-alpha] (April 6, 2023)

- Refactored Environment Variables/Config File
- Ability to specify a destination project ID for workspaces to be created in per tfm copy
- tfm list projects
- Updates to docs

## [0.0.2-pre-alpha] (March 3, 2023)

- Bug fix with state migration

## [0.0.1-pre-alpha] (February 18, 2023)

- Ready for Testing
- tfm copy workspaces
- tfm list
- Migration workflow

## [Unreleased] (November 24, 2022)

**Added**

- initial repository layout/boiler plate created
- copy teams subcommand
- list orgs subcommand
- list teams subcommand

**Changed**

- structure of `cmd` and its subcommands. Each sub command is now its in own directory and considered a package the `cmd` package imports in.
- `tfclient` pkg created at top level, not related to `cmd` package

**Removed**
