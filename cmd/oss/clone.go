// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oss

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/github"
	"github.com/hashicorp-services/tfm/output"
	githubclient "github.com/hashicorp-services/tfm/vcsclients"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	o output.Output

	// `tfm oss clone` command
	cloneCmd = &cobra.Command{
		Use:   "clone",
		Short: "Clone VCS repositories containing terraform code.",
		Long:  "clone VCS repositories containing terraform code. These will be iterated upon by tfm to download state files, read them, and push them to workspaces.",
		RunE: func(cmd *cobra.Command, args []string) error {

			return cloneRepos(
				githubclient.CreateContext())
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			o.Close()
		},
	}
)

func init() {

	cloneCmd.Flags().SetInterspersed(false)

	// Add commands
	OssCmd.AddCommand(cloneCmd)
}

// listRepos lists all repositories for the configured organization or user.
func listRepos(ctx *githubclient.ClientContext) ([]*github.Repository, error) {

	o.AddFormattedMessageUserProvided("Getting list of Repositories from Github organization: \n", ctx.GithubOrganization)

	var allRepos []*github.Repository
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	for {
		repos, resp, err := ctx.GitHubClient.Repositories.ListByOrg(ctx.GithubContext, ctx.GithubOrganization, opt)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	o.AddFormattedMessageCalculated("Found %d Repositories\n", len(allRepos))

	return allRepos, nil
}

// cloneRepos clones the repositories returned by listRepos.
// It clones each repository into a subdirectory under the current working directory.
func cloneRepos(ctx *githubclient.ClientContext) error {
	repos, err := listRepos(ctx)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		// Construct the directory path based on the repository name.
		clonePath := viper.GetString("github_clone_repos_path")

		dir := filepath.Join(clonePath, *repo.Name)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			fmt.Printf("Cloning %s into %s\n", *repo.FullName, dir)
			_, err := git.PlainClone(dir, false, &git.CloneOptions{
				URL:      *repo.CloneURL,
				Progress: os.Stdout,
			})
			if err != nil {
				fmt.Printf("Error cloning %s: %v\n", *repo.FullName, err)

			}
		} else {
			fmt.Printf("Directory %s already exists, skipping clone of %s\n", dir, *repo.FullName)
		}
		o.AddDeferredMessageRead("Cloned\n", *repo.Name)
	}

	return nil
}
