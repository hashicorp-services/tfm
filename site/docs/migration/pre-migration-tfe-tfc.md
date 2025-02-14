# HashiCorp Implementation Services

Migrating from Terraform Enterprise to Terraform Cloud.

Updated copy of questionaire found at [Google Drive](https://docs.google.com/spreadsheets/d/1yi2TRF0G3AN7XTJQxdMneJHX2vTV6-BO4YNQs0F65Bc/edit?usp=sharing).

## Pre-Migration Questionnaire

### Terraform Enterprise, TFE, Current Deployment (Source)

- TFE Version?
- Number of TFE Organizations?
- Number of users/teams using TFE (10/100/1000)?
- Number of workspaces (per organization)?
- Number of Modules in the PMR (per organization)?
- Number of Sentinel Policies/Sets (per organization)?
- Number of TFE Teams (per organization)?
- Which VCS?
- Is the VCS routable from the internet and serving a publicly trusted certificate?
- Which IdP?
- Where are you deploying infrastructure from TFE today?
- Do you have the ability to run Python scripts or a binary tool from a developer workstation?
- Is there an established maintenance window for transitioning from TFE to TFC?

### Terraform Cloud, TFC, New Deployment (Destination)

- How many TFC Organizations?
- Do you want to migrate everything to TFC?
- Do you require Cloud Agents?
- Will the workspaces all be going into the same Org?

### Project Flow

The high level Migration path has 6 key components:

1. Discovery
1. Planning
1. Configuration
1. Technical Validation
1. Migration
1. End-User Validation

#### Discovery

- Determine current Terraform Enterprise landscape
  - Version Control Providers (Which VCS, Count, Distribution of Use)
  - Identity Platform (Which IdP, Number of Teams)
  - Modules in the Private Module Registry (Count, Publishing Method, No. of Versions, Frequency of Change)
  - Policies and Policy Sets (Count, Publishing Method, Frequency of Change)
  - Workspaces (Count, Publishing Method, Frequency of Change)
- Inventory Terraform Enterprise footprint (Aids in tracking completed work)
- Establish Workspace criteria required to be eligible for migration
- Determine if Cloud Agents are required (Typically due to on-premise deployments from Terraform)
  - Identify how Cloud Agents are deployed and configured
  - Determine if any utilities are needed within the agent (Local-Exec, Custom Providers, etc...)
- Discuss new user onboarding in Terraform Cloud (Slight differences from Terraform Enterprise)
- Discuss State Migration on Workspaces (Latest State only)
- Evaluate current use of Workspace Variables marked as "sensitive" (these values are write-only from the API)

#### Planning

- Identify Admins of VCS and IdP that will be needed to configure both in Terraform Cloud
- Determine order in which to migrate Modules, Policies, and most importantly Workspaces
- Establish Workspace Migration Timeline and Priority
- Establish Workspace Deprecation process in Terraform Enterprise (Common Options: Locking, Deleting, Archiving)
- Establish Validation process as agreed upon with the Customer
- Determine required API tokens needed for the migration (must have access to all needed Organizations)
  - Terraform Enterprise token with Owner permissions(Source)
  - Terraform Cloud token with Owner permissions (Destination)
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
