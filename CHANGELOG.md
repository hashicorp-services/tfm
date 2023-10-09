# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

- TFM now retries connections in the event network issues or API rate limiting prevents an operation from taking place. __issue__ #143
- BUG FIX: An earlier introduction of the `tfm copy workspaces --state --last x` command prevented pagination when copying state files. Implemented a fix. TFM will now copy more than 20 state files per worksapce. __issue__ #144
- TFM will now handle errors encountered when migrating a state file for a given workspace. In the event an error is encountered TFM will stop migrating states for the given workspace and proceed to the next workspace. It will output a `workspace_error_log.txt` file in the working directory.

## [0.4.1] (September 14, 2023)

- BUG FIX: Applies after state migration would not upload the generated state file and output a 409 conflict error. TFM was not applying lineage to the workspace, but was copying the state file that contained a lineage. TFM now sets the lineage for the state files at the workspace level to match the one in the migrated state files.  __issue__ #139

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
