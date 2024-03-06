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
| projects | A list of projects to migrate across from TFE to TFC or TFC org to TFC org | Provide a list of source projects in the source TFC/TFE org to migrate. If not "projects-map" if detected, all projects will be migrated | `no` | 
| | | | | 
| | | | | 
| | | | | 
| | | | | 
| | | | | 
| | | | | 
