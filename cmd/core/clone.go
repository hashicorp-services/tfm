// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package core

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

	// `tfm core clone` command
	CloneCmd = &cobra.Command{
		Use:   "clone",
		Short: "Clone VCS repositories containing terraform code.",
		Long:  "clone VCS repositories containing terraform code. These will be iterated upon by tfm to download state files, read them, and push them to workspaces.",
		RunE: func(cmd *cobra.Command, args []string) error {

			return main(
				githubclient.CreateContext())
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			o.Close()
		},
	}
)

func init() {

	CloneCmd.Flags().SetInterspersed(false)

	// Add commands
	CoreCmd.AddCommand(CloneCmd)
}

func listRepos(ctx *githubclient.ClientContext) ([]*github.Repository, error) {
	reposList := viper.GetStringSlice("repos_to_clone")

	if ctx.GithubOrganization == "" || ctx.GithubToken == "" || ctx.GithubUsername == "" {
		return nil, fmt.Errorf("github_organization, github_username, or github_token not provided")
	}

	var allRepos []*github.Repository
	var filteredRepos []*github.Repository
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

	// If repos_to_clone is specified, filter the repositories
	if len(reposList) > 0 {

		o.AddFormattedMessageCalculated("Found %d repos to clone in `repos_to_clone` list.", len(reposList))

		repoMap := make(map[string]bool)
		for _, repoName := range reposList {
			repoMap[repoName] = true
		}

		for _, repo := range allRepos {
			if _, ok := repoMap[*repo.Name]; ok {
				filteredRepos = append(filteredRepos, repo)
			}
		}

		// If repos_to_clone list is empty then clone all repos in the github org
	} else {
		filteredRepos = allRepos

		o.AddFormattedMessageUserProvided("No repos_to_clone list found in config file. Getting All Repositories from Github organization: \n", ctx.GithubOrganization)
	}

	o.AddFormattedMessageCalculated("Found %d Repositories in GitHub org.\n", len(allRepos))

	return filteredRepos, nil
}

// cloneRepos clones the repositories returned by listRepos.
// It clones each repository into a subdirectory under github_cloned_repos_path.
func cloneRepos(ctx *githubclient.ClientContext, repos []*github.Repository) error {

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
		o.AddDeferredMessageRead("Cloned:", *repo.Name)
	}

	return nil
}

func main(ctx *githubclient.ClientContext) error {

	repos, err := listRepos(ctx)
	if err != nil {
		fmt.Printf("Failed to list repositories: %v\n", err)
		return nil
	}

	err = cloneRepos(ctx, repos)
	if err != nil {
		fmt.Printf("Failed to clone repositories: %v\n", err)
		return nil
	}

	return nil
}
