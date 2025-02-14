# Pre-Requisites For Migrating From Terraform Community Edition to TFC/TFE

Note: The Terraform Community Edition migration as part of `tfm` has been deprecated in favor of [tf-migrate](https://developer.hashicorp.com/terraform/cloud-docs/migrate/tf-migrate). The CE migration feature has not been removed from `tfm` however it will not be receiving further developments.

## The following pre-reqs should be completed in the destination TFC/TFE before using tfm

- [Supported](../migration/supported-vcs.md) VCS provisioned in TFC/TFE

## The Following pre-reqs are Required to use the tfm Features for Cloning VCS Repositories

- A VCS token with permissions to read each repository of interestin the GitHub Organization.
- A Github organization or GitLab Project depending on the [supported vcs](../migration/supported-vcs.md) in use.
- A VCS username.

## The Following pre-reqs are Required to use the tfm Features for Removing Backend Configurations From Cloned Repositories

- A VCS token permissions to write contents to repositories.

## The Following pre-reqs are Required to use the tfm Features for Retrieving State Files

- The execution environment must provide credentials to the backend
- Terraform CLI must be installed in the execution environment
- The `tfm core init-repos` command must be run to create a metadata file

## Constraints

The following are environment/configuration constraints where a migration using tfm cannot occur:

- At the time of this writing tfm only supports the cloning of GitHub and GitLab repositories.
- At this time there is no way to handle CLI driven workspace migrations.
- At this time there is no way to handle variable migration.
