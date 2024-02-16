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

		return uploadStateFiles(tfclient.GetClientContexts(), clonePath)
	},
}

func init() {
	CoreCmd.AddCommand(UploadStateCmd) // Make sure coreCmd is your defined root or subgroup command
}

type TerraformState struct {
	Lineage string `json:"lineage"`
}

func uploadStateFiles(c tfclient.ClientContexts, clonePath string) error {
	// List directories directly under clonePath
	dirs, err := os.ReadDir(clonePath)
	if err != nil {
		return fmt.Errorf("error reading directories: %v", err)
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue // Skip files
		}
		path := filepath.Join(clonePath, dir.Name())
		tfstatePath := filepath.Join(path, ".terraform", "pulled_terraform.tfstate")

		if _, err := os.Stat(tfstatePath); err != nil {
			continue // Skip if terraform.tfstate does not exist
		}

		workspaceName := filepath.Base(path)
		workspace, err := c.DestinationClient.Workspaces.Read(c.DestinationContext, c.DestinationOrganizationName, workspaceName)
		if err != nil {
			fmt.Printf("Failed to read workspace %s: %v\n", workspaceName, err)
			continue // Proceed to next directory on error
		}

		_, err = c.DestinationClient.Workspaces.Lock(c.DestinationContext, workspace.ID, tfe.WorkspaceLockOptions{
			Reason: tfe.String("Uploading pulled_terraform.tfstate"),
		})
		if err != nil {
			fmt.Printf("Failed to lock workspace %s: %v\n", workspaceName, err)
			continue // Proceed to next directory on error
		}

		// Ensure workspace unlock in case of error after this point
		defer func() {
			if _, unlockErr := c.DestinationClient.Workspaces.Unlock(c.DestinationContext, workspace.ID); unlockErr != nil {
				fmt.Printf("Failed to unlock workspace %s: %v\n", workspaceName, unlockErr)
			}
		}()

		stateFileContent, err := ioutil.ReadFile(tfstatePath)
		if err != nil {
			fmt.Printf("Failed to read state file %s: %v\n", tfstatePath, err)
			continue // Proceed to next directory on error
		}

		var tfState TerraformState
		if err := json.Unmarshal(stateFileContent, &tfState); err != nil {
			fmt.Printf("Failed to unmarshal pulled_terraform.tfstate: %v\n", err)
			continue // Proceed to next directory on error
		}

		stringState := base64.StdEncoding.EncodeToString(stateFileContent)
		md5String := fmt.Sprintf("%x", md5.Sum([]byte(stateFileContent)))

		_, err = c.DestinationClient.StateVersions.Create(c.DestinationContext, workspace.ID, tfe.StateVersionCreateOptions{
			Serial:  tfe.Int64(1),
			Lineage: tfe.String(tfState.Lineage),
			MD5:     tfe.String(md5String),
			State:   tfe.String(stringState),
		})
		if err != nil {
			fmt.Printf("Failed to upload state file to workspace %s: %v\n", workspaceName, err)
			continue // Proceed to next directory on error
		}

		fmt.Printf("Successfully uploaded pulled_terraform.tfstate to workspace: %s\n", workspaceName)
	}

	return nil
}
