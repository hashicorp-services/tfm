# tfm

![TFM](site/docs/images/TFM-black.png)

HashiCorp Implementation Services (IS) has identified a need to develop a purpose built tool to assist our engagements and customers during migrations from Terraform open source / community edition / core to Terraform Cloud (TFC) or Terraform Enterprise (TFE) and migrations to and from Terraform Enterprise and Terraform Cloud.

> [!Warning]
> This CLI does not have official support, but the code owners will work with partners and interested parties to provide assitance when possible.
> Check out our [case studies](https://hashicorp-services.github.io/tfm/migration/case-studies/).

## Overview

This tool has been develop to assist HashiCorp Implementation Services, Partners and Customers with their migrations to and from HashiCorp services.

Check out the full documentation at [https://hashicorp-services.github.io/tfm/](https://hashicorp-services.github.io/tfm/)

## Installation

Binaries are created as part of a release, check out the [Release Page](https://github.com/hashicorp-services/tfm/releases) for the latest version.

## Migration Type

There are differences between Terraform OSS / CE / Core migrations and Terraform Enterprise & Terraform Cloud Migrations. Accodrinly the documentation has been split between two pages [tfe migration](./tfe-migration.md) and [ce migration](./ce-migration.md)

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

- Use GitHub Action `Release tfm`
- Specify a version number. Please follow semantic versioning for the release.

This action will do the following steps

- Compile TFM for Linux, Mac, Windows with amd64 and arm64 versions
- Upload the artifacts
- Create a new release + tag on the repo at the current main

## Reporting Issues

If you believe you have found a defect in `tfm` or its documentation, use the [GitHub issue tracker](https://github.com/hashicorp-services/tfm/issues) to report the problem to the `tfm` maintainers.
