# tfm core remove-backend

`tfm core remove-backend` is used to assist in removing the `backend{}` configuration block from the `terraform{}` block in terraform configurations that have been migrated.

tfm will iterate through all cloned repositories in the `github_clone_repos_path`. tfm will only look in the root of each repository. tfm will examine all files ending in a `.tf` extension and remove any instances of a `backend{}` configuration. 

tfm will create a branch, commit the branch, and push it to the origin for code owners to create a PR.

The following config file options are required to use this command:

```hcl
commit_message = "commitm message"
commit_author_name = "name"
commit_author_email = "email"
```