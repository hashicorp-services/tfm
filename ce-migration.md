# TFM - Terraform Open Source / Community Edition to TFC/TFE

## Pre-Requisites

The following prerequisites are used when migrating from terraform community edition (also known as open source) to TFC/TFE managed workspaces.

- Terraform - Must be installed in the execution environment and available in the path
- A configuration file
- A terraform cloud or enterprise token with the permissions to create workspaces in an organization
- A Github token with the permissions to read repositories containing terraform code to be migrated

## Config File

`tfm` utilizes a config file OR environment variables. An HCL file with the following is the minimum located at `/home/user/.tfm.hcl` or specified by `--config /path/to/config_file`. Multiple config files can be created to assist with large migrations.

> [!NOTE]
> Use the `tfm generate config` command to generate a sample configuration for quick editing.

```hcl
dst_tfc_hostname="app.terraform.io for TFC or the hostname of your TFE application"
dst_tfc_org="A TFE/TFC organization to create workspaces in"
dst_tfc_token="A TFC/TFE Token with the permissions to create workspaces in the TFC/TFE organization"
github_token = "A Github token with the permissions to read terraform code repositories you wish to migrate"
github_organization = "The github organization containing terrafor code repositories"
github_username = "A github username"
github_clone_repos_path = "/path/on/local/system/to/clone/repos/to"
```

Additional configurations can be provided to assist in the community edition to TFC/TFE migration:

```hcl
commit_message = "A commit message the tfm core remove-backend command uses when removing backend blocks from .tf files and commiting the changes back"
commit_author_name = "the name that will appear as the commit author"
commit_author_email = "the email that will appear for the commit author"
vcs_provider_id = "An Oauth ID of a VCS provider connection configured in TFC/TFE"

# A list of VCS repositories containing terraform code. TFM will clone each repo during the tfm core clone command for migrating opensource/commmunity edition terraform managed code to TFE/TFC.

Organization.
repos_to_clone =  [
 "repo1",
 "repo2",
 "repo3"
]
```

## Environment Variables

If no config file is found, the following environment variables can be set or used to override existing config file values.

```bash
export dst_tfc_hostname="app.terraform.io for TFC or the hostname of your TFE application"
export dst_tfc_org="A TFE/TFC organization to create workspaces in"
export dst_tfc_token="A TFC/TFE Token with the permissions to create workspaces in the TFC/TFE organization"
export github_token = "A Github token with the permissions to read terraform code repositories you wish to migrate"
export github_organization = "The github organization containing terrafor code repositories"
export github_username = "A github username"
export github_clone_repos_path = "/path/on/local/system/to/clone/repos/to"
```

## tfm in a Pipeline

`tfm` can be used in a pipeline to automate migrations. There are a few considerations when using `tfm` in this manner.

- For source and destination credentials use `SRC_TFE_TOKEN` and `DST_TFC_TOKEN` environment variables from the pipeline's secrete manager or other secure means. Don't embed these credentials in `tfm` configuration files
- Several `tfm` commands require confirmation before proceeding, which are listed below. To override these in a pipeline, add the `--autoapprove` flag.
  - `copy workspaces` Only when all workspaces are going to be migrated due to no workspace list, or map defined.
  - `copy workspaces --state --last`
  - `delete workspace`
  - `delete workspaces-vcs`
  - `nuke`
