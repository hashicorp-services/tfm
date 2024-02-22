# Exmple Scenario for Migrating From Terraform Community Edition to TFC/TFE

## Happy Path Scenario

Customer has terraform configurations managed by terraform community edition and would like to start managing these configurations with Terraform Cloud or Terraform Enterprise Workspaces.

### VCS
They are using a tfm supported Version Control System (VCS) to manage their terraform configurations. Only GitHub is supported at the time of this release. That VCS has been connceted to TFC/TFE and the code owners are ready to start creating workspaces to link to the VCS repositories.

### State Files
The state files for the VCS stored terraform configurations are stored in a supported terraform backend such as S3 and a backend{} block is configured within the terraform configurations. Code owners want to have the state managed by a TFE/TFC workspace.

### Terraform VCS Repositories
The following is a list of terraform configuration repositories in GitHub that an example customer has identified for migration:

```bash
application-infra
application-networking
application-database
application-rbac
```

A suitable repository for migration with tfm (at this time) has the following requirements:

- The root of the repository contains the terraform configurations and the backend{} configuration block.
- The initial release of the CE migration feature only supports this configuration. There is no support (yet) for monorepos that contain multiple directories with multiple terraform configurations and backend{} configurations. There is no support for repositories that user the terraform workspace feature to manage multiple state files with one backedn{} configuration.


### Preparing the destination (TFC/TFE organization)

In preparation of TFC, the following are completed to prepare for migration:

- GitHub connected as a [VCS provider](https://developer.hashicorp.com/terraform/cloud-docs/vcs/github-app)
- A variable set is configured with provider credentials to run the validation plan and apply on a workspace after migration.

### Setting up the TFM config file

The customer configures the TFM file for migration.

The following is what a `~/.tfm.hcl` file will look like for `tfm core` commands to clone only the identified repositories for migration.

```hcl
dst_tfc_hostname="app.terraform.io"
dst_tfc_org="organization"
dst_tfc_token="token with permissions to create workspaces"
github_token = "token with read permissions to cloned repos"
github_organization = "organization"
github_username = "username"
github_clone_repos_path = "/opt/tfm-migration/repos"

"repos_to_clone" = [
  "application-infra",
  "application-networking",
  "application-database",
  "application-rbac"
]

# Only used only with the tfm core link-vcs command
vcs_provider_id = "ot-K9ofy9Rr9R9Bo3Nj" 

# Only used with the tfm core remove-backend command
commit_message = "commit message"
commit_author_name = "name"
commit_author_email = "email"
```

### Clone Repos

With the configuration file configured the migration team clones the repos to the local host or execution environment.

`tfm core clone`

tfm populates the `github_clone_repos_path` with the 4 repos defined in the config file above.

### Get State

With the repos cloned locally the migration team retrieves the state files for each repo.

`tfm core getstate`

tfm runs `terraform init` and `terraform state pull > .terraform/pulled_terraform.tfstate` to retrieve that state for each repo. Terraform must be installed locally or in the execution environment and available in the path.

### Create Workspaces
The migration team creates TFC/TFE workspaces for each terraform configuration and state file to be managed by TFE/TFC

`tfm core create-workspaces`

tfm looks at all of the cloned repos and looks for the ones containing `.terraform/pulled_terraform.tfstate` and creates a workpace with the identical name as the cloned repository directory.


### Migrate State
The migration team can now upload state to the workspaces.

`tfm core get-state`

tfm looks at all of the cloned repos and looks for the ones containing `.terraform/pulled_terraform.tfstate` and uploads the state file to the workspace with the same name as the repository directory it belongs to.

### Link Repositories
The migration team can now attach the repository containing the terraform code to the workspace.

`tfm core link-vcs`

tfm looks at all of the cloned repos and looks for a TFC/TFE workspace with the same name as the cloned repository directory. tfm uses the provided `vcs_provider_id` in the tfm configuration file and the name of the repository to update the workspaces VCS connection.
 
### Validation
The migration team can now run a terraform plan and apply on the workspace and expect no changes to be made. If changes are being shown then verify there were no changes expected before migration.


### Post `tfm` Migration Tasks


#### Cleanup

The migration team has verified that everything is working and can begin the cleanup process.

The team wants to remove the stale backend{} configurations from all of the repos.

`tfm core remove-backend` iterates through all of the cloned repos and looks at all files ending in a `.tf` extension for a `backend{}` configuration. tfm creates a branch, removes the backend, commits the change, and pushes the branch. 

tfm DOES NOT create a PR. That is the responsibility of the code owners.

`tfm core cleanup` will remove all of the cloned repos from the clone path defined in the configuration file.


### Example GitHub Actions Pipeline

Coming soon
