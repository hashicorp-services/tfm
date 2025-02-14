# tfm core getstate

## Requirements

- Terraform community edition must be installed and in your path in the environment where this command is run.
- Credentials to authenticate to the configured backend must be configured in the environment where this command is run.
- A `terraform_config_metadata.json` must exist in the tfm working directory. Run `tfm core init-repos` to generate one.

## Get State

`tfm core getstate` will use the `terraform_config_metadata.json` config file to  iterate through all of the cloned repositories in the `clone_repos_path` and metadata `config_paths` to download the state files from the backend.

tfm will use the locally installed terraform binary to perform `terraform init` and `terraform state pull > .terraform/pulled_terraform.tfstate` commands.

If tfm cannot successfully run a `terraform init` for a cloned repo tfm will return an error and continue with the next repository initilization attempt.

## Terraform CE Workspaces

For any `config_path` with `uses_workspaces: true`, tfm will run `tfm workspace select` for each workspace in the `workspace_names` list and `terraform state pull > .terraform/pulled_<worspace name>_terraform.tfstate`. The end result will be multiple state files within the `config_path` for each workspace.
