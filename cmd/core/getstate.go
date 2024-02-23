// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (

	// `tfm core getstate` command
	GetStateCmd = &cobra.Command{
		Use:   "getstate",
		Short: "Initialize and get state from terraform repos in the github_clone_repos_path.",
		Long:  "Initialize and get state from terraform repos in the github_clone_repos_path.",
		RunE: func(cmd *cobra.Command, args []string) error {

			return initializeRepos()
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			o.Close()
		},
	}
)

func init() {

	GetStateCmd.Flags().SetInterspersed(false)

	// Add commands
	CoreCmd.AddCommand(GetStateCmd)
}

// Runs terraform init. Terraform must be installed and in the path
func runTerraformInit(dirPath string) error {
	cmd := exec.Command("terraform", "init")
	cmd.Dir = dirPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Runs terraform state pull. Terraform must be installed and in the path
func pullTerraformState(dirPath, outputPath string) error {
	cmd := exec.Command("terraform", "state", "pull")
	cmd.Dir = dirPath
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(outputPath, output, 0644)
}

func selectTerraformWorkspace(dirPath, ceWorkspaceName string) error {
	cmd := exec.Command("terraform", "workspace", "select", ceWorkspaceName)
	cmd.Dir = dirPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type repoConfig struct {
	RepoName    string           `json:"repo_name"`
	ConfigPaths []ConfigPathInfo `json:"config_paths"`
}

type configPathInfo struct {
	Path          string        `json:"path"`
	WorkspaceInfo WorkspaceInfo `json:"workspace_info"`
}

type workspaceInfo struct {
	UsesWorkspaces bool     `json:"uses_workspaces"`
	WorkspaceNames []string `json:"workspace_names,omitempty"`
}

func loadMetadata() ([]RepoConfig, error) {
	metadataFile := "terraform_config_metadata.json"
	file, err := os.Open(metadataFile)
	if err != nil {
		return nil, fmt.Errorf("error opening metadata file: %v. Run tfm core init-repos first.", err)
	}
	defer file.Close()

	var repoConfigs []RepoConfig
	err = json.NewDecoder(file).Decode(&repoConfigs)
	if err != nil {
		return nil, fmt.Errorf("error decoding metadata file: %v", err)
	}

	return repoConfigs, nil
}

func initializeRepos() error {
	repoConfigs, err := loadMetadata()
	if err != nil {
		return err
	}

	clonePath := viper.GetString("github_clone_repos_path")
	if clonePath == "" {
		return fmt.Errorf("github_clone_repos_path is not configured")
	}

	var initCount int

	for _, repoConfig := range repoConfigs {
		for _, configPathInfo := range repoConfig.ConfigPaths {
			// fullPath is directly constructed from clonePath and configPathInfo.Path
			fullPath := filepath.Join(clonePath, configPathInfo.Path)

			fmt.Printf("Initializing Terraform in: %s\n", fullPath)
			if err := runTerraformInit(fullPath); err != nil {
				fmt.Printf("Failed to initialize Terraform in %s: %v\n", fullPath, err)
				continue
			}

			if configPathInfo.WorkspaceInfo.UsesWorkspaces {
				for _, workspace := range configPathInfo.WorkspaceInfo.WorkspaceNames {
					if workspace == "default" && len(configPathInfo.WorkspaceInfo.WorkspaceNames) == 1 {
						continue // Skip default workspace handling if it's the only one and uses_workspaces is true
					}
					fmt.Printf("Handling workspace '%s' in %s\n", workspace, fullPath)

					// Conditional logic based on workspace presence
					if workspace != "default" {
						if err := selectTerraformWorkspace(fullPath, workspace); err != nil {
							fmt.Printf("Failed to select workspace '%s' in %s: %v\n", workspace, fullPath, err)
							continue
						}
					}

					stateFilePath := filepath.Join(fullPath, fmt.Sprintf(".terraform/pulled_%s_terraform.tfstate", workspace))
					if err := pullTerraformState(fullPath, stateFilePath); err != nil {
						fmt.Printf("Failed to pull Terraform state for workspace '%s' in %s: %v\n", workspace, fullPath, err)
						continue
					}
				}
			} else {
				// Handle non-workspace scenario
				stateFilePath := filepath.Join(fullPath, ".terraform/pulled_terraform.tfstate")
				if err := pullTerraformState(fullPath, stateFilePath); err != nil {
					fmt.Printf("Failed to pull Terraform state in %s: %v\n", fullPath, err)
					continue
				}
			}

			initCount++
		}
	}

	fmt.Printf("Terraform initialization and state processing completed for %d configurations.\n", initCount)
	return nil
}
