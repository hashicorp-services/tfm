# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
