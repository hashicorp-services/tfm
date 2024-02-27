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

var LinkVCSCmd = &cobra.Command{
	Use:   "link-vcs",
	Short: "Link repos in the github_clone_repos_path to their corresponding workspaces in TFE/TFC.",
	Long:  `Iterates over cloned repositories containing .terraform/pulled_terraform.tfstate files, finds the corresponding TFE/TFC workspace, and links it to the GitHub repository.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		clonePath := viper.GetString("github_clone_repos_path") // Ensure this is set
		if clonePath == "" {
			clonePath = "test" // Default path if not specified
		}
		return LinkVCS(tfclient.GetDestinationClientContexts(), clonePath)
	},
}

func init() {
	CoreCmd.AddCommand(LinkVCSCmd)
}


// Loads the metadata file information for use
func loadMetadataLinkVcs(metadataFile string) ([]RepoConfig, error) {
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

// Constructs the workspace names with the same function used in tfm core create-workspaces for mapping of workspaces to VCS repos
func constructWorkspaceNames(repoConfig RepoConfig, configPath ConfigPathInfo) []string {
	var workspaceNames []string
	basePath := strings.TrimPrefix(configPath.Path, repoConfig.RepoName+"/")

	// TFC/TFE workspace names cannot contain '/' so we remove it
	basePath = strings.ReplaceAll(basePath, "/", "-")

	isRootConfigPath := basePath == "" || basePath == repoConfig.RepoName
	basePathForWorkspaceNaming := ""
	if !isRootConfigPath {
		basePathForWorkspaceNaming = "-" + basePath
	}

	if configPath.WorkspaceInfo.UsesWorkspaces {
		for _, workspaceName := range configPath.WorkspaceInfo.WorkspaceNames {

			// Skip the default workspace if it's not the only one
			if workspaceName == "default" && len(configPath.WorkspaceInfo.WorkspaceNames) > 1 {
				continue
			}

			var fullWorkspaceName string

			// Use repo name directly for default workspace at root
			if isRootConfigPath && workspaceName == "default" {
				fullWorkspaceName = repoConfig.RepoName
			} else {
				fullWorkspaceName = fmt.Sprintf("%s%s-%s", repoConfig.RepoName, basePathForWorkspaceNaming, workspaceName)
			}
			workspaceNames = append(workspaceNames, strings.Trim(fullWorkspaceName, "-"))
		}
	} else {

		// For configurations that do not use multiple workspaces
		fullWorkspaceName := repoConfig.RepoName + basePathForWorkspaceNaming
		workspaceNames = append(workspaceNames, strings.Trim(fullWorkspaceName, "-"))
	}

	return workspaceNames
}

func LinkVCS(c tfclient.DestinationContexts, clonePath string) error {

	// Check for necessary configurations
	metadata, err := loadMetadataLinkVcs("terraform_config_metadata.json")
	if err != nil {
		return fmt.Errorf("error loading metadata: %v", err)
	}

	vcsProviderID := viper.GetString("vcs_provider_id")
	if vcsProviderID == "" {
		return fmt.Errorf("vcs_provider_id is not configured")
	}

	githubOrganization := viper.GetString("github_organization")
	if githubOrganization == "" {
		return fmt.Errorf("github_organization is not configured")
	}

	for _, repoConfig := range metadata {
		for _, configPath := range repoConfig.ConfigPaths {
			workspaceNames := constructWorkspaceNames(repoConfig, configPath)
			for _, workspaceName := range workspaceNames {
				repoIdentifier := fmt.Sprintf("%s/%s", githubOrganization, repoConfig.RepoName)
				workingDirectory := ""

				// Only set the working directory if config_path is not the root of the repo
				if !(configPath.Path == repoConfig.RepoName || strings.Trim(configPath.Path, "/") == "") {
					workingDirectory = strings.TrimPrefix(configPath.Path, repoConfig.RepoName+"/")
				}

				fmt.Printf("Linking repo %s to workspace %s with working directory %s.\n", repoIdentifier, workspaceName, workingDirectory)

				workspace, err := c.DestinationClient.Workspaces.Read(c.DestinationContext, c.DestinationOrganizationName, workspaceName)
				if err != nil {
					fmt.Printf("Workspace %s not found: %v\n", workspaceName, err)
					continue
				}

				updateOptions := tfe.WorkspaceUpdateOptions{
					VCSRepo: &tfe.VCSRepoOptions{
						Identifier:   &repoIdentifier,
						OAuthTokenID: &vcsProviderID,
					},
				}

				// Conditionally add the WorkingDirectory if it's not empty
				if workingDirectory != "" {
					updateOptions.WorkingDirectory = &workingDirectory
				}

				_, err = c.DestinationClient.Workspaces.Update(c.DestinationContext, c.DestinationOrganizationName, workspace.Name, updateOptions)
				if err != nil {
					fmt.Printf("[ERROR] Failed to link VCS for workspace %s: %v\n", workspaceName, err)
				} else {
					fmt.Printf("Successfully linked %s to workspace %s with working directory %s\n", repoIdentifier, workspaceName, workingDirectory)
				}
			}
		}
	}

	return nil
}
