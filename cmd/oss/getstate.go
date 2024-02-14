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

	err := filepath.Walk(clonePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == clonePath {
			return nil
		}
		if info.IsDir() {
			hasTfFiles, err := filepath.Glob(filepath.Join(path, "*.tf"))
			if err != nil {
				fmt.Printf("Error checking .tf files in %s: %v\n", path, err)
				return err
			}
			if len(hasTfFiles) > 0 {
				fmt.Printf("Initializing Terraform in: %s\n", path)
				if err := runTerraformInit(path); err != nil {
					fmt.Printf("Failed to initialize Terraform in %s: %v\n", path, err)
					return nil
				}
				initCount++

				// Pull the state and save it to pulled_terraform.tfstate
				pulledStatePath := filepath.Join(path, ".terraform/pulled_terraform.tfstate")
				if err := pullTerraformState(path, pulledStatePath); err != nil {
					fmt.Printf("Failed to pull Terraform state in %s: %v\n", path, err)
					return nil
				}
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking through directories: %v", err)
	}

	fmt.Printf("Terraform initialization and state processing completed for %d repositories.\n", initCount)
	return nil
}
