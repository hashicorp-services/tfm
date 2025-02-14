# tfm core remove-backend

## Requirements

- Using this feature requires the VCS token defined in the configuration file to have write permissions to the contents of the repository.
- Add the following example to your confiuration file, modifying the github values based on your [supported VCS type](../migration/supported-vcs.md):

```hcl
github_token = "api token"
github_organization = "org"
github_username = "username"
commit_message = "Remove Terraform backend configuration"
commit_author_name = "username"
commit_author_email = "user@email.com"
```

`tfm core remove-backend` is used to assist in removing the `backend{}` configuration block from the `terraform{}` block in terraform configurations that have been migrated.

tfm will use the `terraform_config_metadata.json` config file to iterate through all cloned repositories in the `clone_repos_path`. tmf will removed the backend from all `config_paths` for the repo.

tfm will create a branch, commit the branch, and push it to the origin for code owners to create a PR.

The following config file options are required to use this command:

```hcl
commit_message = "commitm message"
commit_author_name = "name"
commit_author_email = "email"
```

## Flags

`--autoapprove` Automatically approve the operation without a confirmation prompt.
`--comment` Will comment out the backend configuration instead of removing it.

## Cleanup

`tfm core cleanup` can be used to remove all cloned repos from the `clone_repos_path`
