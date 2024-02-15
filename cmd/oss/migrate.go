// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oss

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
	Short: "Migrates opensource/community edition Terraform code and state to TFE/TFC.",
	Long: `Executes a sequence of commands to clone repositories, get state, create workspaces, 
upload state, link VCS, and optionally remove backend configurations as part of the OSS migration process.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate included commands before executing the migration process
		for _, includeCmd := range includeCommands {
			if _, valid := validIncludeCommands[includeCmd]; !valid {
				return fmt.Errorf("invalid command specified in --include: %s", includeCmd)
			}
		}

		commonArgs := []string{}

		// Directly invoke the RunE function of each command
		if err := CloneCmd.RunE(cmd, commonArgs); err != nil {
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
				if err := RemoveBackendCmd.RunE(cmd, []string{}); err != nil {
					return err
				}
			}
		}

		return nil
	},
}

func init() {
	OssCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().StringSliceVar(&includeCommands, "include", nil, "Specify additional commands to include in the migration process (e.g., --include remove-backend)")
}
