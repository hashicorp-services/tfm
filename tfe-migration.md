# TFM - Terraform Enterprise to or from Terraform Cloud

## Pre-Requisites

The following prerequisites are used when migrating from or to TFE or TFC from TFE or TFC.

- A [tfm config file](./site/docs/configuration_file/config_file.md)
- A TFC/TFE Owner token with for the source TFE/TFC Organization that you are migrating from
- A TFC/TFE Owner token with for the source TFE/TFC Organization that you are migrating to

## Config File

`tfm` utilizes a [config file](./site/docs/configuration_file/config_file.md) OR environment variables. An HCL file with the following is the minimum located at `/home/user/.tfm.hcl` or specified by `--config /path/to/config_file`. Multiple config files can be created to assist with large migrations.

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

## Environment Variables

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

As part of the HCL config file (`/home/user/.tfm.hcl`), a list of vcs mappings of either `source-vcs-oauth-id=destination-vcs-oauth-id` or `source-vcs-github-app-id=destination-vcs-github-app-id` can be provided. `tfm` will use this list when running `tfm copy workspaces --vcs` to look at all workspaces in the source host with the assigned source VCS ID and assign the matching named workspace in the destination with the mapped destination VCS ID. `tfm` only supports like for like vcs migration, so if the source is a GitHub App VCS connection the destination must use a GitHub App VCS connection.

```hcl
# A list of source=destination VCS IDs. TFM will look at each workspace in the source for the source VCS ID and assign the matching workspace in the destination with the destination VCS ID.
vcs-map=[
  "ot-5uwu2Kq8mEyLFPzP=ot-coPDFTEr66YZ9X9n",
  "ot-gkj2An452kn2flfw=ot-8ALKBaqnvj232GB4",
  "ghain-sc8a3b12S212gy45=ghain-B3asgvX3oF541aDo"
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
   ]
```

## Copy Workspaces into Projects

By default, a workspace will be copied over to the Default Project in the destination (eg TFC).
Users can specify the project ID for the desired project to place all workspaces in the `tfm copy workspace` run.

Utilize `tfm list projects --side destination` to determine the `project id`.

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
