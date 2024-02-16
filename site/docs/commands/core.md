# Core

`tfm core`

The command used to prefix all commands related to assiting in migrating from terraform open source / community edition (also known as terraform core) to TFC/TFE.

```stdout
tfm core -h 
Command used to perform terraform open source (core) to TFE/TFC migration commands

Usage:
  tfm core [command]

Available Commands:
  cleanup           Cleans up all repositories in the clone path
  clone             Clone VCS repositories containing terraform code.
  create-workspaces Create TFE workspaces for each cloned Terraform repo
  getstate          Initialize and get state from terraform VCS repos.
  link-vcs          Link repos to TFE workspaces via VCS
  migrate           Migrates opensource/community edition Terraform code and state to TFE/TFC.
  remove-backend    Remove Terraform backend configurations
  upload-state      Upload pulled_terraform.tfstate files to TFE workspaces

Flags:
  -h, --help   help for core
```

## Copy sub commands

- [`tfm copy teams`](copy_teams.md)
- [`tfm copy varsets`](copy_varsets.md)
- [`tfm copy workspaces`](copy_workspaces.md)


## Possible Future copy commands enhancements

- `tfm copy modules`
- `tfm copy policy-sets`
- `tfm copy workspace --all`

Got an idea for a feature to `tfm`? Submit a [feature request](https://github.com/hashicorp-services/tfm/issues/new?assignees=&labels=&template=feature_request.md&title=)! 