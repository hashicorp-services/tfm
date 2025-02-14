# tfm core clone

`tfm core clone` will clone VCS repositories to the local system into the `clone_repos_path` you have defined in the config file.

## Supported VCS

- github
- gitlab

## Requirements

- A [supported](../migration/supported-vcs.md) `vcs_type` must be configured with one of the tfm supported VCS providers in the tfm config file.

```
vcs_type = github
```

### Credentials

tfm will used the following configuration file settings to authenticate to VCS and clone the repos:

The api token must have read access to the repositories to clone them.

#### github

```
github_token = "api token"
github_organization = "org"
github_username = "username"
```

#### gitlab

```
gitlab_token = "api token"
gitlab_group = "group102109"
gitlab_username = "username"
```

## Clone a List of Repositories

Provide a `repos_to_clone` list in the config file of repositories that you would like to clone for migration.

```hcl
repos_to_clone =  [
 "repos1",
 "repo2",
 "repo3"
]
```

## Clone All Repositories

Not providing a `repos_to_clone` list will result in tfm attempting to clone every repository in the GitHub org or GitLab group.
