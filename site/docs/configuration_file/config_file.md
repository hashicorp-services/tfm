# TFM Configuration

`tfm` uses [Viper](https://github.com/spf13/viper) for configuration. Every parameter in the table below can be supplied via **either** a config file **or** an environment variable — they are fully interchangeable and can be mixed freely. Environment variables always take precedence over config file values.

## Option 1 — `.tfm.hcl` config file

By default `tfm` looks for `.tfm.hcl` in the current directory, then `~/.tfm.hcl`. Pass a custom path with `--config /path/to/file.hcl`.

```hcl
# .tfm.hcl — example for an HCP Terraform to HCP Terraform migration
src_tfe_hostname = "app.terraform.io"
src_tfe_org      = "my-source-org"
src_tfe_token    = "..."

dst_tfc_hostname = "app.terraform.io"
dst_tfc_org      = "my-destination-org"
dst_tfc_token    = "..."

workspaces = ["ws-alpha", "ws-beta"]
```

## Option 2 — Environment variables

Every config key maps to an uppercase environment variable of the same name (no prefix). For example, `src_tfe_token` is read from `SRC_TFE_TOKEN`.

This is the recommended approach for CI/CD pipelines — credentials are supplied by the pipeline's secret manager and never written to disk.

A fully documented template listing every supported variable is provided in [`.env.example`](https://github.com/hashicorp-services/tfm/blob/main/.env.example) at the root of the repository. Copy it to `.env` and populate your values:

```bash
cp .env.example .env
# edit .env — fill in your tokens, hostnames, and org names
export $(grep -v '^#' .env | xargs)
tfm list workspaces
```

> !!!! WARNING:
    Never commit a populated `.env` file. The repository's `.gitignore` already excludes `.env` to protect against accidental credential exposure.

### Quick-reference: key environment variables

| Environment Variable | Config Key | Description |
|---|---|---|
| `SRC_TFE_HOSTNAME` | `src_tfe_hostname` | Source TFE/HCP Terraform hostname (e.g. `app.terraform.io`) |
| `SRC_TFE_ORG` | `src_tfe_org` | Source organisation name |
| `SRC_TFE_TOKEN` | `src_tfe_token` | Source API token |
| `DST_TFC_HOSTNAME` | `dst_tfc_hostname` | Destination TFE/HCP Terraform hostname |
| `DST_TFC_ORG` | `dst_tfc_org` | Destination organisation name |
| `DST_TFC_TOKEN` | `dst_tfc_token` | Destination API token |
| `GITHUB_TOKEN` | `github_token` | GitHub PAT (required for `tfm core` with `vcs_type = "github"`) |
| `GITLAB_TOKEN` | `gitlab_token` | GitLab PAT (required for `tfm core` with `vcs_type = "gitlab"`) |

See [`.env.example`](https://github.com/hashicorp-services/tfm/blob/main/.env.example) for the complete list.

---

## All configuration parameters

| Parameter | Supported Values| Description | Required |
| --------- | --------------- | ----------- | -------- |
| src_tfe_hostname | A hostname such as app.terraform.io | The hostname of a TFE server that you are migrating from | `yes` for TFE to TFC or TFC to TFC migrations |
| src_tfe_org | A TFC/TFE organization name | The TFE/TFC Organization that you are migrating from | `yes` for TFE to TFC or TFC to TFC migrations |
| src_tfe_token | A TFC/TFE Token | A Token for the TFE/TFC Organization that you are migrating from | `yes` for TFE to TFC or TFC to TFC migrations |
| dst_tfc_hostname | A hostname such as app.terraform.io | The hostname of a TFE server or the TFC hostname that you are migrating to | `yes` for all migrations |
| dst_tfc_org | A TFC/TFE organization name | A TFC/TFE organization that you are migrating to | `yes` for all migrations |
| dst_tfc_token | A TFC/TFE Token | A Token for the TFE/TFC Organization that you are migrating to | `yes` for all migrations |
| exclude-ws-remote-state-resources | true/false | Option to skip workspaces which use [remote state data](https://developer.hashicorp.com/terraform/language/state/remote-state-data). | `no`|
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
| github_token | A GitHub personal access token | Used for `tfm core` commands when `vcs_type = "github"` | `yes` only for `tfm core` migrations |
| github_organization | A GitHub organisation name | Used for `tfm core` commands when `vcs_type = "github"` | `yes` only for `tfm core` migrations |
| github_username | A GitHub username | Used for `tfm core` commands when `vcs_type = "github"` | `yes` only for `tfm core` migrations |
| gitlab_token | A GitLab personal access token | Used for `tfm core` commands when `vcs_type = "gitlab"` | `yes` only for `tfm core` migrations |
| gitlab_username | A GitLab username | Used for `tfm core` commands when `vcs_type = "gitlab"` | `yes` only for `tfm core` migrations |
| gitlab_group | A GitLab group | Used for `tfm core` commands when `vcs_type = "gitlab"` | `yes` only for `tfm core` migrations |
| clone_repos_path | A local filesystem path | Path where VCS repositories will be cloned during `tfm core clone` | `yes` only for `tfm core` migrations |
| vcs_type | `github` or `gitlab` | The VCS provider type for `tfm core` migrations | `yes` only for `tfm core` migrations |
| vcs_provider_id | | | `yes` only for `tfm core link-vcs` command migrations |
| agents_map | A list of source=destination agent pool IDs | TFM will look at each workspace in the source for the source agent pool ID and assign the matching workspace in the destination the destination agent pool ID. Conflicts with agent-assignment | `no` |
| agent-assignment-id | An agent Pool ID | An agent pool ID to assign to all workspaces in the destination. Conflicts with agents-map | `no` |
| varsets_map | A list of source=destination variable set names | TFM will look at each source variable set and recreate the variable set with the specified destination name | `no` |
| ssh-map | A list of source=destination SSH IDs | TFM will look at each workspace in the source for the source SSH  ID and assign the matching workspace in the destination with the destination SSH ID | `no` |
| | | | |


