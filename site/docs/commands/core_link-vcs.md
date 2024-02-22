# tfm core link-vcs

`tfm core link-vcs` will add the VCS connection to the created workspace. tfm will iterate through all of the cloned repositories in the `github_clone_repos_path` and look for a matching TFE/TFC workspace with the same name as the cloned repository directory. 

tfm will update the workspace settings using the `vcs_provider_id` defined in the config file and the name of the cloned repository directory to map the repo to the workspace.

The VCS provider must be configured in TFE/TFC and you must provide the VCS providers Oauth ID as the `vcs_provider_id` in the config file.