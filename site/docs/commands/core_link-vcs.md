# tfm core link-vcs

## Requirements

- A `terraform_config_metadata.json` must exist in the tfm working directory. Run `tfm core init-repos` to generate one.
- The VCS provider must be configured in TFE/TFC and you must provide the VCS providers Oauth ID as the `vcs_provider_id` in the config file.
- Configure the following credentials in the tfm config file:

```hcl
dst_tfc_hostname="app.terraform.io"
dst_tfc_org="organization"
dst_tfc_token="A token with permissions to create TFC/TFE workspaces"
```

## Link VCS
`tfm core link-vcs` will add the VCS connection to the workspaces that were created using `tfm core create-workspaces`. tfm will use the `terraform_config_metadata.json` file and look at each `config_path` and find a matching TFC/TFE workspace with a name that matches to the path.

tfm will update the workspace settings using the `vcs_provider_id` defined in the config file and the `repo_name` that the `config_path` belongs to and map the repo to the workspace.

## Out of Band Workspace Creation

You can create workspaces using the terraform tfe provider instead of tfm. As long as the workspace names match the constructed workspace name that tfm is looking for then the state will still be uploaded. See the documentation for the `tfm create-workspaces` command for more information regarding workspace name creation.

<!-- ## Cleaning Up

If something goes wrong and you wish to cleanup the workspaces and start the process over you can run the command `tfm nuke workspaces` to delete any workspaces created with tfm commands. -->
