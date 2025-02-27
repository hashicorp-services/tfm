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

## Migration Type

There are differences between Terraform OSS / CE / Core migrations and Terraform Enterprise & Terraform Cloud Migrations. Accordingly this Readme has been split between two pages [tfe migration](./tfe-migration.md) and [ce migration](./ce-migration.md)

## tfm in a Pipeline

`tfm` can be used in a pipeline to automate migrations. There are a few considerations when using `tfm` in this manner.

- For source and destination credentials use environment variables from the pipeline's secrete manager or other secure means. Don't embed these credentials in `tfm` configuration files
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
