// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oss

import (
	"github.com/spf13/cobra"
)

// `tfm oss` commands
var OssCmd = &cobra.Command{
	Use:   "oss",
	Short: "Command used to perform terraform open source (OSS) to TFE/TFC migration commands",
	Long:  "Command used to perform terraform open source (OSS) to TFE/TFC migration commands",
}

func init() {

}
