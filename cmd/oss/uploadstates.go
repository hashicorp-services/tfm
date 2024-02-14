// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oss

import (
	"crypto/md5"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hashicorp-services/tfm/tfclient"
	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var uploadStateCmd = &cobra.Command{
	Use:   "upload-state",
	Short: "Upload terraform.tfstate files to TFE workspaces",
	Long: `Iterates over directories containing .terraform/terraform.tfstate files, 
           finds corresponding TFE workspaces, locks the workspace, uploads the state file, 
           and then unlocks the workspace.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		clonePath := viper.GetString("github_clone_repos_path") // Ensure this is set

		return uploadStateFiles(tfclient.GetClientContexts(), clonePath)
	},
}

func init() {
	OssCmd.AddCommand(uploadStateCmd) // Make sure OssCmd is your defined root or subgroup command
}

type TerraformState struct {
	Lineage string `json:"lineage"`
}

func uploadStateFiles(c tfclient.ClientContexts, clonePath string) error {

	err := filepath.Walk(clonePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != clonePath {
			tfstatePath := filepath.Join(path, ".terraform", "terraform.tfstate")
			if _, err := os.Stat(tfstatePath); err != nil {
				return nil // Skip if terraform.tfstate does not exist
			}

			workspaceName := filepath.Base(path)
			// Ensure workspace exists
			workspace, err := c.DestinationClient.Workspaces.Read(c.DestinationContext, c.DestinationOrganizationName, workspaceName)
			if err != nil {
				return fmt.Errorf("failed to read workspace %s: %v", workspaceName, err)
			}

			// Lock the workspace
			_, err = c.DestinationClient.Workspaces.Lock(c.DestinationContext, workspace.ID, tfe.WorkspaceLockOptions{
				Reason: tfe.String("Uploading terraform.tfstate"),
			})
			if err != nil {
				return fmt.Errorf("failed to lock workspace %s: %v", workspaceName, err)
			}

			// Schedule workspace unlock before handling any other error
			defer func() {
				if _, unlockErr := c.DestinationClient.Workspaces.Unlock(c.DestinationContext, workspace.ID); unlockErr != nil {
					fmt.Printf("Failed to unlock workspace %s: %v\n", workspaceName, unlockErr)
				}
			}()

			// Read terraform.tfstate file
			stateFileContent, err := ioutil.ReadFile(tfstatePath)
			if err != nil {
				return fmt.Errorf("failed to read state file %s: %v", tfstatePath, err)
			}

			// Unmarshal the state file to extract lineage
			var tfState TerraformState
			if err := json.Unmarshal(stateFileContent, &tfState); err != nil {
				return fmt.Errorf("failed to unmarshal terraform.tfstate: %v", err)
			}

			// Base64 encode the state as a string
			stringState := b64.StdEncoding.EncodeToString(stateFileContent)

			// Get the MD5 hash of the state
			md5String := fmt.Sprintf("%x", md5.Sum([]byte(stateFileContent)))

			// Upload state file to the workspace
			_, err = c.DestinationClient.StateVersions.Create(c.DestinationContext, workspace.ID, tfe.StateVersionCreateOptions{
				Type:             "",
				Serial:           tfe.Int64(1),
				Lineage:          tfe.String(tfState.Lineage),
				MD5:              tfe.String(md5String),
				State:            tfe.String(string(stringState)),
				Force:            new(bool),
				JSONState:        new(string),
				JSONStateOutputs: new(string),
			})
			if err != nil {
				return fmt.Errorf("failed to upload state file to workspace %s: %v", workspaceName, err)
			}

			fmt.Printf("Successfully uploaded terraform.tfstate to workspace: %s\n", workspaceName)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error processing directories: %v", err)
	}
	return nil
}
