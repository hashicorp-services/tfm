# TFM Configuration File Settings

| Parameter | Supported Values| Description | Required |
| --------- | --------------- | ----------- | -------- |
| src_tfe_hostname | A hostname such as app.terraform.io | The hostname of a TFE server that you are migrating from | `yes` for TFE to TFC or TFC to TFC migrations | 
| src_tfe_org | A TFC/TFE organization name | The TFE/TFC Organization that you are migrating from | `yes` for TFE to TFC or TFC to TFC migrations | 
| src_tfe_token | A TFC/TFE Token | A Token for the TFE/TFC Organization that you are migrating from | `yes` for TFE to TFC or TFC to TFC migrations | 
| dst_tfc_hostname | A hostname such as app.terraform.io | The hostname of a TFE server or the TFC hostname that you are migrating to | `yes` for all migrations | 
| dst_tfc_org | A TFC/TFE organization name | A TFC/TFE organization that you are migrating to | `yes` for all migrations | 
| dst_tfc_token | A TFC/TFE Token | | `yes` for all migrations | 
| repos_to_clone | A list of VCS repository names | Used with the`tfm core clone` command to clone a set of VCS repositories. If not provided, all VCS repos will be cloned | `no` | 
| vcs-map | A list of source=destination VCS oauth IDs | TFM will look at each workspace in the source for the source VCS oauth ID and assign the matching workspace in the destination with the destination VCS oauth ID | `yes` for `tfm copy workspaces --vcs` |
| workspaces | A list of workspaces to migrate from TFE to TFC or TFC org to TFC org | Provide a list of source workspaces in the source TFC/TFE org to migrate. If not provided and no "workspaces-map" is detected, all workspaces will be migrated. | `no` |
| exclude-workspaces | A list of workspaces to exclude when copying. | Conflicts with workspaces and workspaces-map. | `no` |
| projects | A list of projects to migrate across from TFE to TFC or TFC org to TFC org | Provide a list of source projects in the source TFC/TFE org to migrate. If not "projects-map" if detected, all projects will be migrated | `no` | 
| projects-map | A list of source=destination project names | TFM will look at each project in the source for the source project name and recreate the project in the destination with the new destination project name. Takes precedence over "projects" list. | `no` | 
| workspaces-map | A list of source=destination workspace names | TFM will look at each source workspace and recreate the workspace with the specified destination name | `no` | 
| commit_message | A commit message | Used when creating a branch for the `tfm core remove-backend` command | `yes` only for the `tfm core remove-backend` command | 
| commit_author_name | Author name to appear on commits | Used when creating a branch for the `tfm core remove-backend` command | `yes` only for the `tfm core remove-backend` command | 
| commit_author_email | Author email to appear on commits | Used when creating a branch for the `tfm core remove-backend` command | `yes` only for the `tfm core remove-backend` command | 
| github_token | A github token | Used for `tfm core` commands when `vcs_type = "github"` | `yes` only for `tfm core` migrations | 
| github_organization | A github organization | Used for `tfm core` commands when `vcs_type = "github" | `yes` only for `tfm core` migrations |
| github_username | A github username | Used for `tfm core` commands when `vcs_type = "github" | `yes` only for `tfm core` migrations |
| gitlab_token | A gitlab username | Used for `tfm core` commands when `vcs_type = "gitlab" | `yes` only for `tfm core` migrations |
| gitlab_username | A gitlab username | Used for `tfm core` commands when `vcs_type = "gitlab"` | `yes` only for `tfm core` migrations |
| gitlab_group | A gitlab group | Used for `tfm core` commands when `vcs_type = "gitlab"` | `yes` only for `tfm core` migrations |
| clone_repos_path | | | `yes` only for `tfm core` migrations |
| vcs_type | | | `yes` only for `tfm core` migrations |
| vcs_provider_id | | | `yes` only for `tfm core link-vcs` command migrations |
| agents_map | A list of source=destination agent pool IDs | TFM will look at each workspace in the source for the source agent pool ID and assign the matching workspace in the destination the destination agent pool ID. Conflicts with agent-assignment | `no` |
| agent-assignment-id | An agent Pool ID | An agent pool ID to assign to all workspaces in the destination. Conflicts with agents-map | `no` |
| varsets_map | A list of source=destination variable set names | TFM will look at each source variable set and recreate the variable set with the specified destination name | `no` |
| ssh-map | A list of source=destination SSH IDs | TFM will look at each workspace in the source for the source SSH  ID and assign the matching workspace in the destination with the destination SSH ID | `no` |
| | | | |


