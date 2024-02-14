// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oss

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (

	// `tfm oss getstate` command
	getstateCmd = &cobra.Command{
		Use:   "getstate",
		Short: "Initialize and get state from terraform VCS repos.",
		Long:  "Initialize and get state from terraform VCS repos cloned by tfm.",
		RunE: func(cmd *cobra.Command, args []string) error {

			return initializeRepos()
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			o.Close()
		},
	}
)

func init() {

	getstateCmd.Flags().SetInterspersed(false)

	// Add commands
	OssCmd.AddCommand(getstateCmd)
}

func initializeRepos() error {
	clonePath := viper.GetString("github_clone_repos_path")

	// Initialize a counter to keep track of initialized repos
	var initCount int

	err := filepath.Walk(clonePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip the root directory quickly
		if path == clonePath {
			return nil
		}
		if info.IsDir() {
			// Check for .tf files to determine if it's a Terraform directory
			hasTfFiles, err := filepath.Glob(filepath.Join(path, "*.tf"))
			if err != nil {
				fmt.Printf("Error checking .tf files in %s: %v\n", path, err)
				return err
			}
			if len(hasTfFiles) > 0 {
				fmt.Printf("Initializing Terraform in: %s\n", path)
				cmd := exec.Command("terraform", "init")
				cmd.Dir = path
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					fmt.Printf("Failed to initialize Terraform in %s: %v\n", path, err)
				} else {
					initCount++
				}
				// Skip further files in this directory since we've already run `terraform init`
				return filepath.SkipDir
			}
		}
		// Continue walking into subdirectories
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking through directories: %v", err)
	}

	fmt.Printf("Terraform initialization completed for %d repositories.\n", initCount)
	return nil
}
