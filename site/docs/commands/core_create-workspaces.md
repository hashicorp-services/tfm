# tfm core create-workspaces

## Requirements

- A `terraform_config_metadata.json` must exist in the tfm working directory. Run `tfm core init-repos` to generate one.
- Configure the following credentials in the tfm config file:

```hcl
dst_tfc_hostname="app.terraform.io"
dst_tfc_org="organization"
dst_tfc_token="A token with permissions to create TFC/TFE workspaces"
```

## Create Workspaces

`tfm core create-workspaces` will use the `terraform_config_metadata.json` config file to create a TFC/TFE workspace in the `dst_tfc_org` defined in the config file for each repository.

Workspace names are generated using the metadata file in the following format:

- Config path using terraform ce workspaces:
`repo_name+config_path+workspace_name`

- Config path without terraform ce workspaces:
`repo_name+config_path`

As an example, the below metadata would create 2 TFC/TFE workspaces with the names:

- `isengard-infra-east-primary-newisengard`
- `isengard-infra-east-primary-oldisengard`

```json
  {
    "repo_name": "isengard",
    "config_paths": [
      {
        "path": "isengard/infra/east/primary",
        "workspace_info": {
          "uses_workspaces": true,
          "workspace_names": [
            "default",
            "newisengard",
            "oldisengard"
          ]
        }
      },
```

If a workspace already exists with the name of the repository tfm will return an error and continue on to the next workspace creation attempt.

<!-- ## Cleaning Up

If something goes wrong and you wish to cleanup the workspaces and start the process over you can run the command `tfm nuke workspaces` to delete any workspaces created with tfm commands. -->

## Out of Band Workspace Creation

You can create workspaces using the terraform tfe provider instead of tfm. As long as the workspace names match the constructed workspace name that tfm is looking for then the state will still be uploaded.

## Future Updates

Future updates are planned to also allow a variable set containing credentials to be assigned to the workspace the the time of creation. This will allow a plan to be run post migration to verify no changes are expected.
