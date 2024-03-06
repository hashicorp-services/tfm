# TFM - Terraform Open Source / Community Edition to TFC/TFE

## Pre-Requisites

The following prerequisites are used when migrating from terraform community edition (also known as open source) to TFC/TFE managed workspaces.

- Terraform - Must be installed in the execution environment and available in the path
- A [configuration file](./site/docs/configuration_file/config_file.md)
- Terraform backend authentication credentials must be configured in the execution environment.
- A terraform cloud or enterprise token with the permissions to create workspaces in an organization
- A [supported VCS](./site/docs/migration/supported-vcs.md) token with the permissions to read repositories containing terraform code to be migrated

## Config File

`tfm` utilizes [a config file](./site/docs/configuration_file/config_file.md) OR environment variables. An HCL file with the following is the minimum located at `/home/user/.tfm.hcl` or specified by `--config /path/to/config_file`. Multiple config files can be created to assist with large migrations.

> [!NOTE]
> Use the `tfm generate config` command to generate a sample configuration for quick editing.

```hcl
dst_tfc_hostname="app.terraform.io for TFC or the hostname of your TFE application"
dst_tfc_org="A TFE/TFC organization to create workspaces in"
dst_tfc_token="A TFC/TFE Token with the permissions to create workspaces in the TFC/TFE organization"
vcs_type = " A [supported vcs_type](./site/docs/migration/supported-vcs.md) "
github_token = "A Github token with the permissions to read terraform code repositories you wish to migrate"
github_organization = "The github organization containing terraform code repositories"
github_username = "A github username"
gitlab_username = "A gitlab username"
gitlab_token = "A gitlab token"
gitlab_group = "A gitlab group"
clone_repos_path = "/path/on/local/system/to/clone/repos/to"
```

Additional configurations can be provided to assist in the community edition to TFC/TFE migration:

```hcl
commit_message = "A commit message the tfm core remove-backend command uses when removing backend blocks from .tf files and committing the changes back"
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
export github_token="A Github token with the permissions to read terraform code repositories you wish to migrate"
export github_organization="The github organization containing terraform code repositories"
export github_username="A github username"
export clone_repos_path="/path/on/local/system/to/clone/repos/to"
export vcs_type="A [supported VCS](./site/docs/migration/supported-vcs.md)"
export gitlab_username="A gitlab username"
export gitlab_token="A gitlab token"
export gitlab_group="A gitlab group"
```
