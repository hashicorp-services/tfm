# TFM FAQs

## Who is `tfm` developed for?

Engineers/Operators that manage/admin Terraform Enterprise/Cloud organizations that need to perform a migration of workspaces. 

## What is `tfm` intended to do?

`tfm` will assist with migration of TFE/TFC workspaces from one TFE/TFC instance/organization to another TFE/TFC instance/organization. 

From a customer journey perspective, it will be used initially to prep any workspaces that need to be migrated (`tfm copy workspaces`). Migrations are usually a planned projects that will occur over time. Often there will need to be a final cutover where `tfm` can be used to update any changes from the source OR not all workspaces can be migrated initially that are to happen in a later phase.

Check out [customer journey example](./migration/journey.md) using `tfm` and what a Professions Services engagement looks like. 



## Can `tfm` perform a TFE to TFE migration?

Yes, we developed `tfm` to utilise the `go-tfe` library which is used for both Terraform Enterprise as well as Terraform Cloud. The following is what is capable

- TFE to TFC (Primary use case)
- TFE to TFE
- TFC to TFC
- TFC to TFE
- TFC ORG 1 to TFC ORG 2


## What constraints are there with migration to TFC?

The following are environment/configuration constraints where a migration of workspaces cannot occur:

- TFE Instances utilising a [Custom (Alternative Terraform Build Worker image)](https://developer.hashicorp.com/terraform/enterprise/install/interactive/installer#custom-image) as TFC does not support this feature. 
- TFE environments utilising [Network Mirror Provider protocol](https://developer.hashicorp.com/terraform/internals/provider-network-mirror-protocol)
    - A strategy for this is to change the workspace configuration in TFC to utilise Cloud Agents which requires further strategy and planning.
- Workspaces [pre 0.12](https://developer.hashicorp.com/terraform/cloud-docs/agents/requirements#supported-terraform-versions) cannot use Cloud Agents in TFC.
    - They would need to be upgraded by workspace owners before migrating to TFC.
- Customers with ONLY private Version Control Systems (VCS), TFC doees have a [list of supported VCS](https://developer.hashicorp.com/terraform/cloud-docs/vcs) solutions, however if private, some features of TFC may not work as intented.
- Workspaces that utilise `local-exec` or `remote-exec` [provisioner](https://developer.hashicorp.com/terraform/enterprise/install/interactive/installer#custom-image). 


## Will this work on a very old version of Terraform Enterprise?

In all honesty, we have not tested in anger what versions of `go-tfe` will not work with `tfm`.  Internal HashiCorp engineers do have the ability to spin up an older version of TFE. Let us know if you need help, we have a test-pipeline our github actions/test directory that can help populate TFE. 


## Is `tfm` supported by our HashiCorp Global Support Team?

Currently there is no official support whatsoever for `tfm`. This project was developed purposely built intially to assist Implementation Engineers if a migration project was to occur as we knew a few key customers had been asking for it. 



