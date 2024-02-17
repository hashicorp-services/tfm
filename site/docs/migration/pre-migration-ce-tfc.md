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
- Establish Workspace criteria required to be eligible for migration
    - tfm names workspaces it creates after the repo it clones.
- Discuss State Migration on Workspaces (Latest State only)
#### Planning
- Gather the requied credentials for migrating with tfm
    - GitHub token with read permissions for each repo of interest
    - (Optional) GitHub token with write permissions for each repo of interest. This is only needed if you intend to use the remove-backend command to assist in removing backend configurations post migraion.
    - Terraform Cloud or Enterprise token with permissiosn to create workspaces.
- Determine order in which to migrate Modules, Policies, and most importantly Workspaces
- Establish Workspace Migration Timeline and Priority
- Establish Workspace Deprecation process in Terraform Enterprise (Common Options: Locking, Deleting, Archiving)
- Establish Validation process as agreed upon with the Customer
- Determine required API tokens needed for the migration (must have access to all needed Organizations)
    - Terraform Enterprise (Source)
    - Terraform Cloud (Destination)
- Determine if Workspace Code changes are required (Module Sourcing, Provider Initialization)
- Design Cloud Agent Pool structure and required infrastructure, including networking and authentication routes
- Determine how to update destination Workspace Variables that are marked as "sensitive" in the source Workspace
- Establish a control plane to run python from to perform migration tasks (Requires access to both Terraform Enterprise and Terraform Cloud)
#### Configuration
- Configure Terraform Cloud for Version Control
- Configure Terraform Cloud for Single Sign-On
- Create Teams in Terraform Cloud
- (If Needed) Create Cloud Agent Pools, Deploy Agents, and register the Agents with their respective Agent Pool


#### Technical Validation (Proof of Concept)
- Leveraging a subset of Modules, Policies, and Workspaces, migrate enough to verify the migration process
- Test validation steps on the items migrated
#### Migrate
- Perform Migration of Modules in the Private Module Registry, and Validate
- Perform Migration of Sentinel Policies and Policy Sets, and Validate
- Perform Migration of Workspaces in priority order, and Validate
    - Run Plan on Terraform Enterprise
    - Perform any code changes to the Terraform consumed by the Workspace
    - Migrate Workspace to Terraform Cloud
    - Run Plan on Terraform Cloud (verify it matches the Plan from Terraform Enterprise)
    - Validate
    - Deprecation of Terraform Enterprise Workspace
#### End-User Validation
- Verify end users can access and leverage Workspaces in Terraform Cloud as they would have in Terraform Enterprise

### Migration Flow
The high level migration flow has 6 key steps:
1. Clone terraform configuration repositories
2. Retreieve state files
3. Create TFC/TFE workspaces
4. Upload state files to workspaces
5. Link configuration repositories to workspaces
6. Cleanup configuration code

