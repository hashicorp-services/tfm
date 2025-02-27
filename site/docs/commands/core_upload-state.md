# tfm core upload-state


## Requirements

- A `terraform_config_metadata.json` must exist in the tfm working directory. Run `tfm core init-repos` to generate one.

## Upload State

`tfm core upload-state` is used to upload the state files that were downloaded using the `tfm core getstate` command to workspace created with the `tfm core create-worksapces` command. tfm will use the `terraform_config_metadata.json` config file to iterate through all of the `config_paths`. Any config path containing a `.terraform/pulled_terraform.tfstate` or `.terraform/pulled_workspaceName_terraform.tfstate` file will have the state file uploaded to a workspace that matches with the `config_path`.

Running this command multiple times will result in the same state file being uploaded multiple times.

## Out of Band Workspace Creation

You can create workspaces using the terraform tfe provider instead of tfm. As long as the workspace names match the constructed workspace name that tfm is looking for then the state will still be uploaded. See the documentation for the `tfm create-workspaces` command for more information regarding workspace name creation.

<!-- ## Cleaning Up

If something goes wrong and you wish to cleanup the workspaces and start the process over you can run the command `tfm nuke workspaces` to delete any workspaces created with tfm commands.
 -->
