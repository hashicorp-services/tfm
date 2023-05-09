// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package list

import (
	"github.com/spf13/cobra"
)

var (
	side    string
	
	ListCmd = &cobra.Command{
		Use:   "list",
		Short: "List command",
		Long:  "List objects in an org",
	}
)

func init() {

	ListCmd.PersistentFlags().StringVar(&side, "side", "", "Specify source or destination side to process")
}
