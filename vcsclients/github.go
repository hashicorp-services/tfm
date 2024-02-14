package githubclient

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func NewGitHubClient(ctx context.Context) *github.Client {
	// Retrieve the GitHub token from viper
	token := viper.GetString("github_token") // Ensure your HCL structure is reflected here

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
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
	githubClient := NewGitHubClient(ctx)

	return &ClientContext{
		GitHubClient:       githubClient,
		GithubContext:      ctx,
		GithubToken:        viper.GetString("github_token"),
		GithubOrganization: viper.GetString("github_organization"),
		GithubUsername:     viper.GetString("github_username"),
	}
}
