// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5/plumbing/transport/http"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/github"
	"github.com/hashicorp-services/tfm/output"
	vcsclients "github.com/hashicorp-services/tfm/vcsclients"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
)

var (
	o output.Output

	// `tfm core clone` command
	CloneCmd = &cobra.Command{
		Use:   "clone",
		Short: "Clone VCS repositories containing terraform code.",
		Long:  "clone VCS repositories containing terraform code. These will be iterated upon by tfm to download state files, read them, and push them to workspaces.",
		RunE: func(cmd *cobra.Command, args []string) error {

			vcsType := viper.GetString("vcs_type")

			switch vcsType {
			case "github":
				return main(vcsclients.CreateContext())
			case "gitlab":
				return main(vcsclients.CreateContextGitlab())
			default:
				return fmt.Errorf("unsupported VCS type: %s", vcsType)
			}
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

func listReposGithub(ctx *vcsclients.ClientContext) ([]*github.Repository, error) {
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

func cloneReposGithub(ctx *vcsclients.ClientContext, repos []*github.Repository) error {

	for _, repo := range repos {
		// Construct the directory path based on the repository name.
		clonePath := viper.GetString("clone_repos_path")

		if clonePath == "" {
			return fmt.Errorf("clone_repos_path is not configured")
		}

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

func listReposGitLab(ctx *vcsclients.GitLabClientContext) ([]*gitlab.Project, error) {
	reposList := viper.GetStringSlice("repos_to_clone")

	if ctx.GitLabGroup == "" || ctx.GitLabToken == "" || ctx.GitLabUsername == "" {
		return nil, fmt.Errorf("gitlab_group, gitlab_username, or gitlab_token not provided")
	}

	var allRepos []*gitlab.Project
	var filteredRepos []*gitlab.Project

	listOptions := gitlab.ListGroupProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 10,
			Page:    1,
		},
	}

	var allProjects []*gitlab.Project
	for {
		projects, response, err := ctx.GitLabClient.Groups.ListGroupProjects(ctx.GitLabGroup, &listOptions)
		if err != nil {
			return nil, err
		}
		allProjects = append(allProjects, projects...)

		if response.CurrentPage >= response.TotalPages {
			break
		}
		listOptions.Page = response.NextPage
	}

	// If repos_to_clone is specified, filter the repositories
	if len(reposList) > 0 {

		o.AddFormattedMessageCalculated("Found %d repos to clone in `repos_to_clone` list.", len(reposList))

		repoMap := make(map[string]bool)
		for _, repoName := range reposList {
			repoMap[repoName] = true
		}

		for _, project := range allProjects {
			if _, ok := repoMap[*&project.Name]; ok {
				filteredRepos = append(filteredRepos, project)
			}
		}

		// If repos_to_clone list is empty then clone all repos in the gitlab group
	} else {
		filteredRepos = allRepos

		o.AddFormattedMessageUserProvided("No repos_to_clone list found in config file. Getting All Repositories from GitLab Group: \n", ctx.GitLabGroup)
	}

	o.AddFormattedMessageCalculated("Found %d Repositories in GitLab group.\n", len(allProjects))

	return filteredRepos, nil
}

func cloneReposGitLab(ctx *vcsclients.GitLabClientContext, repos []*gitlab.Project) error {
	clonePath := viper.GetString("clone_repos_path")
	if clonePath == "" {
		return fmt.Errorf("clone_repos_path is not configured")
	}

	auth := &http.BasicAuth{
		Username: ctx.GitLabUsername,
		Password: ctx.GitLabToken,
	}

	for _, repo := range repos {
		dir := filepath.Join(clonePath, repo.Name)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			fmt.Printf("Cloning %s into %s\n", repo.WebURL, dir)
			_, err := git.PlainClone(dir, false, &git.CloneOptions{
				URL:      repo.HTTPURLToRepo,
				Progress: os.Stdout,
				Auth:     auth,
			})
			if err != nil {
				fmt.Printf("Error cloning %s: %v\n", repo.WebURL, err)
				continue
			}
		} else {
			fmt.Printf("Directory %s already exists, skipping clone of %s\n", dir, repo.WebURL)
		}
		o.AddDeferredMessageRead("Cloned:", *&repo.Name)

	}

	return nil
}

func main(ctx interface{}) error {
	vcsType := viper.GetString("vcs_type")

	var err error
	switch vcsType {
	case "github":
		githubCtx, ok := ctx.(*vcsclients.ClientContext)
		if !ok {
			return fmt.Errorf("invalid context for GitHub")
		}
		repos, err := listReposGithub(githubCtx)
		if err != nil {
			fmt.Printf("Failed to list GitHub repositories: %v\n", err)
			return err
		}
		err = cloneReposGithub(githubCtx, repos)
		if err != nil {
			fmt.Printf("Failed to cloneGitHub repos: %v\n", err)
			return err
		}

	case "gitlab":
		gitlabCtx, ok := ctx.(*vcsclients.GitLabClientContext)
		if !ok {
			return fmt.Errorf("invalid context for GitLab")
		}
		repos, err := listReposGitLab(gitlabCtx)
		if err != nil {
			fmt.Printf("Failed to list GitLab repositories: %v\n", err)
			return err
		}
		err = cloneReposGitLab(gitlabCtx, repos)
		if err != nil {
			fmt.Printf("Failed to clone GitLab repos: %v\n", err)
			return err
		}
	default:
		return fmt.Errorf("unsupported VCS type: %s", vcsType)
	}

	if err != nil {
		fmt.Printf("Failed to clone repositories: %v\n", err)
		return err
	}

	return nil
}
