// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/hashicorp-services/tfm/tfclient"
	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// tfm core create-workspaces command
var CreateWorkspacesCmd = &cobra.Command{
	Use:   "create-workspaces",
	Short: "Create TFE/TFC workspaces for each cloned repo in the clone_repos_path that contains a pulled_terraform.tfstate file.",
	Long:  `Create TFE/TFC workspaces for each cloned repo in the clone_repos_path that contains a pulled_terraform.tfstate file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		clonePath := viper.GetString("clone_repos_path")

		return CreateWorkspaces(tfclient.GetDestinationClientContexts(), clonePath)
	},
}

func init() {
	CoreCmd.AddCommand(CreateWorkspacesCmd)
}

// Main function that creates the workspaces in the destination TFC/TFE organization.
func CreateWorkspaces(c tfclient.DestinationContexts, clonePath string) error {
	if c.DestinationOrganizationName == "" || c.DestinationHostname == "" || c.DestinationToken == "" {
		return fmt.Errorf("Destination TFC/TFE Organization, hostname, or token not configured.")
	}

	// The metadata file needs to have been generated using the tfm core init-repos command
	metadataFile := "terraform_config_metadata.json"
	metadata, err := loadMetadataCreateWS(metadataFile)
	if err != nil {
		return fmt.Errorf("error loading metadata: %v", err)
	}

	// Slice to store names of created workspaces for output later
	var createdWorkspaces []string

	// for each repo in the metadata file get the config paths
	for _, repoConfig := range metadata {
		for _, configPath := range repoConfig.ConfigPaths {
			basePath := strings.TrimPrefix(configPath.Path, repoConfig.RepoName+"/")

			// Replace all '/' with '-' because TFC/TFE workspace names cannot contain '/'
			basePath = strings.ReplaceAll(basePath, "/", "-")

			isRootConfigPath := basePath == "" || basePath == repoConfig.RepoName
			basePathForWorkspaceNaming := ""
			if !isRootConfigPath {
				basePathForWorkspaceNaming = "-" + basePath
			}

			// If workspaces are being used then construct the workspace name based on repo_name+config_path+workspace_name
			if configPath.WorkspaceInfo.UsesWorkspaces {
				for _, workspaceName := range configPath.WorkspaceInfo.WorkspaceNames {

					// Skip the default workspace if it's not the only one
					if workspaceName == "default" && len(configPath.WorkspaceInfo.WorkspaceNames) > 1 {
						continue
					}

					fullWorkspaceName := fmt.Sprintf("%s%s-%s", repoConfig.RepoName, basePathForWorkspaceNaming, workspaceName)

					// Special handling for root configs to avoid appending the repo name twice
					if isRootConfigPath && workspaceName == "default" {

						// Use repo name as workspace name only for default workspace at root
						fullWorkspaceName = repoConfig.RepoName
					}

					if err := createWorkspace(c, strings.Trim(fullWorkspaceName, "-")); err != nil {
						fmt.Printf("Failed to create workspace %s: %v\n", fullWorkspaceName, err)

					} else {

						// Add the successfully created workspace name to the slice for output later
						createdWorkspaces = append(createdWorkspaces, fullWorkspaceName)
					}
				}
			} else {

				// If CE workspaces aren't in use then build the workspace name using repo_name+config_path
				fullWorkspaceName := repoConfig.RepoName + basePathForWorkspaceNaming
				if err := createWorkspace(c, strings.Trim(fullWorkspaceName, "-")); err != nil {
					fmt.Printf("Failed to create workspace %s: %v\n", fullWorkspaceName, err)

					// Log the successful creation of a workspace
					o.AddDeferredMessageRead("Created workspace:", fullWorkspaceName)
				} else {

					// Add the successfully created workspace name to the slice
					createdWorkspaces = append(createdWorkspaces, fullWorkspaceName)
				}
			}
		}
	}

	// Output the names of all created workspaces
	fmt.Println("Workspaces created successfully:")
	for _, wsName := range createdWorkspaces {
		fmt.Println(wsName)
	}
	return nil
}

// Loads the metadata file information for use
func loadMetadataCreateWS(metadataFile string) ([]RepoConfig, error) {
	var metadata []RepoConfig
	file, err := ioutil.ReadFile(metadataFile)
	if err != nil {
		return nil, fmt.Errorf("error reading metadata file: %v. Run tfm core init-repos first.", err)
	}
	err = json.Unmarshal(file, &metadata)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling metadata: %v", err)
	}
	return metadata, nil
}

// This function creates the workspaces
func createWorkspace(c tfclient.DestinationContexts, workspaceName string) error {
	var tag []*tfe.Tag
	tag = append(tag, &tfe.Tag{Name: "tfm"})

	_, err := c.DestinationClient.Workspaces.Create(c.DestinationContext, c.DestinationOrganizationName, tfe.WorkspaceCreateOptions{
		Name:       &workspaceName,
		Tags:       tag,
	})
	return err
}
