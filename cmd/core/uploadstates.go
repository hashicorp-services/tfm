// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package core

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp-services/tfm/tfclient"
	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var UploadStateCmd = &cobra.Command{
	Use:   "upload-state",
	Short: "Upload .terraform/pulled_terraform.tfstate files from repos cloned into the github_clone_repos_path to TFE/TFC workspaces.",
	Long: `Iterates over directories containing .terraform/pulled_terraform.tfstate files, 
           finds corresponding TFE workspaces, locks the workspace, uploads the state file, 
           and then unlocks the workspace.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		clonePath := viper.GetString("github_clone_repos_path") // Ensure this is set

		return uploadStateFiles(tfclient.GetDestinationClientContexts(), clonePath)
	},
}

func init() {
	CoreCmd.AddCommand(UploadStateCmd)
}

type TerraformState struct {
	Lineage string `json:"lineage"`
}

// Loads the metadata file information for use
func loadMetadataUploadState(metadataFile string) ([]RepoConfig, error) {
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

// // Constructs the workspace names with the same function used in tfm core create-workspaces for mapping of workspaces to VCS repos
// func constructWorkspaceNames2(repoConfig RepoConfig, configPath ConfigPathInfo) []string {
// 	var workspaceNames []string
// 	basePath := strings.TrimPrefix(configPath.Path, repoConfig.RepoName+"/")

// 	// TFC/TFE workspace names cannot contain '/' so we remove it
// 	basePath = strings.ReplaceAll(basePath, "/", "-")

// 	isRootConfigPath := basePath == "" || basePath == repoConfig.RepoName
// 	basePathForWorkspaceNaming := ""
// 	if !isRootConfigPath {
// 		basePathForWorkspaceNaming = "-" + basePath
// 	}

// 	if configPath.WorkspaceInfo.UsesWorkspaces {
// 		for _, workspaceName := range configPath.WorkspaceInfo.WorkspaceNames {

// 			// Skip the default workspace if it's not the only one
// 			if workspaceName == "default" && len(configPath.WorkspaceInfo.WorkspaceNames) > 1 {
// 				continue
// 			}

// 			var fullWorkspaceName string

// 			// Use repo name directly for default workspace at root
// 			if isRootConfigPath && workspaceName == "default" {
// 				fullWorkspaceName = repoConfig.RepoName
// 			} else {
// 				fullWorkspaceName = fmt.Sprintf("%s%s-%s", repoConfig.RepoName, basePathForWorkspaceNaming, workspaceName)
// 			}
// 			workspaceNames = append(workspaceNames, strings.Trim(fullWorkspaceName, "-"))
// 		}
// 	} else {

// 		// For configurations that do not use multiple workspaces
// 		fullWorkspaceName := repoConfig.RepoName + basePathForWorkspaceNaming
// 		workspaceNames = append(workspaceNames, strings.Trim(fullWorkspaceName, "-"))
// 	}

// 	return workspaceNames
// }

// func uploadStateFiles(c tfclient.DestinationContexts, clonePath string) error {
// 	if c.DestinationOrganizationName == "" || c.DestinationHostname == "" || c.DestinationToken == "" {
// 		return fmt.Errorf("Destination TFC/TFE Organization, hostname, or token not configured.")
// 	}

// 	metadata, err := loadMetadataUploadState("terraform_config_metadata.json")
// 	if err != nil {
// 		return fmt.Errorf("error loading metadata: %v", err)
// 	}

// 	for _, repoConfig := range metadata {
// 		for _, configPath := range repoConfig.ConfigPaths {
// 			workspaceNames := constructWorkspaceNames2(repoConfig, configPath)
// 			for _, workspaceName := range workspaceNames {
// 				// Dynamically adjust state file name based on workspace usage
// 				stateFileName := "pulled_terraform.tfstate"
// 				if configPath.WorkspaceInfo.UsesWorkspaces {
// 					for _, wsName := range configPath.WorkspaceInfo.WorkspaceNames {
// 						if wsName != "default" || len(configPath.WorkspaceInfo.WorkspaceNames) == 1 {
// 							stateFileName = fmt.Sprintf("pulled_%s_terraform.tfstate", wsName)
// 						}
// 					}
// 				}

// 				// Adjust the path construction to avoid repo name duplication
// 				tfstatePath := ""
// 				if strings.HasPrefix(configPath.Path, repoConfig.RepoName) {
// 					// If configPath already includes the repoName, do not duplicate it
// 					tfstatePath = filepath.Join(clonePath, configPath.Path, ".terraform", stateFileName)
// 				} else {
// 					tfstatePath = filepath.Join(clonePath, repoConfig.RepoName, configPath.Path, ".terraform", stateFileName)
// 				}

// 				if _, err := os.Stat(tfstatePath); err != nil {
// 					fmt.Printf("State file does not exist: %s\n", tfstatePath)
// 					continue
// 				}

// 				workspace, err := c.DestinationClient.Workspaces.Read(c.DestinationContext, c.DestinationOrganizationName, workspaceName)
// 				if err != nil {
// 					fmt.Printf("Failed to read workspace %s: %v\n", workspaceName, err)
// 					continue
// 				}

// 				// Lock the workspace before uploading state
// 				_, err = c.DestinationClient.Workspaces.Lock(c.DestinationContext, workspace.ID, tfe.WorkspaceLockOptions{
// 					Reason: tfe.String("Uploading state file"),
// 				})
// 				if err != nil {
// 					fmt.Printf("Failed to lock workspace %s: %v\n", workspaceName, err)
// 					continue
// 				}
// 				// Ensure workspace unlock in case of error after this point
// 				defer func(workspaceID string) {
// 					if _, unlockErr := c.DestinationClient.Workspaces.Unlock(c.DestinationContext, workspaceID); unlockErr != nil {
// 						fmt.Printf("Failed to unlock workspace %s: %v\n", workspaceName, unlockErr)
// 					}
// 				}(workspace.ID)

// 				stateFileContent, err := ioutil.ReadFile(tfstatePath)
// 				if err != nil {
// 					fmt.Printf("Failed to read state file %s: %v\n", tfstatePath, err)
// 					continue
// 				}

// 				var tfState TerraformState
// 				if err := json.Unmarshal(stateFileContent, &tfState); err != nil {
// 					fmt.Printf("Failed to unmarshal state file: %v\n", err)
// 					continue
// 				}

// 				stringState := base64.StdEncoding.EncodeToString(stateFileContent)
// 				md5String := fmt.Sprintf("%x", md5.Sum([]byte(stateFileContent)))

// 				// Upload state file to the workspace
// 				_, err = c.DestinationClient.StateVersions.Create(c.DestinationContext, workspace.ID, tfe.StateVersionCreateOptions{
// 					Serial:  tfe.Int64(1),
// 					Lineage: tfe.String(tfState.Lineage),
// 					MD5:     tfe.String(md5String),
// 					State:   tfe.String(stringState),
// 				})
// 				if err != nil {
// 					fmt.Printf("Failed to upload state file to workspace %s: %v\n", workspaceName, err)
// 				} else {
// 					o.AddFormattedMessageCalculated2("Successfully uploaded state file %s to workspace: %s\n", tfstatePath, workspaceName)
// 				}
// 			}
// 		}
// 	}

// 	return nil
// }

func constructWorkspaceNameAndStatePath(repoConfig RepoConfig, configPath ConfigPathInfo, wsName string, clonePath string) (string, string) {
	var workspaceName, stateFileName, tfstatePath string
	basePath := strings.TrimPrefix(configPath.Path, repoConfig.RepoName+"/")
	basePath = strings.ReplaceAll(basePath, "/", "-")

	isRootConfigPath := basePath == "" || basePath == repoConfig.RepoName
	if isRootConfigPath && wsName != "default" {
		// Root config with a workspace that is not 'default'
		workspaceName = fmt.Sprintf("%s-%s", repoConfig.RepoName, wsName)
		stateFileName = fmt.Sprintf("pulled_%s_terraform.tfstate", wsName)
	} else if !isRootConfigPath {
		if configPath.WorkspaceInfo.UsesWorkspaces && wsName != "default" {
			// Non-root config with specific workspace
			workspaceName = fmt.Sprintf("%s-%s-%s", repoConfig.RepoName, basePath, wsName)
			stateFileName = fmt.Sprintf("pulled_%s_terraform.tfstate", wsName)
		} else {
			// Non-root config without specific workspaces or 'default' workspace
			workspaceName = fmt.Sprintf("%s-%s", repoConfig.RepoName, basePath)
			stateFileName = "pulled_terraform.tfstate"
		}
	} else {
		// Root config with 'default' workspace
		workspaceName = repoConfig.RepoName
		stateFileName = "pulled_terraform.tfstate"
	}

	// Construct the path to the state file
	if strings.HasPrefix(configPath.Path, repoConfig.RepoName) {
		tfstatePath = filepath.Join(clonePath, configPath.Path, ".terraform", stateFileName)
	} else {
		tfstatePath = filepath.Join(clonePath, repoConfig.RepoName, configPath.Path, ".terraform", stateFileName)
	}

	return workspaceName, tfstatePath
}

func uploadStateFiles(c tfclient.DestinationContexts, clonePath string) error {
	if c.DestinationOrganizationName == "" || c.DestinationHostname == "" || c.DestinationToken == "" {
		return fmt.Errorf("Destination TFC/TFE Organization, hostname, or token not configured.")
	}

	metadata, err := loadMetadataUploadState("terraform_config_metadata.json")
	if err != nil {
		return fmt.Errorf("error loading metadata: %v", err)
	}

	for _, repoConfig := range metadata {
		for _, configPath := range repoConfig.ConfigPaths {
			for _, wsName := range configPath.WorkspaceInfo.WorkspaceNames {
				// Skip default workspace if it's not the only one.
				if wsName == "default" && len(configPath.WorkspaceInfo.WorkspaceNames) > 1 {
					continue
				}

				workspaceName, tfstatePath := constructWorkspaceNameAndStatePath(repoConfig, configPath, wsName, clonePath)
				if _, err := os.Stat(tfstatePath); err != nil {
					fmt.Printf("State file does not exist: %s\n", tfstatePath)
					continue
				}

				workspace, err := c.DestinationClient.Workspaces.Read(c.DestinationContext, c.DestinationOrganizationName, workspaceName)
				if err != nil {
					fmt.Printf("Failed to read workspace %s: %v\n", workspaceName, err)
					continue
				}

				_, err = c.DestinationClient.Workspaces.Lock(c.DestinationContext, workspace.ID, tfe.WorkspaceLockOptions{Reason: tfe.String("Uploading state file")})
				if err != nil {
					fmt.Printf("Failed to lock workspace %s: %v\n", workspaceName, err)
					continue
				}

				defer func() {
					if _, unlockErr := c.DestinationClient.Workspaces.Unlock(c.DestinationContext, workspace.ID); unlockErr != nil {
						fmt.Printf("Failed to unlock workspace %s: %v\n", workspaceName, unlockErr)
					}
				}()

				stateFileContent, err := ioutil.ReadFile(tfstatePath)
				if err != nil {
					fmt.Printf("Failed to read state file %s: %v\n", tfstatePath, err)
					continue
				}

				var tfState TerraformState
				if err = json.Unmarshal(stateFileContent, &tfState); err != nil {
					fmt.Printf("Failed to unmarshal state file: %v\n", err)
					continue
				}

				stringState := base64.StdEncoding.EncodeToString(stateFileContent)
				md5String := fmt.Sprintf("%x", md5.Sum(stateFileContent))

				_, err = c.DestinationClient.StateVersions.Create(c.DestinationContext, workspace.ID, tfe.StateVersionCreateOptions{
					Serial:  tfe.Int64(1),
					Lineage: tfe.String(tfState.Lineage),
					MD5:     tfe.String(md5String),
					State:   tfe.String(stringState),
				})
				if err != nil {
					fmt.Printf("Failed to upload state file to workspace %s: %v\n", workspaceName, err)
					continue
				}

				fmt.Printf("Successfully uploaded state file %s to workspace: %s\n", tfstatePath, workspaceName)
			}
		}
	}

	return nil
}
