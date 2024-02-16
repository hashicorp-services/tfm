// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var autoApprove bool

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Cleans up all repositories in the clone path",
	Long:  `Deletes all repositories that were cloned into the specified clone path, cleaning up the workspace.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		clonePath := viper.GetString("github_clone_repos_path")
		if clonePath == "" {
			return fmt.Errorf("clone path is not specified in the configuration")
		}

		if !autoApprove {
			fmt.Printf("Are you sure you want to delete all repositories in %s? [yes/no]: ", clonePath)
			var response string
			fmt.Scanln(&response)
			if response != "yes" {
				fmt.Println("Cleanup aborted.")
				return nil
			}
		}

		return cleanupRepos(clonePath)
	},
}

func init() {
	CoreCmd.AddCommand(cleanupCmd)
	cleanupCmd.Flags().BoolVar(&autoApprove, "autoapprove", false, "Automatically approve the operation without a confirmation prompt")
}

func cleanupRepos(clonePath string) error {
	dirs, err := os.ReadDir(clonePath)
	if err != nil {
		return fmt.Errorf("error reading clone path directories: %v", err)
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			dirPath := filepath.Join(clonePath, dir.Name())
			fmt.Printf("Removing directory: %s\n", dirPath)
			err := os.RemoveAll(dirPath)
			if err != nil {
				return fmt.Errorf("failed to remove directory %s: %v", dirPath, err)
			}
		}
	}

	fmt.Println("Cleanup completed successfully.")
	return nil
}
