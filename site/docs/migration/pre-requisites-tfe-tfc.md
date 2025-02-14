# Pre-Requisites

## The following pre-reqs should be completed in the destination TFC/TFE before using tfm

- A TFC/TFE Token with owner permissions is required
- Existing Workspaces should have a recent clean TF Plan/Apply
- VCS provisioned
  - VCS Map provided as configuration file
- Teams created
  - Team map provided as configuration file
- Agent Pools created
  - Agent map provided as configuration file
- Variable Sets created
  - Variable Set map provided as configuration file
- Variables with secrets known OR can be regenerated

## Constraints

The following are environment/configuration constraints where a migration of workspaces cannot occur:

- TFE Instances utilising a [Custom (Alternative Terraform Build Worker image)](https://developer.hashicorp.com/terraform/enterprise/install/interactive/installer#custom-image) as TFC does not support this feature.
- TFE environments utilising [Network Mirror Provider protocol](https://developer.hashicorp.com/terraform/internals/provider-network-mirror-protocol)
  - A strategy for this is to change the workspace configuration in TFC to utilize Cloud Agents which requires further strategy and planning.
- Workspaces [pre 0.12](https://developer.hashicorp.com/terraform/cloud-docs/agents/requirements#supported-terraform-versions) cannot use Cloud Agents in TFC.
  - They would need to be upgraded by workspace owners before migrating to TFC.
- Customers with ONLY private Version Control Systems (VCS), TFC doees have a [list of supported VCS](https://developer.hashicorp.com/terraform/cloud-docs/vcs) solutions, however if private, some features of TFC may not work as intented.
- Workspaces that utilize `local-exec` or `remote-exec` [provisioner](https://developer.hashicorp.com/terraform/enterprise/install/interactive/installer#custom-image).
