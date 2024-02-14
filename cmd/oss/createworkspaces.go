// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oss

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp-services/tfm/tfclient"
	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Assuming OssCmd is your root or relevant subcommand group
var createWorkspacesCmd = &cobra.Command{
	Use:   "create-workspaces",
	Short: "Create TFE workspaces for each cloned Terraform repo",
	Long:  `Iterates over all directories in the specified clone path and creates a TFE workspace for each.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		clonePath := viper.GetString("github_clone_repos_path")

		return createWorkspaces(tfclient.GetClientContexts(), clonePath)
	},
}

func init() {
	OssCmd.AddCommand(createWorkspacesCmd)
}

// createWorkspaces iterates over directories in clonePath and creates TFE workspaces.
func createWorkspaces(c tfclient.ClientContexts, clonePath string) error {

	err := filepath.Walk(clonePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != clonePath { // Skip the root clonePath itself
			// Construct path to expected terraform.tfstate file
			tfstatePath := filepath.Join(path, ".terraform", "terraform.tfstate")
			if _, err := os.Stat(tfstatePath); os.IsNotExist(err) {
				// terraform.tfstate does not exist, skip this directory
				return nil
			}
			// File exists, proceed to create workspace
			workspaceName := filepath.Base(path)
			fmt.Printf("Creating workspace for repository with terraform.tfstate: %s\n", workspaceName)

			// Create TFE Workspace
			_, err := c.DestinationClient.Workspaces.Create(c.DestinationContext, c.DestinationOrganizationName, tfe.WorkspaceCreateOptions{
				Name: &workspaceName,
			})
			if err != nil {
				fmt.Printf("Failed to create workspace %s: %v\n", workspaceName, err)

			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error iterating directories: %v", err)
	}
	return nil
}
