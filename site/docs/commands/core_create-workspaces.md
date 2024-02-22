# tfm core create-workspaces

`tfm core create-workspaces` will iterate through all of the cloned repositories in the `github_clone_repos_path` and create  a TFC/TFE workspace in the `dst_tfc_org` defined in the config file for each repository. 

tfm only creates workspaces for repositories that it detects a `.terraform/pulled_terraform.tfstate` file in. Workspaces will be created with the same name as the github repository cloned directroy name.

If a workspace already exists with the name of the repository tfm will return an error and continue on to the next workspace creation attempt.

## Future Updates

Future updates are planned to also allow a variable set containing credentials to be assigned to the workspace the the time of creation. This will allow a plan to be run post migration to verify no changes are expected. 