# tfm core migrate


## Requirements

- The VCS provider must be configured in TFE/TFC and you must provide the VCS providers Oauth ID as the `vcs_provider_id` in the config file.
- Configure the `clone_repos_path` in the config file.
- Configure the `vcs_type` with a [supported vcs types](../migration/supported-vcs.md) in the config file.
- Authentication credentials for the cloned terraform configuration backends must be configured in the environment.
- Terraform CLI must be installed in the environment and in the path.
- Configure the following credentials in the tfm config file:

```hcl
dst_tfc_hostname="app.terraform.io"
dst_tfc_org="organization"
dst_tfc_token="A token with permissions to create TFC/TFE workspaces"
```

- Configure the VCS credentials in the config file required for your [supported vcs type](../migration/supported-vcs.md)

`tfm core migrate` will sequentially run all of the commands required to migrate terraform open source / community edition configurations to TFE/TFC workspace management.

tfm will run the following commands in the following order when the migrate command is used:

`tfm core clone`
`tfm core init-repos`
`tfm core getstate`
`tfm core create-worksapces`
`tfm core upload-state`
`tfm core link-vcs`

## Flags

`--include remove-backend` will add the `tfm core remove-backend` command to be run last as part of the `tfm core migrate` command. This requires a VCS API token with write permissions to the VCS repositories.

<!-- ## Cleaning Up

If something goes wrong and you wish to cleanup the workspaces and start the process over you can run the command `tfm nuke workspaces` to delete any workspaces created with tfm commands. -->
