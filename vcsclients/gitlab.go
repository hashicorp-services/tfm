package vcsclients

import (
	"context"

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

type GitLabClientContext struct {
	GitLabClient   *gitlab.Client
	GitLabContext  context.Context
	GitLabToken    string
	GitLabGroup    string
	GitLabUsername string
}

func CreateContextGitlab() *GitLabClientContext {
	ctx := context.Background()
	gitlabClient := NewGitLabClient(ctx)

	return &GitLabClientContext{
		GitLabClient:   gitlabClient,
		GitLabContext:  ctx,
		GitLabToken:    viper.GetString("gitlab_token"),
		GitLabGroup:    viper.GetString("gitlab_group"),
		GitLabUsername: viper.GetString("gitlab_username"),
	}
}
