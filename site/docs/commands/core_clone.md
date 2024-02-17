# tfm core clone

`tfm core clone` will clone VCS repositories to the local system into the `github_clone_repos_path` you have defined in the config file.

## Supported VCS

- Github

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

Not providing a `repos_to_clone` list will result in tfm attempting to clone every repository in the GitHub org.

## Credentials

tfm will used the following configuration file settings to authenticate to GitHub and clone the repos:

The api token must have read access to the repositories to clone them.

```
github_token = "api token"
github_organization = "org"
github_username = "username"
```