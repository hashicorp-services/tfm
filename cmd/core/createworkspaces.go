// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp-services/tfm/tfclient"
	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Assuming coreCmd is your root or relevant subcommand group
var CreateWorkspacesCmd = &cobra.Command{
	Use:   "create-workspaces",
	Short: "Create TFE/TFC workspaces for each cloned repo in the github_clone_repos_path that contains a pulled_terraform.tfstate file.",
	Long:  `Create TFE/TFC workspaces for each cloned repo in the github_clone_repos_path that contains a pulled_terraform.tfstate file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		clonePath := viper.GetString("github_clone_repos_path")

		return CreateWorkspaces(tfclient.GetClientContexts(), clonePath)
	},
}

func init() {
	CoreCmd.AddCommand(CreateWorkspacesCmd)
}

// createWorkspaces iterates over directories in clonePath and creates TFE workspaces.
func CreateWorkspaces(c tfclient.ClientContexts, clonePath string) error {

	dirs, err := os.ReadDir(clonePath)
	if err != nil {
		return fmt.Errorf("error reading directories: %v", err)
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		repoPath := filepath.Join(clonePath, dir.Name())
		tfstatePath := filepath.Join(repoPath, ".terraform", "pulled_terraform.tfstate")

		if _, err := os.Stat(tfstatePath); os.IsNotExist(err) {
			continue
		}

		// Terraform state file exists, proceed to create workspace
		workspaceName := dir.Name()
		fmt.Printf("Creating workspace for repository with pulled_terraform.tfstate: %s\n", workspaceName)

		var tag []*tfe.Tag
		tag = append(tag, &tfe.Tag{Name: "tfm"})

		// Create TFE Workspace
		_, err := c.DestinationClient.Workspaces.Create(c.DestinationContext, c.DestinationOrganizationName, tfe.WorkspaceCreateOptions{
			Name: &workspaceName,
			Tags: tag,
		})
		if err != nil {
			fmt.Printf("Failed to create workspace %s: %v\n", workspaceName, err)
		}
	}

	return nil
}
