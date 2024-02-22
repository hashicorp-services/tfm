# Open Source / Community Edition to TFC/TFE Feature(s)

## Phases
This feature will be implemented in phases

### Phase 1
- GitHub only support
- Only migration of VCS repositories with terraform configurations in the ROOT of the repository will be supported
- No support for 3rd party (terragrunt) managed workspaces
- No support for monorepos with a directroy structure containing multiple terraform configurations
- No support for configurations using the terraform workspace feature to use 1 backend for multiple stat files

### Phase 2
- Support for monorepos with multiple directories containing different terraform configurations and backend blocks.

Ideas:

- Modify the getstate command to recursively walk through each repository's directory structure.
- For each directory, check if it contains .tf files. If it does, treat it as a Terraform configuration directory.
- Execute terraform init and terraform state pull within each Terraform configuration directory found.
- Save the state files with a path or naming convention that reflects their directory structure within the repository, ensuring uniqueness.
- When creating workspaces, use the information collected during the getstate command execution to determine the correct working directories for each workspace.
- Ensure the workspace creation process includes setting the workspace's working directory to match the subdirectory structure where the Terraform configuration was found.
- This might involve modifying the naming convention for workspaces to reflect their structure within the repository (e.g., repo-deployments-prod).
- Make sure the VCS can still be connceted to the proper workspace with the working directroy set proper.

- Configuration File or Flags: Allowing users to specify paths to search for Terraform configurations within repositories via a configuration file or command-line flags. This could help handle cases with non-standard structures without hardcoding paths.

```hcl
terraform_paths = [
  "deployments/prod",
  "deployments/dev",
  "infra"
]
```

```go
type Config struct {
    TerraformPaths []string `mapstructure:"terraform_paths"`
}
```

- Metadata Tracking: Implement a system to track metadata about each discovered Terraform configuration, such as its path within the repository and the associated workspace name. This metadata will be crucial for linking workspaces with VCS and setting correct working directories.

```go
type TerraformConfigMetadata struct {
    RepoPath       string
    ConfigPath     string
    WorkspaceName  string
}

var discoveredConfigs []TerraformConfigMetadata

discoveredConfigs = append(discoveredConfigs, TerraformConfigMetadata{
    RepoPath:       "path/to/repo",
    ConfigPath:     "deployments/prod",
    WorkspaceName:  "repo-deployments-prod",
})
```

- Flexible Workspace Naming: Develop a flexible and consistent naming scheme for workspaces that can accommodate various directory structures, potentially incorporating the repository name and path to the Terraform configuration.

```go
for _, config := range discoveredConfigs {
    workspaceName := generateWorkspaceName(config.RepoPath, config.ConfigPath)
    // Use workspaceName to create the workspace
}

// generateWorkspaceName creates a workspace name based on the repository path and Terraform configuration path
func generateWorkspaceName(repoName, configPath string) string {
    if configPath == "/" || configPath == "" {
        // Root-level configuration; no specific path needed in the name
        return repoName
    }
    // For configurations in subdirectories, include the path in the name
    configPathFlat := strings.ReplaceAll(strings.Trim(configPath, "/"), "/", "-")
    return fmt.Sprintf("%s-%s", repoName, configPathFlat)
}
```

- Adjust the VCS command to support the new workspace nameing convention

```go
// Placeholder function to list all workspaces. Implement accordingly based on your setup.
func listWorkspaces() ([]Workspace, error) {
    // Implementation to call the Terraform Cloud API and list workspaces
    return nil, nil // Return actual workspaces and error
}

// Workspace structure based on Terraform Cloud API response
type Workspace struct {
    Name string
}

// parseWorkspaceName parses the workspace name to extract the repository name and config path
func parseWorkspaceName(workspaceName string) (string, string) {
    parts := strings.SplitN(workspaceName, "-", 2)
    if len(parts) < 2 {
        return workspaceName, ""
    }
    repoName := parts[0]
    configPath := strings.ReplaceAll(parts[1], "-", "/")
    return repoName, configPath
}

func linkVCSWorkspaces() error {
    workspaces, err := listWorkspaces()
    if err != nil {
        return err
    }

    for _, workspace := range workspaces {
        repoName, configPath := parseWorkspaceName(workspace.Name)
        // Repository name (repoName) and the Terraform configuration path (configPath)
        // Use this information to match the workspace with its VCS repository and set the working directory
    }
    return nil
}
```

### Phase 3
- Support for terraform configurations using the terraform workspaces command'

Ideas:

- How do you link workspaces to these state files? You can't. This requires an entire refactor of the terraform configuration repository to provide a separate working directory for each state.
- How to help customers refactor repos quickly?
- terraform workspace list
- terraform workspace select <workspace>
- terraform init
- terraform state pull

Refactoring for Terraform Cloud:

If your current setup heavily relies on Terraform CLI workspaces to manage different environments within the same configuration, transitioning to Terraform Cloud might require some refactoring.
Ideally, each Terraform Cloud workspace should be aligned with a specific environment or configuration set. This might mean splitting up configurations that are currently managed with CLI workspaces into separate directories or repositories, each linked to its own Terraform Cloud workspace.

To effectively use TFC workspaces with a VCS-linked repository, consider organizing your repository to clearly separate different sets of configurations that will be managed as distinct environments in Terraform Cloud.
This might mean restructuring your repository so that each major environment or set of configurations has its own directory, which can then be directly mapped to a TFC workspace with a corresponding working directory.

### Phase 4
- Support for additional VCSs.
  - Targetting GitLab
- Support to allow workspace renaming during migration

### Phase 5
- Continued development