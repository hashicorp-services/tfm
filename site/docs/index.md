# Welcome to TFM Docs

![TFM](./images/TFM-black.png)

!!! warning ""
    Note: This CLI currently does not have official support!

    TFM is currently being developed and tested by interested parties.
    
    The following is a [summary of a real example use case from a real TFE customer](https://docs.google.com/presentation/d/e/2PACX-1vSFN9osZARMitvG8HcbnR37nSbVXnK2GMlMvXOz7GsNceFDYp4-98Ko4xZ89-Rtvkf1_YqBmx338er3/pub?start=true&loop=true&delayms=3000). 
    
    Please reach out to one of our listed [contact methods](#contacts).

_tfm_ is a standalone CLI for Terraform Community Edition, Terraform Cloud, Terraform Enterprise migrations.

HashiCorp Implementation Services (IS) has identified a need to develop a purpose built tool to assist our engagements, partners, and customers during:

- Terraform open source / community edition / core to TFC/TFE
- TFE to TFC
- TFC to TFE
- 1 TFC Organization to another TFC Organization

Note: The Terraform Community Edition migration as part of `tfm` has been deprecated in favor of [tf-migrate](https://developer.hashicorp.com/terraform/cloud-docs/migrate/tf-migrate). The CE migration feature has not been removed from `tfm` however it will not be receiving further developments.

## Installation

Binaries are created as part of a release, check out the [Release Page](https://github.com/hashicorp-services/tfm/releases) for the latest version.

### MacOs Installation amd64

```sh
version="x.x.x"
curl -L -o tfm "https://github.com/hashicorp-services/tfm/releases/download/${version}/tfm_darwin_x86_64"
chmod +x tfm
```

### MacOs Installation arm64

```sh
version="x.x.x"
curl -L -o tfm "https://github.com/hashicorp-services/tfm/releases/download/${version}/tfm_darwin_arm64"
chmod +x tfm
```

!!! note ""
    Note: `tfm` CLI is currently not developer-signed or notarized and you will run into an initial issue where `tfm` is not allowed to run. Please follow "[safely open apps on your mac](https://support.apple.com/en-au/HT202491#:~:text=View%20the%20app%20security%20settings%20on%20your%20Mac&text=In%20System%20Preferences%2C%20click%20Security,%E2%80%9CAllow%20apps%20downloaded%20from.%E2%80%9D)" to allow `tfm` to run

### Linux Installation

```sh
  version="x.x.x"
  curl -L -o tfm "https://github.com/hashicorp-services/tfm/releases/download/${version}/tfm_linux_x86_64"
  chmod +x tfm
```

### Windows Installation

```sh
  version="x.x.x"
  curl -L -o tfm.exe "https://github.com/hashicorp-services/tfm/releases/download/${version}/tfm_windows_x86_64.exe"
```

## Usage

### Print the TFM CLI help

`tfm -h`

```bash
tfm -h
A CLI to assist with Terraform community edition, Terraform Cloud, and Terraform Enterprise migrations.

Usage:
  tfm [command]

Available Commands:
  copy        Copy command
  core        Command used to perform terraform open source (core) to TFE/TFC migration commands
  delete      delete command
  generate    generate command for generating .tfm.hcl config template
  help        Help about any command
  list        List command
  lock        Lock
  unlock      Unlock

Flags:
      --autoapprove     Auto approve the tfm run. --autoapprove=true . false by default
      --config string   Config file, can be used to store common flags, (default is ~/.tfm.hcl).
  -h, --help            help for tfm
      --json            Print the output in JSON format
  -v, --version         version for tfm

Use "tfm [command] --help" for more information about a command.
```

## Pre-Requisites

The following prerequisites are used when migrating from or to TFE or TFC from TFE or TFC.

- A tfm config file
- A TFC/TFE token with for the source TFE/TFC Organization that you are migrating from.
- A TFC/TFE token with for the source TFE/TFC Organization that you are migrating to.

For which token type to use, we recommend either Team or Personal token. We do not recommend Organization token due to the limited permissions Org tokens have. Please refer to [Access Levels](https://developer.hashicorp.com/terraform/cloud-docs/users-teams-organizations/api-tokens#access-levels) for the different Token types and their permissions.

### Environment Variables

The following environment variables can be set or used to override existing config file values.

```bash
export SRC_TFE_HOSTNAME="tf.local.com"
export SRC_TFE_ORG="companyxyz"
export SRC_TFE_TOKEN="<user token from source TFE/TFC with permissions>"
export DST_TFC_HOSTNAME="app.terraform.io"
export DST_TFC_ORG="companyxyz"
export DST_TFC_TOKEN="<user token from source TFE/TFC with permissions>"
export DST_TFC_PROJECT_ID="Destination Project ID for workspaces being migrated by tfm. If this is not set, then Default Project is chosen"
```

### Config File

A HCL file with the following as the minimum located at `/home/user/.tfm.hcl` or specified by `--config config_file`. You can also run `tfm generate config` to create a tempalte config file for use.

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

### Copy Workspaces into Projects

By default, a workspace will be copied over to the Default Project in the destination (eg TFC).
Users can specify the project ID for the desired project to place all workspaces in the `tfm copy workspace` run.

Utilize `tfm list projects --side destination` to determine the `project id`.

Set either the environment variable:

```bash
export DST_TFC_PROJECT_ID=prj-XXXX
```

or specify the following in your `~/.tfm.hcl` configuration file.

```hcl
dst_tfc_project_id=prj-xxx 
```

## Workspace List

As part of the HCL config file (`/home/user/.tfm.hcl`), a list of workspaces from the source TFE can be specified. `tfm` will use this list when running `tfm copy workspaces` and ensure the workspace exists or is created in the target.

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
_NOTE: Using this configuration in your HCL config file will take precedence over the other Workspace List which only lists source workspace names._

```hcl
# A list of source=destination workspace names. TFM will look at each source workspace and recreate the workspace with the specified destination name.
"workspaces-map" = [
   "tf-demo-workflow=dst-demo-workflow",
   "api-test=dst-api-test"
   ]
```

## Assign SSH

As part of the HCL config file (`/home/user/.tfm.hcl`), a list of `source-ssh-key-id=destination-ssh-key-id` can be provided. `tfm` will use this list when running `tfm copy workspaces --ssh` to look at all workspaces in the source host with the assigned source SSH key ID and assign the matching named workspace in the destination with the mapped SSH key ID.

```hcl
# A list of source=destination SSH IDs. TFM will look at each workspace in the source for the source SSH  ID and assign the matching workspace in the destination with the destination SSH ID.
ssh-map=[
  "sshkey-sPLAKMcqnWtHPSgx=sshkey-CRLmPJpoHwsNFAoN",
]
```

## Contacts

- Initial Slack Channel for developement: [#ps-offering-tfm](https://hashicorp.slack.com/archives/C046STDBXNC)
- Google Group TFM Dev Team: [svc-github-team-tfm@hashicorp.com](svc-github-team-tfm@hashicorp.com
)
- Got an idea for a feature to `tfm`? Submit a [feature request](https://github.com/hashicorp-services/tfm/issues/new?assignees=&labels=&template=feature_request.md&title=) or provide us some [feedback](./feedback.md).
