// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oss

import (
	"fmt"
	"io/ioutil"
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

func runTerraformInit(dirPath string) error {
	cmd := exec.Command("terraform", "init")
	cmd.Dir = dirPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func pullTerraformState(dirPath, outputPath string) error {
	cmd := exec.Command("terraform", "state", "pull")
	cmd.Dir = dirPath
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(outputPath, output, 0644)
}

func initializeRepos() error {
	clonePath := viper.GetString("github_clone_repos_path")

	var initCount int

	// Read directories directly under clonePath
	dirs, err := os.ReadDir(clonePath) // Using os.ReadDir, which is the recommended way since Go 1.16
	if err != nil {
		return fmt.Errorf("error reading directories: %v", err)
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue // Skip files, process only directories
		}
		repoPath := filepath.Join(clonePath, dir.Name())

		// Check for .tf files directly in the root of the repository
		hasTfFiles, err := filepath.Glob(filepath.Join(repoPath, "*.tf"))
		if err != nil {
			fmt.Printf("Error checking .tf files in %s: %v\n", repoPath, err)
			continue // Proceed to next directory on error
		}
		if len(hasTfFiles) > 0 {
			fmt.Printf("Initializing Terraform in: %s\n", repoPath)
			if err := runTerraformInit(repoPath); err != nil {
				fmt.Printf("Failed to initialize Terraform in %s: %v\n", repoPath, err)
				continue // Proceed to next directory on error
			}
			initCount++

			// Pull the state and save it to pulled_terraform.tfstate
			pulledStatePath := filepath.Join(repoPath, ".terraform/pulled_terraform.tfstate")
			if err := pullTerraformState(repoPath, pulledStatePath); err != nil {
				fmt.Printf("Failed to pull Terraform state in %s: %v\n", repoPath, err)
				continue // Proceed to next directory despite the error
			}

			// Here you can add additional processing for the pulled state file as needed
		}
	}

	o.AddFormattedMessageCalculated("Terraform initialization and state processing completed for %d repositories.\n", initCount)
	return nil
}
