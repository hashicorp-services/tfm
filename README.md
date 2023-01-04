# tfm

HashiCorp IS has identified a need to develop a purpose built tool to assist our engagments and customers during a TFE to TFC migration. 

## Overview

This tool has been develop to assist HashiCorp Implementation Services and customer engagements during an migration of TFE to TFC(or another TFE). Having a tool allows us the ability to offer a standardized offering to our customers. 

## Pre-Requisites

`tfm` utilize a config file OR environment variables.

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

<<<<<<< HEAD
## Workspace List
As part of the HCL config file (`/home/user/.tfx.hcl`), a list of workspaces from the source TFE can be specified. `tfemigrate` will use this list when running `tfemigrate copy workspaces` and ensure the workspace exists or is created in the target. 

```hcl
#List of Workspaces to create/check are migrated across to new TFC
"workspaces" = [
  "appAFrontEnd",
  "appABackEnd",
  "appBDataLake",
  "appBInfra"
]

```

=======
### Environment Variables
>>>>>>> main

If no config file is found, the following environment variables can be set or used to override existing config file values.

```bash
export SOURCEHOSTNAME="tf.local.com"
export SOURCEORGANIZATION="companyxyz"
export SOURCETOKEN="<user token from source TFE/TFC with owner permissions>"
export DESTINATIONHOSTNAME="app.terraform.io"
export DESTINATIONORGANIZATION="companyxyz"
export DESTINATIONTOKEN="<user token from source TFE/TFC with owner permissions>"
```

## Docs

Check out our documentation page (coming soon)

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

## Reporting Issues

If you believe you have found a defect in `tfm` or its documentation, use the [GitHub issue tracker](https://github.com/hashicorp-services/tfm/issues) to report the problem to the `tfm` maintainers.
