# HashiCorp Implementation Services
Migrating from Terraform Community Edition to Terraform Cloud or Terraform Enterprise.

## Pre-Migration Questionnaire

### Software

- Are you permitted to run TFM in an environment that can access the supported VCS and TFC/TFE at the same time?
- Terraform CLI must be installed in the execution environment to use the tfm features for downloading state.
- The local execution environment running tfm must be able to authenticate to the backend store state files to use the tfm features for downloading stae.

### Terraform Cloud or Terraform Enterprise
- TFE Version?  (If applicable)
- Number of TFE/TFC Organizations that you wish to split terraform configurations amongst if more than 1.
- Which VCS? Only GitHub is supported at this time.
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
    - Estimated number of repositories and terraform configurations.
- Establish Workspace criteria required to be eligible for migration. Migrating everything or just some.
- Discuss workspace creation
  - tfm only has the capabilities to create workspaces with the same name as the repo at this time.
  - Future releases aim to allow workspace names to be modified during creation or map to existing worksapces created with the terraform tfe provider.
- Discuss State Migration on Workspaces (Latest State only)

#### Planning
- Gather the requied credentials for migrating with tfm
    - GitHub token with read permissions for each repo of interest
    - (Optional) GitHub token with write permissions for each repo of interest. This is only needed if you intend to use the remove-backend command to assist in removing backend configurations post migraion.
    - Terraform Cloud or Enterprise token with permissiosn to create workspaces.
- Determine if TFE/TFC workspaces will have the ability to manage the resources after migration. Will TFE/TFC have network access and authentication capabilities to manage the migrated resources?
- Establish nigration Timeline and Priority
    - Sometimes multiple configuration files are created to migrate sets of configurations to TFE/TFC at intervals.
- Establish Validation process as agreed upon with the Customer.
    - Running a plan and apply after the migration to verify no changes are expected.
    - How to assign credentials to all of the workspaces after migration.
- Establish a control plane to run tfm from to perform migration tasks. Requires access to the Terraform Enterprise/Terraform Cloud destination org and the VCS GitHub org at the same time.

#### Configuration

- Configure the TFE/TFC Version Control connection
- Configure the tfm configuration file

#### Technical Validation (Proof of Concept)

- Leveraging a subset of VCS repos, migrate enough to verify the migration process.
- Test validation steps on the items migrated configurations by running a plan and apply in the worksapce.

#### Migrate
- Perform migration of terraform configurations in the priority order identified in planning.
    - ( Optional) Run tfm core migrate to migrate an entire GitHub org at once.
    - ( Recommended ) Run tfm core commands in the order defined:
        - tfm core clone
        - tfm core getstate
        - tfm core create-workspaces
        - tfm core upload-state
        - tfm core link-vcs
    - Run plan and apply on TFC/TFE workspaces and verify no changes are expected.
    - Validate.

#### End-User Validation
- Run plan and apply on TFC/TFE workspaces and verify no changes are expected.
    - Validate. 

### Migration Flow
The high level migration flow has 6 key steps:
1. Clone terraform configuration repositories
2. Retreieve state files
3. Create TFC/TFE workspaces
4. Upload state files to workspaces
5. Link configuration repositories to workspaces
6. Cleanup configuration code

