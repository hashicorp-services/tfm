# TFM FAQs

## Who is `tfm` developed for?

Engineers/Operators that manage/admin Terraform Enterprise/Cloud organizations that need to perform a migration of workspaces. 

## What is `tfm` intended to do?

`tfm` will assist with migration of TFE/TFC workspaces from one TFE/TFC instance/organization to another TFE/TFC instance/organization. 


### Migrations of workspaces
From a customer journey perspective, it will be used initially to copy any workspaces that need to be migrated (`tfm copy workspaces`). Each aspect of a workspace can be copied/migrated to the destination by specifying flags. This allows flexibility and control. Check out the [pre-requisites](./migration/pre-requisites.md) when migrating a workspace. 

### Idempotent Migrations
We have designed `tfm` to be run more than once on an existing source workspace/s. This will allow users to not only update any changes from the source workspace, but keep tabs on workspaces that have yet to be migrated to the destination. 

### CLI Usage
The main usage will be used by TFE admins carefully migrating select workspaces for migration to TFE/TFC as shown in our [example migration scenario](../docs/migration/example-scenarios.md).   

We also envision some users to use the CLI tool in a pipeline and add workspaces over time for migration by using the `workspace-map` or `workspace-list` options in the config file. Check out an [example github actions pipeline](../docs/migration/example-scenarios.md#example-github-actions-pipeline). This allows teams to plan gradually which workspaces are ready for migration. 


### Future Migrations
Migrations are usually a planned projects that will occur over time. Often there will need to be a final cutover where `tfm` can be used to update any changes from the source OR not all workspaces can be easily migrated initially that require more technical preparations before migration.

### Customer Journey

Check out [customer journey example](./migration/journey.md) using `tfm` and what a Professional Services engagement looks like. 



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


## Does tfm store the state file to disk?

No tfm streams the statefile from the source to memory and then sends it to the destination. No statefile or any kind of file is written to disk. The statefile is also encrypted when fetched from the TFE API. 


## Will this work on a very old version of Terraform Enterprise?

In all honesty, we have not tested in anger what versions of `go-tfe` will not work with `tfm`.  Internal HashiCorp engineers do have the ability to spin up an older version of TFE test. Let us know if you need help, we have a test-pipeline in the project's github actions/test directory that can help populate TFE. 


## Is `tfm` supported by our HashiCorp Global Support Team?

Currently there is *no official support* whatsoever for `tfm`. This project was developed intially to assist TFE/TFC Administrators/Engineers if a migration project was to occur.


## I have a feature request

Check out [Future Work/Roadmap](./code/future.md) or our [Github Projects page](https://github.com/orgs/hashicorp-services/projects/6).

Got an idea for a feature to `tfm`? Submit a [feature request](https://github.com/hashicorp-services/tfm/issues/new?assignees=&labels=&template=feature_request.md&title=)! 
