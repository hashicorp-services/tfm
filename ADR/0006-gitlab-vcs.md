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

Need create a struct similiar to the following:

> [!IMPORTANT]
> Gitlab uses groups instead of organizations

```go
type ClientContext struct {
	GitLabClient       *gitlab.Client
    GitLabContext      context.Context
	GitLabToken        string
	GitLabGroup        string 
	GitLabUsername     string
}
```

### TFM Config File

The following new configurations will be added and required for use with GitLab:

- gitlab_token
- gitlab_group
- gitlab_username

### Testing

GitLab offers a free tier that can be used for nightly/weekly testing with GitLab. A service account can be created for implementing a nightly/weekly test for interaction with gitlab.

### Changes to Existing Code

- `tfm core clone` must be modified to clone repos based on VCS.
  - We can add an additional config file option `vcs_type` and use that as an input for the clone function.
- The configuration file option `clone_repos_path` needs to be modified to `clone_repo_path`.
- Many functions enforce the requirements of `github_username` `github_token` `github_organization`. Functions should me refactored to be VCS agnostic or, if required, a seprate function should be created for each VCS.

## Decision

- A new `gitlab.go` will be created as part of the `vcsclients` package and will build a context for client creation and use with Viper.
- Additional config file items will be added for GitLab support.
- `clone_repos_path` will be modified to be VCS agnostic.
- Existing functions will be made VCS agnostic.
- A new config file option `vcs_type` will be added to the config file.

## Consequences

- Users will have the ability to use GitLab with the `tfm core clone` and `tfm core remove-backend` commands.
- TFM will be made VCS agnostic for future implementation of supported VCS providers.
