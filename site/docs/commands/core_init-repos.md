# tfm core init-repos

## Requirements

- Terraform community edition must be installed and in your path in the environment where this command is run.
- Credentials to authenticate to the configured backend must be configured in the environment where this command is run.

## Init Repos
`tfm core init-repos` will iterate through the cloned repositories in the `clone_repos_path` and build a `terraform_config_metadata.json` file that contains information about how the repositories are configured. It will determine all paths within the repo that contain terraform configurations and if terraform CE workspaces are being used.

Any directory containing a `.tf` file with a `backend {}` configuration contained within a `terraform {}` configuration will be added to the metadata file as a `config_path`.

```json
"repo_name": "mordor2",
    "config_paths": [
      {
        "path": "mordor2/deployments/dev",
        "workspace_info": {
          "uses_workspaces": false,
          "workspace_names": [
            "default"
          ]
        }
      },
      {
        "path": "mordor2/deployments/prod",
        "workspace_info": {
          "uses_workspaces": false,
          "workspace_names": [
            "default"
          ]
        }
      },
```

Any `config_path` that displays a terraform workspace in addition to the default one will be considered to be using terraform CE workspaces and `uses_workspaces` will be `true`. `workspace_names` will be populated with all of the terraform CE workspaces in use for the `config_path`.

```json
  {
    "repo_name": "rivendell",
    "config_paths": [
      {
        "path": "rivendell",
        "workspace_info": {
          "uses_workspaces": true,
          "workspace_names": [
            "default",
            "newrivendell",
            "oldrivendell"
          ]
        }
      }
    ]
  }
```
