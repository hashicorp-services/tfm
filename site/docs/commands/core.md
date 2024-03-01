# Core

`tfm core`

The command used to prefix all commands related to assiting in migrating from terraform open source / community edition (also known as terraform core) to TFC/TFE.

```stdout
tfm core -h
Command used to perform terraform open source (core) to TFE/TFC migration commands

Usage:
  tfm core [command]

Available Commands:
  cleanup           Removes up all cloned repositories from the github_repos_clone_path.
  clone             Clone VCS repositories containing terraform code.
  create-workspaces Create TFE/TFC workspaces for each cloned repo in the github_clone_repos_path that contains a pulled_terraform.tfstate file.
  getstate          Initialize and get state from terraform repos in the github_clone_repos_path.
  init-repos        Scan cloned repositories for Terraform configurations and build metadata
  link-vcs          Link repos in the github_clone_repos_path to their corresponding workspaces in TFE/TFC.
  migrate           Migrates opensource/community edition Terraform code and state to TFE/TFC in 1 continuous workflow.
  remove-backend    Create a branch, remove Terraform backend configurations from cloned repos in github_clone_repos_path, commit the changes, and push to the origin.
  upload-state      Upload .terraform/pulled_terraform.tfstate files from repos cloned into the github_clone_repos_path to TFE/TFC workspaces.

Flags:
  -h, --help   help for core

Global Flags:
      --autoapprove     Auto approve the tfm run. --autoapprove=true . false by default
      --config string   Config file, can be used to store common flags, (default is ~/.tfm.hcl).
      --json            Print the output in JSON format
```

Got an idea for a feature to `tfm`? Submit a [feature request](https://github.com/hashicorp-services/tfm/issues/new?assignees=&labels=&template=feature_request.md&title=)! 