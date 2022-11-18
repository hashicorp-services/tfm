# Repository Layout

Information on how this repository is structured.

## Files

`main.go` - Entry point for the CLI, not much code here
`go.mod` and `go.sum` - Go package dependencies
`CHANGELOG.md` - Repository Changes by release

`version/version.go` - Isolated package to create a struct for Version information
- It may be possible to roll this into the main package

`cmds/` - Directory for all commands/subcommands
`helper_viper.go` - This file creates an easy set of functions to cleanly retrieve flag values regardless of how they were set
`foo.go` - All the `tfe-mig foo` commands and related functions
