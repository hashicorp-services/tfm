# tfm

![TFM](site/docs/images/TFM-black.png)

HashiCorp Implementation Services (IS) has identified a need to develop a purpose built tool to assist our engagements and customers during:

- Terraform open source / community edition / core to Terraform Cloud (TFC) & Terraform Enterprise (TFE)
- TFE to TFC
- TFC to TFE
- 1 TFC Organization to another TFC Organization
- TFE Server / Organization Consolidation

> [!Warning]
> This CLI does not have official support, but the code owners will work with partners and interested parties to provide assitance when possible.
> Check out our [case studies](https://hashicorp-services.github.io/tfm/migration/case-studies/).

## Overview

This tool has been develop to assist HashiCorp Implementation Services, Partners and Customers with their migrations to HashiCorp services. Having a tool allows us the ability to offer a standardized offering to our customers.

Check out the full documentation at [https://hashicorp-services.github.io/tfm/](https://hashicorp-services.github.io/tfm/)

Note: The Terraform Community Edition migration as part of `tfm` has been deprecated in favor of [tf-migrate](https://developer.hashicorp.com/terraform/cloud-docs/migrate/tf-migrate). The CE migration feature has not been removed from `tfm` however it will not be receiving further developments.

## Installation

Binaries are created as part of a release, check out the [Release Page](https://github.com/hashicorp-services/tfm/releases) for the latest version.

## Configuration

`tfm` supports two equivalent ways to supply configuration — use whichever fits your workflow:

### Option 1 — `.tfm.hcl` config file (recommended for interactive use)

By default `tfm` looks for `.tfm.hcl` in the current directory (or `~/.tfm.hcl`). Pass a custom path with `--config`.

```hcl
# .tfm.hcl
src_tfe_hostname = "app.terraform.io"
src_tfe_org      = "my-source-org"
src_tfe_token    = "..."

dst_tfc_hostname = "app.terraform.io"
dst_tfc_org      = "my-destination-org"
dst_tfc_token    = "..."
```

See the full documentation at [https://hashicorp-services.github.io/tfm/](https://hashicorp-services.github.io/tfm/) for all available keys.

### Option 2 — Environment variables (recommended for CI/CD pipelines)

Config keys that map cleanly to uppercase environment variable names can be overridden with no prefix required. For example, `src_tfe_token` is read from `SRC_TFE_TOKEN`.

> **Note:** Hyphenated config keys such as `projects-map`, `vcs-map`, and `exclude-ws-remote-state-resources` are not automatically translated to underscore-style environment variable names by the current implementation. Set those values in `.tfm.hcl` instead.

A template for commonly used environment-variable-based settings is provided at [`.env.example`](./.env.example). Copy it to `.env` and populate your values:
```bash
cp .env.example .env
# edit .env with your tokens and org names
export $(grep -v '^#' .env | xargs)
tfm list workspaces
```

> **Note:** `.env` is listed in `.gitignore` — never commit a populated `.env` file.

Key variables at a glance:

| Variable | Description |
|---|---|
| `SRC_TFE_HOSTNAME` | Source TFE/HCP Terraform hostname (e.g. `app.terraform.io`) |
| `SRC_TFE_ORG` | Source organisation name |
| `SRC_TFE_TOKEN` | Source API token |
| `DST_TFC_HOSTNAME` | Destination TFE/HCP Terraform hostname |
| `DST_TFC_ORG` | Destination organisation name |
| `DST_TFC_TOKEN` | Destination API token |
| `GITHUB_TOKEN` | GitHub PAT (required for `core` VCS migration commands) |
| `GITLAB_TOKEN` | GitLab PAT (required for `core` VCS migration commands) |

See [`.env.example`](./.env.example) for the complete list including VCS, GitLab, and core migration variables.

## Migration Type

There are differences between Terraform OSS / CE / Core migrations and Terraform Enterprise & Terraform Cloud Migrations. Accordingly this Readme has been split between two pages [tfe migration](./tfe-migration.md) and [ce migration](./ce-migration.md)

## tfm in a Pipeline

`tfm` can be used in a pipeline to automate migrations. There are a few considerations when using `tfm` in this manner.

- For source and destination credentials, use environment variables sourced from your pipeline's secret manager or another secure mechanism. See [`.env.example`](./.env.example) for the full list of supported variables. Don't embed credentials in `tfm` configuration files.
- Several `tfm` commands require confirmation before proceeding, which are listed below. To override these in a pipeline, add the `--autoapprove` flag.
  - `copy workspaces` Only when all workspaces are going to be migrated due to no workspace list, or map defined.
  - `copy workspaces --state --last`
  - `delete workspace`
  - `delete workspaces-vcs`

## Architectural Decisions Record (ADR)

An architecture decision record (ADR) is a document that captures an important architecture decision made along with its context and consequences.

This project will store ADRs in [docs/ADR](docs/ADR/) as a historical record.

More information about [ADRs](docs/ADR/index.md).

## To build

```bash
make build-local
./tfm -v
```

-or-

```bash
go run . -v
```

## To release

To create a new release of TFM

- Locally Create a new Tag
- Push tag to Origin
- GitHub workflow `release.yml` will create a build and make the release in GitHub

## Reporting Issues

If you believe you have found a defect in `tfm` or its documentation, use the [GitHub issue tracker](https://github.com/hashicorp-services/tfm/issues) to report the problem to the `tfm` maintainers.
