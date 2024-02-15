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

var LinkVCSCmd = &cobra.Command{
	Use:   "link-vcs",
	Short: "Link repos to TFE workspaces via VCS",
	Long:  `Iterates over cloned repositories containing .terraform/terraform.tfstate files, finds the corresponding TFE workspace, and links it to the GitHub repository.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		clonePath := viper.GetString("github_clone_repos_path") // Ensure this is set
		if clonePath == "" {
			clonePath = "test" // Default path if not specified
		}
		return LinkVCS(tfclient.GetClientContexts(), clonePath)
	},
}

func init() {
	OssCmd.AddCommand(LinkVCSCmd) // Make sure OssCmd is your defined root or subgroup command
}

func LinkVCS(c tfclient.ClientContexts, clonePath string) error {
	vcsProviderID := viper.GetString("vcs_provider_id")
	githubOrganization := viper.GetString("github_organization")

	if c.DestinationOrganizationName == "" || vcsProviderID == "" {
		return fmt.Errorf("TFE organization or VCS provider ID not specified in configuration")
	}

	o.AddFormattedMessageCalculated3("Using VCS Provider %s in org %s\n", vcsProviderID, c.DestinationOrganizationName)

	err := filepath.Walk(clonePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != clonePath {
			tfstatePath := filepath.Join(path, ".terraform", "terraform.tfstate")
			if _, err := os.Stat(tfstatePath); err != nil {
				return nil
			}

			workspaceName := filepath.Base(path)
			// Construct the GitHub repo URL; adjust as needed based on your repo naming conventions
			repoIdentifier := fmt.Sprintf("%s/%s", githubOrganization, workspaceName)

			o.AddFormattedMessageCalculated3("Linking repo %s to workspace %s.\n", repoIdentifier, workspaceName)

			workspace, err := c.DestinationClient.Workspaces.Read(c.DestinationContext, c.DestinationOrganizationName, workspaceName)
			if err != nil {
				fmt.Printf("Workspace %s not found: %v\n", workspaceName, err)
				return nil
			}

			// Attempt to create a VCS connection
			_, err = c.DestinationClient.Workspaces.Update(c.DestinationContext, c.DestinationOrganizationName, workspace.Name, tfe.WorkspaceUpdateOptions{
				Type: "",
				VCSRepo: &tfe.VCSRepoOptions{
					Branch:            new(string),
					Identifier:        &repoIdentifier,
					IngressSubmodules: new(bool),
					OAuthTokenID:      &vcsProviderID,
					TagsRegex:         new(string),
					GHAInstallationID: new(string),
				},
			})
			if err != nil {
				fmt.Printf("Failed to link VCS for workspace %s: %v\n", workspaceName, err)
			} else {
				fmt.Printf("Successfully linked %s to workspace %s\n", repoIdentifier, workspaceName)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error processing directories: %v", err)
	}
	return nil
}
