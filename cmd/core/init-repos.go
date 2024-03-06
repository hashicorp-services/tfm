package core

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Define the command
var InitReposCmd = &cobra.Command{
	Use:   "init-repos",
	Short: "Scan cloned repositories for Terraform configurations and build metadata",
	Long: `Scans all cloned repositories based on the 'clone_repos_path' from the configuration file,
identifies directories containing Terraform configurations, and builds a metadata file summarizing these findings.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return initRepos()
	},
}

func init() {
	CoreCmd.AddCommand(InitReposCmd)
}

type WorkspaceInfo struct {
	UsesWorkspaces bool     `json:"uses_workspaces"`
	WorkspaceNames []string `json:"workspace_names,omitempty"`
}

type ConfigPathInfo struct {
	Path          string        `json:"path"`
	WorkspaceInfo WorkspaceInfo `json:"workspace_info"`
}

type RepoConfig struct {
	RepoName    string           `json:"repo_name"`
	ConfigPaths []ConfigPathInfo `json:"config_paths"`
}

func initRepos() error {
	clonedReposPath := viper.GetString("clone_repos_path")
	if clonedReposPath == "" {
		return fmt.Errorf("clone_repos_path is not configured")
	}

	var repoConfigs []RepoConfig

	err := filepath.WalkDir(clonedReposPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".tf" {
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
	fullConfigDir := filepath.Join(viper.GetString("clone_repos_path"), configPath)

	workspaceInfo := checkTerraformWorkspaces(fullConfigDir)

	for i, rc := range *repoConfigs {
		if rc.RepoName == repoName {
			(*repoConfigs)[i].ConfigPaths = append(rc.ConfigPaths, ConfigPathInfo{
				Path:          configPath,
				WorkspaceInfo: workspaceInfo,
			})
			return
		}
	}
	*repoConfigs = append(*repoConfigs, RepoConfig{
		RepoName:    repoName,
		ConfigPaths: []ConfigPathInfo{{Path: configPath, WorkspaceInfo: workspaceInfo}},
	})
}

func parseWorkspaceListOutput(output string) []string {
	var workspaces []string
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = strings.TrimPrefix(line, "* ")
		if line != "" {
			workspaces = append(workspaces, line)
		}
	}
	return workspaces
}

func runTerraformInit2(configDir string) error {
	initCmd := exec.Command("terraform", "init", "-input=false", "-no-color")
	initCmd.Dir = configDir
	initCmd.Stdout = os.Stdout
	initCmd.Stderr = os.Stderr
	return initCmd.Run()
}

func listTerraformWorkspaces(configDir string) ([]string, error) {
	cmd := exec.Command("terraform", "workspace", "list", "-no-color")
	cmd.Dir = configDir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("error listing Terraform workspaces: %w", err)
	}

	// Parse output to extract workspace names
	workspaces := parseWorkspaceListOutput(out.String())
	return workspaces, nil
}

func checkTerraformWorkspaces(configDir string) WorkspaceInfo {
	// run terraform init
	if err := runTerraformInit(configDir); err != nil {
		fmt.Printf("Error initializing Terraform: %s\n", err)
		return WorkspaceInfo{} // Consider how you want to handle init errors
	}

	// list workspaces
	workspaces, err := listTerraformWorkspaces(configDir)
	if err != nil {
		fmt.Printf("Error listing Terraform workspaces: %s\n", err)
		return WorkspaceInfo{}
	}

	// Determine if using workspaces beyond 'default'
	usesWorkspaces := len(workspaces) > 1 || (len(workspaces) == 1 && workspaces[0] != "default")

	return WorkspaceInfo{
		UsesWorkspaces: usesWorkspaces,
		WorkspaceNames: workspaces,
	}
}

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
