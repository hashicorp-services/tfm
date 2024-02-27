// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package core

import (
	"fmt"

	"github.com/spf13/cobra"
)

var includeRemoveBackend bool
var includeCommands []string

var validIncludeCommands = map[string]bool{
	"remove-backend": true,
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrates opensource/community edition Terraform code and state to TFE/TFC in 1 continuous workflow.",
	Long: `Executes a sequence of commands to clone repositories, get state, create workspaces, 
upload state, link VCS, and optionally remove backend configurations as part of the core migration process.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate included commands before executing the migration process
		for _, includeCmd := range includeCommands {
			if _, valid := validIncludeCommands[includeCmd]; !valid {
				return fmt.Errorf("invalid command specified in --include: %s", includeCmd)
			}
		}

		// Check for auto-approval
		if !autoApprove {
			promptMessage := `
This command will run all of the commands listed below in order.:
	1. tfm core clone
	2. tfm core init-repos
	3. tfm core getstate
	4. tfm core create-workspaces
	5. tfm core upload-state
	6. tfm core link-vcs
	7. (Optional) If you provided the flag --include remove-backend then tfm core remove-backen will run.
			
	Are you sure you want to proceed? Type 'yes' to continue: `
			o.AddPassUserProvided(promptMessage)
			var response string
			_, err := fmt.Scanln(&response)
			if err != nil || response != "yes" {
				fmt.Println("Operation aborted by the user.")
				return nil // Exit if the user does not confirm
			}
		}


		commonArgs := []string{}

		// Directly invoke the RunE function of each command
		if err := CloneCmd.RunE(cmd, commonArgs); err != nil {
			return err
		}
		if err := InitReposCmd.RunE(cmd, commonArgs); err != nil {
			return err
		}
		if err := GetStateCmd.RunE(cmd, commonArgs); err != nil {
			return err
		}
		if err := CreateWorkspacesCmd.RunE(cmd, commonArgs); err != nil {
			return err
		}
		if err := UploadStateCmd.RunE(cmd, commonArgs); err != nil {
			return err
		}
		if err := LinkVCSCmd.RunE(cmd, commonArgs); err != nil {
			return err
		}

		// Dynamically execute additional commands based on --include flag
		for _, includeCmd := range includeCommands {
			switch includeCmd {
			case "remove-backend":
				RemoveBackendCmd.Flags().Set("autoapprove", "true")
				if err := RemoveBackendCmd.RunE(cmd, []string{}); err != nil {
					return err
				}
			}
		}

		return nil
	},
}

func init() {
	CoreCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().BoolVar(&autoApprove, "autoapprove", false, "Automatically approve the operation without a confirmation prompt")
	migrateCmd.Flags().StringSliceVar(&includeCommands, "include", nil, "Specify additional commands to include in the migration process (e.g., --include remove-backend)")
}
