# HashiCorp Implementation Services

Migrating from Terraform Community Edition to Terraform Cloud or Terraform Enterprise.

Note: The Terraform Community Edition migration as part of `tfm` has been deprecated in favor of [tf-migrate](https://developer.hashicorp.com/terraform/cloud-docs/migrate/tf-migrate). The CE migration feature has not been removed from `tfm` however it will not be receiving further developments.

## Pre-Migration Questionnaire

### Software

- Are you permitted to run TFM in an environment that can access the supported VCS and TFC/TFE at the same time?
- Terraform CLI must be installed in the execution environment to use the tfm features.
- The local execution environment running tfm must be able to authenticate to the backend store state files to use the tfm features for downloading state.

### Terraform Cloud or Terraform Enterprise

- TFE Version?  (If applicable)
- Number of TFE/TFC Organizations that you wish to split terraform configurations amongst if more than 1.
- Which VCS? Only [supported vcs types](../migration/supported-vcs.md) at this time.
- Is the VCS routable from the internet and serving a publicly trusted certificate?
- Once configurations are migrated, will TFE/TFC have network connectivity to manage the inrastructure and download the required providers and modules defined in your code?

### Terraform Community Edition Configurations

- How many configurations do you believe you have?
- Where are state files being stored today? (s3, azure, etc.)
- How are VCS repositories that contain terraform community edition configurations structured?
- Are there repositories that contain multiple configurations split between miltiple directories?
- Are you using terraform workspaces and managing multiple state files with 1 backend configuration in the same terraform configuration?
- Are you using third party tools like terragrunt to manage terraform configurations?

### Project Flow

The migration propject has 6 key phases:

1. Discovery
1. Planning
1. Configuration
1. Technical Validation
1. Migration
1. End-User Validation

#### Discovery

- Determine current Terraform Community Edition landscape
  - Version Control Providers (Which VCS)
  - Structure(s) of VCS repositories (monorepos, terraform workspaces in use, terragrunt or other 3rd party tools being used)
  - Estimated number of repositories and terraform configurations per repository.
  - Estimated number of total terraform configurations to migrate.
- Establish Workspace criteria required to be eligible for migration. Migrating everything or just some.
- Discuss workspace creation
  - tfm only has the capabilities to create workspaces with constructed from the metadata file tfm generates with the following format: `repo_name+config_path+workspace_name` for configurations using terraform ce workspaces and `repo_name+config_path` for configurations not using terraform ce workspaces.
  - Future releases aim to allow workspace names to be modified during creation or map to existing worksapces created with the terraform tfe provider.
- Discuss State Migration on Workspaces (Latest State only)

#### Planning

- Gather the requied credentials for migrating with tfm
  - VCS token with read permissions for each repo of interest
  - (Optional) VCS token with write permissions for each repo of interest. This is only needed if you intend to use the remove-backend command to assist in commenting out backend configurations post migraion.
  - Terraform Cloud or Enterprise token with permissiosn to create workspaces.
- Determine if TFE/TFC workspaces will have the ability to manage the resources after migration. Will TFE/TFC have network access and authentication capabilities to manage the migrated resources?
- Establish nigration Timeline and Priority
  - Sometimes multiple configuration files are created to migrate sets of configurations to TFE/TFC at intervals.
- Establish Validation process as agreed upon with the Customer.
  - Running a plan and apply after the migration to verify no changes are expected.
  - How to assign credentials to all of the workspaces after migration.
- Establish a control plane to run tfm from to perform migration tasks. Requires access to the Terraform Enterprise/Terraform Cloud destination org and the VCS at the same time. Requires credentials to authenticate to the terraform backends.

#### Configuration

- Configure the TFE/TFC Version Control connection
- Configure the tfm configuration file
- Install terraform in the working environment
- Configure backend authentication credentials in the working environment

#### Technical Validation (Proof of Concept)

- Leveraging a subset of VCS repos, migrate enough to verify the migration process.
- Test validation steps on the items migrated configurations by running a plan and apply in the worksapce.

#### Migrate

- Perform migration of terraform configurations in the priority order identified in planning.
  - ( Optional) Run tfm core migrate to migrate an entire GitHub org or Gitlab Project at once.
  - ( Recommended ) Run tfm core commands in the order defined:
    - tfm core clone
    - tfm core init-repos
    - tfm core getstate
    - tfm core create-workspaces
    - tfm core upload-state
    - tfm core link-vcs
  - Run plan and apply on TFC/TFE workspaces and verify no changes are expected.
  - Validate.

  - tfm core remove-backend if desired to cleanup the code and push branches.

#### End-User Validation

- Run plan and apply on TFC/TFE workspaces and verify no changes are expected.
  - Validate.

### Migration Flow

The high level migration flow has 6 key steps:

1. Clone terraform configuration repositories
2. Build the metadata file describing the layout of each repo
3. Retreieve state files
4. Create TFC/TFE workspaces
5. Upload state files to workspaces
6. Link configuration repositories to workspaces
7. Cleanup configuration code
