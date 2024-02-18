package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Define the command
var initReposCmd = &cobra.Command{
	Use:   "init-repos",
	Short: "Scan cloned repositories for Terraform configurations and build metadata",
	Long: `Scans all cloned repositories based on the 'github_cloned_repos_path' from the configuration file,
identifies directories containing Terraform configurations, and builds a metadata file summarizing these findings.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return initRepos()
	},
}

func init() {
	// Assuming rootCmd is your application's root Cobra command
	CoreCmd.AddCommand(initReposCmd)
}

// RepoConfig represents the metadata for a repository's Terraform configurations
type RepoConfig struct {
	RepoName    string   `json:"repo_name"`
	ConfigPaths []string `json:"config_paths"`
}

// initRepos scans repositories for Terraform configurations containing a backend block and generates metadata
func initRepos() error {
	clonedReposPath := viper.GetString("github_clone_repos_path")

	if clonedReposPath == "" {
		return fmt.Errorf("github_clone_repos_path is not configured")
	}

	var repoConfigs []RepoConfig

	err := filepath.WalkDir(clonedReposPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".tf" {
			// Check if the file contains a backend block within a terraform block
			containsBackend, err := fileContainsBackendBlock(path)
			if err != nil {
				return err
			}
			if containsBackend {
				relPath, err := filepath.Rel(clonedReposPath, filepath.Dir(path))
				if err != nil {
					return err
				}
				repoName := strings.Split(relPath, string(os.PathSeparator))[0]
				addRepoConfig(&repoConfigs, repoName, relPath)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error scanning repositories: %w", err)
	}

	return saveMetadata(repoConfigs)
}

func fileContainsBackendBlock(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var terraformBlockStarted bool
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "terraform {") {
			terraformBlockStarted = true
		}
		if terraformBlockStarted && strings.Contains(line, "backend") && strings.Contains(line, "{") {
			return true, nil
		}
		if terraformBlockStarted && strings.Contains(line, "}") {
			terraformBlockStarted = false
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}

func addRepoConfig(repoConfigs *[]RepoConfig, repoName, configPath string) {
	for i, rc := range *repoConfigs {
		if rc.RepoName == repoName {
			if !contains(rc.ConfigPaths, configPath) {
				(*repoConfigs)[i].ConfigPaths = append(rc.ConfigPaths, configPath)
			}
			return
		}
	}
	*repoConfigs = append(*repoConfigs, RepoConfig{
		RepoName:    repoName,
		ConfigPaths: []string{configPath},
	})
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// saveMetadata saves the repositories' Terraform configuration metadata to a file
func saveMetadata(repoConfigs []RepoConfig) error {
	metadataFile, err := os.Create("terraform_config_metadata.json")
	if err != nil {
		return fmt.Errorf("error creating metadata file: %w", err)
	}
	defer metadataFile.Close()

	encoder := json.NewEncoder(metadataFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(repoConfigs); err != nil {
		return fmt.Errorf("error writing metadata to file: %w", err)
	}

	fmt.Println("Metadata file 'terraform_config_metadata.json' created successfully.")
	return nil
}
