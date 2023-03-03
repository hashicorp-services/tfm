# Welcome to TFM Docs

![TFM](../../images/TFM-black.png)

!!! warning ""
    Note: This CLI is still under active development and subject to change!

_tfm_ is a standalone CLI for Terraform Cloud and Terraform Enterprise migrations.

HashiCorp Implementation Services (IS) has identified a need to develop a purpose built tool to assist our engagements and customers during a TFE to TFC migration.

This tool has been develop to assist HashiCorp Implementation Services and customer engagements during an migration of TFE to TFC(or another TFE). Having a tool allows us the ability to offer a standardized offering to our customers.



## Installation

Binaries are created as part of a release, check out the [Release Page](https://github.com/hashicorp-services/tfm/releases) for the latest version.

**MacOs Installation**
```sh
version="x.x.x"
curl -L -o tfm "https://github.com/hashicorp-services/tfm/download/${version}/tfm_darwin_amd64"
chmod +x tfm
```

**Linux Installation**
```sh
version="x.x.x"
curl -L -o tfm "https://github.com/hashicorp-services/tfm/download/${version}/tfm_linux_amd64"
chmod +x tfm
```

**Windows Installation**
```sh
version="x.x.x"
curl -L -o tfm.exe "https://github.com/hashicorp-services/tfm/download/${version}/tfm_windows_amd64"
```

## Usage

### Print the TFM CLI help

`tfm -h`

```
A CLI to assist with TFE Migration.

Usage:
  tfm [command]

Available Commands:
  copy        Copy command
  help        Help about any command
  list        List command

Flags:
      --config string   Config file, can be used to store common flags, (default is ./.tfm.hcl).
  -h, --help            help for tfm
  -v, --version         version for tfm

Use "tfm [command] --help" for more information about a command.
```

## Pre-Requisites

`tfm` utilize a config file AND/OR environment variables.
We recommend using environment variables for sensitive tokens from TFE/TFC instead of storing it in a config file. 

### Environment Variables

The following environment variables can be set or used to override existing config file values.

```bash
export SOURCEHOSTNAME="tf.local.com"
export SOURCEORGANIZATION="companyxyz"
export SOURCETOKEN="<user token from source TFE/TFC with owner permissions>"
export DESTINATIONHOSTNAME="app.terraform.io"
export DESTINATIONORGANIZATION="companyxyz"
export DESTINATIONTOKEN="<user token from source TFE/TFC with owner permissions>"
```

### Config File

A HCL file with the following is the minimum located at `/home/user/.tfm.hcl` or specified by `--config config_file`.

```hcl
sourceHostname="tf.local.com"
sourceOrganization="companyxyz"
sourceToken="<user token from source TFE/TFC with owner permissions>"
destinationHostname="app.terraform.io"
destinationOrganization="companyxyz"
destinationToken="<user token from destination TFE/TFC with owner permissions>"
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
*NOTE: Using this configuration in your HCL config file will take precedence over the other Workspace List which only lists source workspace names.*

```hcl
# A list of source=destination workspace names. TFM will look at each source workspace and recreate the workspace with the specified destination name.
"workspace-map" = [
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


## TFM Demo

<iframe src="https://drive.google.com/file/d/1IG8B3yqPFoX642vYnE13lFLPAoUd7Wz6/preview" width="640" height="480" allow="autoplay"></iframe>


## Contacts

- Initial Slack Channel for developement: [#ps-offering-tfe-migration](https://hashicorp.slack.com/archives/C046STDBXNC)
- Google Group TFM Dev Team: [svc-github-team-tfm@hashicorp.com](svc-github-team-tfm@hashicorp.com
)

