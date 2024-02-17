# Pre-Requisites For Migrating From Terraform Community Edition to TFC/TFE

## The following pre-reqs should be completed in the destination TFC/TFE before using tfm:

- GitHub VCS provisioned in TFC/TFE

## The Following pre-reqs are Required to use the tfm Features for Cloning VCS Repositories

- A Github token with permissions to read each repository of interestin the GitHub Organization.
- A Github organization
- A Github username

## The Following pre-reqs are Required to use the tfm Features for Removing Backend Configurations From Cloned Repositories

- A GitHub token permissions to write contents to repositories.

## The Following pre-reqs are Required to use the tfm Features for Retrieving State Files

- The execution environment must provide credentials to the backend
- Terraform CLI must be installed in the execution environment

## Constraints

The following are environment/configuration constraints where a migration using tfm cannot occur:

- At the time of this writing tfm only supports the cloning of GitHub repositories.
- At the time of this writing tfm does not account for mono repo type of configurations in which multiple terraform configurations are nested in 1 repo, each with its own state file. tfm only looks at and acts upon the root of the repository.
- At the time of this writing tfm does not account for configurations using the terraform workspace feature using 1 backend configuration to store multiple state files.



 