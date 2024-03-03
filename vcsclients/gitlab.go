package main

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/spf13/viper"
	gitlab "github.com/xanzy/go-gitlab"
)

// NewGitLabClient creates a new GitLab client
func NewGitLabClient(ctx context.Context) *gitlab.Client {

	token := viper.GetString("gitlab_token")

	client, err := gitlab.NewClient(token)
	if err != nil {
		panic(err)
	}
	return client
}

type ClientContext struct {
	GitHubClient       *github.Client
	GithubContext      context.Context
	GithubToken        string
	GithubOrganization string
	GithubUsername     string
}

func CreateContext() *ClientContext {
	ctx := context.Background()
	gitlabClient := NewGitLabClient(ctx)

	return &ClientContext{
		GitHubClient:       gitlabClient,
		GithubContext:      ctx,
		GithubToken:        viper.GetString("github_token"),
		GithubOrganization: viper.GetString("github_organization"),
		GithubUsername:     viper.GetString("github_username"),
	}
}
