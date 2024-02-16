// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package core

import (
	"github.com/spf13/cobra"
)

// `tfm core` commands
var CoreCmd = &cobra.Command{
	Use:   "core",
	Short: "Command used to perform terraform open source (core) to TFE/TFC migration commands",
	Long:  "Command used to perform terraform open source (core) to TFE/TFC migration commands",
}

func init() {

}
