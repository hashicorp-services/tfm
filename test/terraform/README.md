# Pre-reqs for tfm testing

For TFM to complete unit testing, Terraform infrastructure needs to be created in source and destination TFC orgs and source workspaces with runs need to be created.

This Terraform code will build out the following resources

- Workspaces
- Agent Pools
- SSH Keys
- Teams
- VCS Integration with GitHub
- Variable Sets
- Workspace Variables

Once completed with unit testing, a terraform destroy will run removing all these resources, leaving the orgs ready to be used again.

A workspace named `unit-test-baseline` in the `hc-implementation-services` organization maintains the state file for these Terraform resources.

#### Workspace Variables
```
destination_tfe_token = tfm-testing-destination token
source_tfe_token = tfm-testing-source token
tfe_token = hc-implementation-services token 
gh_token = the GitHub OAuth token for the
```

The `tfm-testing-source` and `tfm-testing-destination` organizations are upgraded TFC organizations used for unit testing. If these organizations are deleted, recreate them and request them to be upgraded in the `#team-se-trial-requests` HashiCorp slack channel.

>Note: Org tokens cannot be used to create SSH keys.

### Generating State Files
After resources are created run `/tests/state/create_state.sh` to generate state files within the workspaces.