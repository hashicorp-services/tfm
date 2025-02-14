# Exmple Scenario for Migrating From Terraform Community Edition to TFC/TFE

## Happy Path Scenario

Customer has terraform configurations managed by terraform community edition and would like to start managing these configurations with Terraform Cloud or Terraform Enterprise Workspaces.

### VCS

They are using a tfm [supported](../migration/supported-vcs.md) Version Control System (VCS) to manage their terraform configurations. That VCS has been connceted to TFC/TFE and the code owners are ready to start creating workspaces to link to the VCS repositories.

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

- The repository contains the terraform configurations and the backend{} configuration block.
- The repository can be a monorepo with many configurations in many directory paths.
- The repository can be using 1 backend to share multiple state files using terraform CE workspaces.

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
vcs_type = "github"
github_token = "token with read permissions to cloned repos"
github_organization = "organization"
github_username = "username"
clone_repos_path = "/opt/tfm-migration/repos"

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

tfm populates the `clone_repos_path` with the 4 repos defined in the config file above.

### Build the Metadata File

With the repos cloned locally the migration team builds a metadata file with information about how the repositories are using terraform.

`tfm core init-repos`

tfm looks through all paths in the repo and identifies paths using a `.tf` file with a `terraform { backend {} }` configuration. tfm also runs `tfm workspace list` and determines if the path is using terraform ce workspaces. tfm builds a `terraform_config_metadata.json` file in the tfm working directory that contains information about each repo.

```json
  {
    "repo_name": "isengard",
    "config_paths": [
      {
        "path": "isengard/infra/east/primary",
        "workspace_info": {
          "uses_workspaces": true,
          "workspace_names": [
            "default",
            "newisengard",
            "oldisengard"
          ]
        }
      },
      {
        "path": "isengard/infra/east/secondary",
        "workspace_info": {
          "uses_workspaces": false,
          "workspace_names": [
            "default"
          ]
        }
      },
```

### Get State

With the repos cloned locally and the metadata file built the migration team retrieves the state files for each repo.

`tfm core getstate`

tfm will use the `terraform_config_metadata.json` config file to  iterate through all of the cloned repositories in the `clone_repos_path` and metadata `config_paths` to download the state files from the backend.

tfm will use the locally installed terraform binary to perform `terraform init` and `terraform state pull > .terraform/pulled_terraform.tfstate` commands.

If tfm cannot successfully run a `terraform init` for a cloned repo tfm will return an error and continue with the next repository initilization attempt.

For any `config_path` with `uses_workspaces: true`, tfm will run `tfm workspace select` for each workspace in the `workspace_names` list and `terraform state pull > .terraform/pulled_<worspace name>_terraform.tfstate`. The end result will be multiple state files within the `config_path` for each workspace.

### Create Workspaces

The migration team creates TFC/TFE workspaces for each terraform configuration and state file to be managed by TFE/TFC

`tfm core create-workspaces`

tfm will use the `terraform_config_metadata.json` config file to create a TFC/TFE workspace in the `dst_tfc_org` defined in the config file for each repository.

Workspace names are generated using the metadata file in the following format:

- Config path using terraform ce workspaces:
`repo_name+config_path+workspace_name`

- Config path without terraform ce workspaces:
`repo_name+config_path`

As an example, the below metadata would create 2 TFC/TFE workspaces with the names:

- `isengard-infra-east-primary-newisengard`
- `isengard-infra-east-primary-oldisengard`

```json
  {
    "repo_name": "isengard",
    "config_paths": [
      {
        "path": "isengard/infra/east/primary",
        "workspace_info": {
          "uses_workspaces": true,
          "workspace_names": [
            "default",
            "newisengard",
            "oldisengard"
          ]
        }
      },
```

### Migrate State

The migration team can now upload state to the workspaces.

`tfm core upload-state`

`tfm core upload-state` is used to upload the state files that were downloaded using the `tfm core getstate` command to workspace created with the `tfm core create-worksapces` command. tfm will use the `terraform_config_metadata.json` config file to iterate through all of the `config_paths`. Any config path containing a `.terraform/pulled_terraform.tfstate` or `.terraform/pulled_workspaceName_terraform.tfstate` file will have the state file uploaded to a workspace that matches with the `config_path`.

### Link Repositories

The migration team can now attach the repository containing the terraform code to the workspace.

`tfm core link-vcs`

tfm will use the `terraform_config_metadata.json` file and look at each `config_path` and find a matching TFC/TFE workspace with a name that matches to the path.

tfm will update the workspace settings using the `vcs_provider_id` defined in the config file and the `repo_name` that the `config_path` belongs to and map the repo to the workspace.

### Validation

The migration team assigns a TFC/TFE variable set to each workspace with the correct credentials required to authenticate with the provider.

The migration team can now run a terraform plan and apply on the workspace and expect no changes to be made. If changes are being shown then verify there were no changes expected before migration.

### Post `tfm` Migration Tasks

#### Cleanup

The migration team has verified that everything is working and can begin the cleanup process.

The team wants to remove the stale backend{} configurations from all of the repos.

`tfm core remove-backend` iterates through all of the cloned repos and looks at all files ending in a `.tf` extension for a `backend{}` configuration. tfm creates a branch, comments out the backend, commits the change, and pushes the branch.

tfm DOES NOT create a PR. That is the responsibility of the code owners.

`tfm core cleanup` will remove all of the cloned repos from the clone path defined in the configuration file.

### Example GitHub Actions Pipeline

Coming soon
