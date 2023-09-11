# tfm

![TFM](site/docs/images/TFM-black.png)

HashiCorp Implementation Services (IS) has identified a need to develop a purpose built tool to assist our engagements and customers during a TFE to TFC migration.

> **Warning**
> This CLI is in beta release and currently does not have official support!
> TFM is currently being developed and tested by interested parties.
> Check out our [case studies](https://hashicorp-services.github.io/tfm/migration/case-studies/). 


## Overview

This tool has been develop to assist HashiCorp Implementation Services and customer engagements during an migration of TFE to TFC(or another TFE). Having a tool allows us the ability to offer a standardized offering to our customers.

Check out the full documentation at [https://hashicorp-services.github.io/tfm/](https://hashicorp-services.github.io/tfm/)

## Installation

Binaries are created as part of a release, check out the [Release Page](https://github.com/hashicorp-services/tfm/releases) for the latest version.

## Pre-Requisites

`tfm` utilize a config file OR environment variables.

### Config File

A HCL file with the following is the minimum located at `/home/user/.tfm.hcl` or specified by `--config config_file`.

```hcl
src_tfe_hostname="tf.local.com"
src_tfe_org="companyxyz"
src_tfe_token="<user token from source TFE/TFC with owner permissions>"
dst_tfc_hostname="app.terraform.io"
dst_tfc_org="companyxyz"
dst_tfc_token="<user token from destination TFE/TFC with owner permissions>"
dst_tfc_project_id="Destination Project ID for workspaces being migrated by tfm. If this is not set, then Default Project is chosen"
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

Alternatively if the source workspaces are not configured to use an agent pool and all destination workspaces should be configured to use an agent pool, a single agent pool ID can be specified instead of an `agents-map` configuration. Note: an `agents-map` config and an `agent-assignment` config can not be specified at the same time.

```hcl
agent-assignment="apool-h896pi2MeP4JJvsB"
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

## Rename Workspaces in destination during a copy

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

### Environment Variables

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

### Delete

This command can be used to delete a resource in either the source or destination side. Today there are two resources that can be deleted

#### workspace

A  workspace can be deleted by the ID or name. This does not look for the hidden tfm source on workspaces that have been created by tfm, so this will let you delete any workspace.
Note: This command will ask for confirmation, however this is a destructive operation and the workspace can not be recovered.

#### Workspaces VCS Connection

This command will look at the `.tfm` configuration for the workspaces on the source that should have their VCS connection removed. This would be utilized after a migration can been completed, and the existing VCS connections on the source should be removed from the workspaces, so runs are no longer triggered by VCS
Note: This command will ask for confirmation. The source workspaces will not be deleted, only their VCS connection, which can be added back manually if needed.

### Nuke

If workspaces that have been created in the destination organization need to be destroyed, `tfm nuke workspace` can be used to remove all workspaces that tfm created. This is done by listing all workspaces in the destination organization and checking if the `SourceName` is set to `tfm`. This command will prompt for confirmation. If Confirmed tfm will delete the workspaces. This is a destructive operation and the workspaces can not be recovered.

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
