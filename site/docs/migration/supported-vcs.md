# Supported VCS Types for TFM CE to TFC Migrations

Note: The Terraform Community Edition migration as part of `tfm` has been deprecated in favor of [tf-migrate](https://developer.hashicorp.com/terraform/cloud-docs/migrate/tf-migrate). The CE migration feature has not been removed from `tfm` however it will not be receiving further developments.

## Supported VCS Types

The following VCS types are supported values for the `vcs_type` configuration in the tfm configuration file at this time.

- github
- gitlab

## Required Github Configuration File Settings

```hcl
github_token = "api token"
github_organization = "org"
github_username = "username"
```

## Required Gitlab Configuration File Settings

```hcl
gitlab_token = "api token"
gitlab_group = "group102109"
gitlab_username = "username"
```