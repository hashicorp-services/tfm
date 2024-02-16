// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	githubclient "github.com/hashicorp-services/tfm/vcsclients"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RemoveBackendCmd = &cobra.Command{
	Use:   "remove-backend",
	Short: "Create a branch, remove Terraform backend configurations from cloned repos in github_clone_repos_path, commit the changes, and push to the origin.",
	Long:  `Searches through .tf files in the root of cloned repositories to remove backend configurations and commit them back on a new branch.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		clonePath := viper.GetString("github_clone_repos_path")
		branchName := "update-backend-" + time.Now().Format("20060102")

		// Step 1: Create branches as needed
		reposWithNewBranches, err := createBranchIfNeeded(clonePath, branchName)
		if err != nil {
			return err
		}

		// Step 2: Remove backend configurations
		err = removeBackendFromRepos(clonePath)
		if err != nil {
			return err
		}

		// Step 3: Commit changes in repos that had new branches created
		// Only proceed to commit if there are repos with new branches
		if len(reposWithNewBranches) > 0 {
			err = commitChangesInRepos(reposWithNewBranches, branchName, "Remove backend configuration")
			err = pushBranches(githubclient.CreateContext(), reposWithNewBranches, branchName, "origin")
			if err != nil {
				fmt.Println("Error committing or pushing branches:", err)
			}

		} else {
			fmt.Println("No new branches were created, skipping commit step.")
		}

		return nil
	},
}

func init() {
	CoreCmd.AddCommand(RemoveBackendCmd)
}

// detectBackendBlocks checks if there's a backend block in any .tf file within the repo.
func detectBackendBlocks(repoPath string) (bool, error) {
	files, err := ioutil.ReadDir(repoPath)
	if err != nil {
		return false, err
	}

	backendRegexp := regexp.MustCompile(`(?s)backend\s+"[^"]+"\s+\{.*?\}`)
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".tf" {
			continue
		}

		content, err := ioutil.ReadFile(filepath.Join(repoPath, file.Name()))
		if err != nil {
			return false, err
		}

		if backendRegexp.Match(content) {
			return true, nil
		}
	}

	return false, nil
}

func createBranchIfNeeded(clonePath, branchName string) ([]string, error) {
	var reposWithNewBranches []string

	dirs, err := os.ReadDir(clonePath)
	if err != nil {
		return nil, fmt.Errorf("error reading clone path directories: %v", err)
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		repoPath := filepath.Join(clonePath, dir.Name())
		hasBackend, err := detectBackendBlocks(repoPath)
		if err != nil {
			fmt.Printf("Error detecting backend blocks in %s: %v\n", repoPath, err)
			continue
		}

		if hasBackend {
			repo, err := git.PlainOpen(repoPath)
			if err != nil {
				fmt.Printf("Failed to open repo: %v\n", err)
				continue
			}

			// Getting the HEAD reference to find the current commit hash
			headRef, err := repo.Head()
			if err != nil {
				fmt.Printf("Failed to get HEAD reference: %v\n", err)
				continue
			}

			// Creating a new branch reference
			branchRefName := plumbing.NewBranchReferenceName(branchName)
			ref := plumbing.NewHashReference(branchRefName, headRef.Hash())

			// Check if the branch already exists
			_, err = repo.Reference(branchRefName, false)
			if err == nil {
				fmt.Printf("Branch '%s' already exists in %s\n", branchName, repoPath)
				continue
			} else if err != plumbing.ErrReferenceNotFound {
				fmt.Printf("Error checking for branch existence: %v\n", err)
				continue
			}

			// Creating the branch with Config
			err = repo.Storer.SetReference(ref)
			if err != nil {
				fmt.Printf("Failed to create branch '%s': %v\n", branchName, err)
				continue
			}

			// Checkout to the newly created branch
			w, err := repo.Worktree()
			if err != nil {
				fmt.Printf("Failed to get worktree: %v\n", err)
				continue
			}

			err = w.Checkout(&git.CheckoutOptions{
				Branch: branchRefName,
				Create: false, // Not creating a new branch here, just switching to it
			})
			if err != nil {
				fmt.Printf("Failed to checkout branch '%s': %v\n", branchName, err)
				continue
			}

			fmt.Printf("Branch '%s' created and checked out in %s\n", branchName, repoPath)
			reposWithNewBranches = append(reposWithNewBranches, repoPath)
		} else {
			fmt.Printf("No backend block found in %s, skipping branch creation\n", repoPath)
		}
	}

	return reposWithNewBranches, nil
}

func removeBackendFromRepos(clonePath string) error {
	dirs, err := os.ReadDir(clonePath)
	if err != nil {
		return fmt.Errorf("error reading clone path directories: %v", err)
	}

	backendRegexp := regexp.MustCompile(`(?s)backend\s+"[^"]+"\s+\{.*?\}`)

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		repoPath := filepath.Join(clonePath, dir.Name())
		files, err := ioutil.ReadDir(repoPath)
		if err != nil {
			fmt.Printf("Error reading repo directory: %v\n", err)
			continue
		}

		repoModified := false

		for _, file := range files {
			if filepath.Ext(file.Name()) != ".tf" {
				continue
			}

			filePath := filepath.Join(repoPath, file.Name())
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				fmt.Printf("Error reading .tf file: %v\n", err)
				continue
			}

			modifiedContent := backendRegexp.ReplaceAll(content, []byte(""))

			if len(modifiedContent) != len(content) {
				err = ioutil.WriteFile(filePath, modifiedContent, file.Mode())
				if err != nil {
					fmt.Printf("Error writing modified .tf file: %v\n", err)
					continue
				}
				fmt.Printf("Removed backend configuration from: %s\n", filePath)
				repoModified = true
			}
		}

		if !repoModified {
			fmt.Printf("No backend blocks found in: %s\n", repoPath)
		}
	}

	return nil
}

func commitChanges(repoPath, branchName string) error {
	commitMessage := viper.GetString("commit_message")
	if commitMessage == "" {
		commitMessage = "Removed backend configuration for migration to TFC/TFE."
	}

	authorName := viper.GetString("commit_author_name")
	authorEmail := viper.GetString("commit_author_email")

	// Open the existing repository
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return fmt.Errorf("failed to open repo: %v", err)
	}

	// Get the worktree to stage changes
	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %v", err)
	}

	// Add changes to the staging area
	// Using "." to add all changes in the repository
	_, err = w.Add(".")
	if err != nil {
		return fmt.Errorf("failed to add changes to staging area: %v", err)
	}

	// Commit the changes
	_, err = w.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  authorName,
			Email: authorEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to commit changes: %v", err)
	}

	return nil
}

func commitChangesInRepos(reposWithNewBranches []string, branchName, commitMessage string) error {
	for _, repoPath := range reposWithNewBranches {
		fmt.Printf("Committing changes in %s on branch '%s'\n", repoPath, branchName)
		err := commitChanges(repoPath, branchName)
		if err != nil {
			fmt.Printf("Error committing changes in %s: %v\n", repoPath, err)
		}
	}
	return nil
}

// pushBranches takes a list of repository paths where commits were made and pushes the specified branch to the remote
func pushBranches(ctx *githubclient.ClientContext, reposWithNewBranches []string, branchName, remoteName string) error {
	for _, repoPath := range reposWithNewBranches {

		repo, err := git.PlainOpen(repoPath)
		if err != nil {
			fmt.Printf("Failed to open repo at %s: %v\n", repoPath, err)
			continue
		}

		// Get the branch reference
		refName := plumbing.NewBranchReferenceName(branchName)

		// Push the changes using GitHub token for authentication
		err = repo.Push(&git.PushOptions{
			RemoteName: remoteName,
			RefSpecs: []config.RefSpec{
				config.RefSpec(refName + ":" + refName),
			},
			Auth: &http.BasicAuth{
				Username: ctx.GithubUsername,
				Password: ctx.GithubToken,
			},
		})
		if err != nil {
			if err == git.NoErrAlreadyUpToDate {
				fmt.Printf("Branch '%s' in repo at %s is already up-to-date with remote '%s'\n", branchName, repoPath, remoteName)
			} else {
				fmt.Printf("Failed to push branch '%s' in repo at %s: %v\n", branchName, repoPath, err)
			}
			continue
		}

		fmt.Printf("Branch '%s' in repo at %s pushed successfully to remote '%s'\n", branchName, repoPath, remoteName)
	}

	return nil
}
