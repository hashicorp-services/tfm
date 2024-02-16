# tfm

![TFM](site/docs/images/TFM-black.png)

HashiCorp Implementation Services (IS) has identified a need to develop a purpose built tool to assist our engagements and customers during:

- Terraform open source / community edition / core to TFC/TFE
- TFE to TFC
- TFC to TFE
- 1 TFC Organization to another TFC Organization

> **Warning**
> This CLI does not have official support, but the code owners will work with partners and interested parties to provide assitance when possible.
> Check out our [case studies](https://hashicorp-services.github.io/tfm/migration/case-studies/). 


## Overview

This tool has been develop to assist HashiCorp Implementation Services and customer engagements during an migration of Terraform open source to TFE or TFC or with a migration from TFE to TFC, TFC to TFE, or 1 TFC organization to another TFC organization. Having a tool allows us the ability to offer a standardized offering to our customers.

Check out the full documentation at [https://hashicorp-services.github.io/tfm/](https://hashicorp-services.github.io/tfm/)

## Installation

Binaries are created as part of a release, check out the [Release Page](https://github.com/hashicorp-services/tfm/releases) for the latest version.

## Pre-Requisites - TFE to TFC, TFC to TFE, or TFC to TFC Migrations

The following prerequisites are used when migrating from or to TFE or TFC from TFE or TFC.

- A tfm config file
- A TFC/TFE Owner token with for the source TFE/TFC Organization that you are migrating from
- A TFC/TFE Owner token with for the source TFE/TFC Organization that you are migrating to

## Config File - TFE to TFC, TFC to TFE, or TFC to TFC Migrations

`tfm` utilizes a config file OR environment variables. An HCL file with the following is the minimum located at `/home/user/.tfm.hcl` or specified by `--config /path/to/config_file`. Multiple config files can be created to assist with large migrations.

> [!NOTE]
> Use the `tfm generate config` command to generate a sample configuration for quick editing.


```hcl
src_tfe_hostname="tf.local.com"
src_tfe_org="companyxyz"
src_tfe_token="<user token from source TFE/TFC with owner permissions>"
dst_tfc_hostname="app.terraform.io"
dst_tfc_org="companyxyz"
dst_tfc_token="<user token from destination TFE/TFC with owner permissions>"
dst_tfc_project_id="Destination Project ID for workspaces being migrated by tfm. If this is not set, then Default Project is chosen"
```

## Environment Variables - TFE to TFC, TFC to TFE, or TFC to TFC Migrations

If no config file is found, the following environment variables can be set or used to override existing config file values.

```bash
export SRC_TFE_HOSTNAME="tf.local.com"
export SRC_TFE_ORG="companyxyz"
export SRC_TFE_TOKEN="<user token from source TFE/TFC with owner permissions>"
export DST_TFC_HOSTNAME="app.terraform.io"
export DST_TFC_ORG="companyxyz"
export DST_TFC_TOKEN="<user token from source TFE/TFC with owner permissions>"
export DST_TFC_PROJECT_ID="Destination Project ID for workspaces being migrated by tfm. If this is not set, then Default Project is chosen"
```

## Workspace List

As part of the HCL config file (`/home/user/.tfm.hcl`), a list of workspaces from the source TFE can be specified. `tfm` will use this list when running `tfm copy workspaces` and ensure the workspace exists or is created in the target with the same name.

```hcl
#List of Workspaces to create/check are migrated across to new TFC
"workspaces" = [
  "appAFrontEnd",
  "appABackEnd",
  "appBDataLake",
  "appBInfra"
]

```

## Assign Agent Pools to Workspaces

As part of the HCL config file (`/home/user/.tfm.hcl`), a list of `source-agent-pool-ID=destination-agent-pool-ID` can be provided. `tfm` will use this list when running `tfm copy workspaces --agents` to look at all workspaces in the source host with the assigned source agent pool ID and assign the matching named workspace in the destination with the mapped destination agent pool ID.

```hcl
# A list of source=destination agent pool IDs TFM will look at each workspace in the source for the source agent pool ID and assign the matching workspace in the destination the destination agent pool ID.
agents-map = [
  "apool-DgzkahoomwHsBHcJ=apool-vbrJZKLnPy6aLVxE",
  "apool-DgzkahoomwHsBHc3=apool-vbrJZKLnPy6aLVx4",
  "apool-DgzkahoomwHsB125=apool-vbrJZKLnPy6adwe3",
  "test=beep"
]
```

Alternatively if the source workspaces are not configured to use an agent pool and all destination workspaces should be configured to use an agent pool, a single agent pool ID can be specified instead of an `agents-map` configuration. Note: an `agents-map` config and an `agent-assignment-id` config can not be specified at the same time.

```hcl
agent-assignment-id="apool-h896pi2MeP4JJvsB"
```

## Copy Variable Sets

To copy ALL variable sets from the source to the destination run the command:
`tfm copy varsets`

To copy only desired variable sets, provide an HCL list in the `.tfm.hcl` configuration file using the snyntax `"source-varset-name=destination-varset-name"`. This list will be converted to a map. tfm will copy only the source variable sets provided on the left side of the `=`. The right side of the `=` can optionally be a different name to allow you to copy the variable set with a new name. Both sides of the `=` must be populated and `varsets-map` cannot be empty if it is defined.

Example configuration file:

```hcl
varsets-map = [
  "Azure-creds=New-Azure-Creds",
  "aws-creds2=New-AWS-Creds",
  "SourceVarSet=DestVarSet"
 ]
 ```

## Assign VCS

As part of the HCL config file (`/home/user/.tfm.hcl`), a list of `source-vcs-oauth-ID=destination-vcs-oauth-id-ID` can be provided. `tfm` will use this list when running `tfm copy workspaces --vcs` to look at all workspaces in the source host with the assigned source VCS oauth ID and assign the matching named workspace in the destination with the mapped destination VCS oauth ID.

```hcl
# A list of source=destination VCS oauth IDs. TFM will look at each workspace in the source for the source VCS oauth ID and assign the matching workspace in the destination with the destination VCS oauth ID.
vcs-map=[
  "ot-5uwu2Kq8mEyLFPzP=ot-coPDFTEr66YZ9X9n",
  "ot-gkj2An452kn2flfw=ot-8ALKBaqnvj232GB4",

]
```

## Rename Workspaces in Destination During a Copy

As part of the HCL config file (`/home/user/.tfm.hcl`), a list of `source-workspace-name=destination-workspace-name` can be provided. `tfm` will use this list when running `tfm copy workspace` to look at all workspaces in the source host and rename the destination workspace name.
*NOTE: Using this configuration in your HCL config file will take precedence over the other Workspace List which only lists source workspace names.*

```hcl
# A list of source=destination workspace names. TFM will look at each source workspace and recreate the workspace with the specified destination name.
"workspaces-map" = [
   "tf-demo-workflow=dst-demo-workflow",
   "api-test=dst-api-test"
```

## Copy Workspaces into Projects

By default, a workspace will be copied over to the Default Project in the destination (eg TFC).
Users can specify the project ID for the desired project to place all workspaces in the `tfm copy workspace` run.

Utilise `tfm list projects --side destination` to determine the `project id`.

Set either the environment variable:

```bast
export DST_TFC_PROJECT_ID=prj-XXXX
```

or specify the following in your `~/.tfm.hcl` configuration file.

```terraform
dst_tfc_project_id=prj-xxx 
```

## Assign SSH

As part of the HCL config file (`/home/user/.tfm.hcl`), a list of `source-ssh-key-id=destination-ssh-key-id` can be provided. To obtain the ssh-ids the `tfm list ssh --side=[source/destination]` command can be used.

`tfm` will use this list when running `tfm copy workspaces --ssh` to look at all workspaces in the source host with the assigned source SSH key ID and assign the matching named workspace in the destination with the mapped SSH key ID.

```hcl
# A list of source=destination SSH IDs. TFM will look at each workspace in the source for the source SSH  ID and assign the matching workspace in the destination with the destination SSH ID.
ssh-map=[
  "sshkey-sPLAKMcqnWtHPSgx=sshkey-CRLmPJpoHwsNFAoN",
]
```

## Additional Commands

### Delete

This command can be used to delete a resource in either the source or destination side. Today there are two resources that can be deleted

### Workspace

A  workspace can be deleted by the ID or name. This does not look for the hidden tfm source on workspaces that have been created by tfm, so this will let you delete any workspace.
Note: This command will ask for confirmation, however this is a destructive operation and the workspace can not be recovered.

### Workspaces VCS Connection

This command will look at the `.tfm` configuration for the workspaces on the source that should have their VCS connection removed. This would be utilized after a migration can been completed, and the existing VCS connections on the source should be removed from the workspaces, so runs are no longer triggered by VCS
Note: This command will ask for confirmation. The source workspaces will not be deleted, only their VCS connection, which can be added back manually if needed.

### Nuke

If workspaces that have been created in the destination organization need to be destroyed, `tfm nuke workspace` can be used to remove all workspaces that tfm created. This is done by listing all workspaces in the destination organization and checking if the `SourceName` is set to `tfm`. This command will prompt for confirmation. If Confirmed tfm will delete the workspaces. This is a destructive operation and the workspaces can not be recovered.

### Lock & Unlock

The `tfm lock workspaces` & `tfm unlock workspaces` commands can be used to lock and unlock a workspace in either the source or destination as needed. This will use the workspaces as configured in the `tfm` config file and either lock them or unlock them. If a workspace is already locked it will skip trying to lock the workspace, and same for the unlock command. It will default to the source side, however with `--side destination` it will lock or unlock the destination side.

## Pre-Requisites - Terraform Open Source / Community Edition to TFC/TFE

The following prerequisites are used when migrating from terraform community edition (also known as open source) to TFC/TFE managed workspaces.

- Terraform - Must be installed in the execution environment and avaible in the path
- A configuration file
- A terraform cloud or enterprise token with the permissions to create worksapces in an organization
- A Github token with the permissions to read repositories containing terraform code to be migrated

## Config File - Terraform Open Source / Community Edition to TFC/TFE

`tfm` utilizes a config file OR environment variables. An HCL file with the following is the minimum located at `/home/user/.tfm.hcl` or specified by `--config /path/to/config_file`. Multiple config files can be created to assist with large migrations.

> [!NOTE]
> Use the `tfm generate config` command to generate a sample configuration for quick editing.

```
dst_tfc_hostname="app.terraform.io for TFC or the hostname of your TFE application"
dst_tfc_org="A TFE/TFC organization to create workspaces in"
dst_tfc_token="A TFC/TFE Token with the permissions to create workspaces in the TFC/TFE organization"
github_token = "A Github token with the permissions to read terraform code repositories you wish to migrate"
github_organization = "The github organization containing terrafor code repositories"
github_username = "A github username"
github_clone_repos_path = "/path/on/local/system/to/clone/repos/to"
```

Additional configurations can be provided to assist in the community edition to TFC/TFE migration:

```
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
## Environment Variables - Terraform Open Source / Community Edition to TFC/TFE

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


## Architectural Decisions Record (ADR)

An architecture decision record (ADR) is a document that captures an important architecture decision made along with its context and consequences.

This project will store ADRs in [docs/ADR](docs/ADR/) as a historical record.

More information about [ADRs](docs/ADR/index.md).

## To build

```bash
make build-local
./tfm -v
```

-or-

```bash
go run . -v
```

## To release

To create a new release of TFM

- Use GitHub Action `Release tfm`
- Specify a version number. Please follow semantic versioning for the release.

This action will do the following steps

- Compile TFM for Linux, Mac, Windows with amd64 and arm64 versions
- Upload the artifacts
- Create a new release + tag on the repo at the current main

## Reporting Issues

If you believe you have found a defect in `tfm` or its documentation, use the [GitHub issue tracker](https://github.com/hashicorp-services/tfm/issues) to report the problem to the `tfm` maintainers.
