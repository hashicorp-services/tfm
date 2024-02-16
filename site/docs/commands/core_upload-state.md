# tfm core upload-state

`tfm core upload-state` is used to upload the state files that were downloaded to the created workspaces. tfm will iterate through all of the repos in the `github_clone_repos_path`. Any repositories containing a `.terraform/pulled_terraform.tfstate` file will have the state file uploaded to a workspace with a matching name.

Running this command multiple times will result in the same state file being uploaded multiple times.