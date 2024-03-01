# ADR #0000: Title, a short present tense phrase

Date: YYYY-MM-DD

## Responsible Architect
The Architect most closely aligned with this decision.

## Author

The person who wrote this document.

## Contributors

* List the names of the contributors here
* Name 1
* Name 2
* etc

## Lifecycle

POC (Proof of Concept), Pilot, Beta, GA (General Availability), Sunset

## Status

Status of the decision made. A decision may be "Proposed" if the project stakeholders haven't agreed with it yet, or "Accepted" once it is agreed. If a later ADR changes or reverses a decision, it may be marked as "Deprecated" or "Superseded" with a reference to its replacement. Once the decision has been implemented, it should be marked as "Implemented".

## Context

TFM only supports the GitHub VCS as of release 0.8.0. There are multiple VCS providers that are compatible with Terraform Cloud and Enterprise that customers use to store Terraform Community Edition code. TFM needs to support all of the same VCS providers for the `tfm core clone` and `tfm core remove-backend` commands along with any future commands that interact with a VCS provider. We should only aim to support VCS providers that TFC/TFE suport [documented here](https://developer.hashicorp.com/terraform/cloud-docs/vcs#supported-vcs-providers).

An ADR should be created for each VCS provider as they require different client auth methods and some have unique paths for storing VCS respoitories.

### Go-Gitlab

[A go-gitlab](https://pkg.go.dev/github.com/xanzy/go-gitlab) package exists to assist in developing this feature.

### Client Context

A new file `vcsclients/gitlab.go` should create the gitlab client context for use with tfm functions that interact with GitLab.

## Decision

This section describes our response to these forces and what we are actually proposing on doing. It is stated in full sentences, with active voice. "We will ..."

## Consequences

What becomes easier or more difficult because of the decision made. This section describes the resulting context, after applying the decision. All consequences should be listed here, not just the "positive" ones. A particular decision may have positive, negative, and neutral consequences, but all of them affect the team and project in the future.