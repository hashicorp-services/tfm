# tfm core getstate

`tfm core getstate` will iterate through all of the cloned repositories in the `github_clone_repos_path` and download the state files from the backend. 

tfm will use the locally installed terraform binary to perform `terraform init` and `terraform state pull > .terraform/pulled_terraform.tfstate` commands. It will only run these commands in the root directory of the repositories. It will only run these commands if a `.tf` file extension is detected in the root of the repository.

If tfm cannot successfully run a `terraform init` for a cloned repo tfm will return an error and continue with the next repository initilization attempt.