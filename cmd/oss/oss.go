// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oss

import (
	"github.com/spf13/cobra"
)

// `tfm oss` commands
var OssCmd = &cobra.Command{
	Use:   "oss",
	Short: "oss command",
	Long:  "Command used to perform terraform open source (OSS) to TFE/TFC migration commands",
}

func init() {

}

// Tfm oss clone-repos
// Tfm oss init-repos
// Tfm oss create-workspaces
// Tfm oss connect-vcs
// Tfm lock workspaces
// Tfm oss push-state
// Tfm unlock workspaces
// Tfm oss remove-backend

// github {
// 	github_token = "your_github_token_here"
// 	github_organization = "optional_default_organization"
//   }
