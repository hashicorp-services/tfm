// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package unlock

import (
	"github.com/spf13/cobra"
)

// `tfm lock` commands
var (
	side    string

	UnlockCmd = &cobra.Command{
		Use:   "unlock",
		Short: "Unlock",
		Long:  "Unlocks objects in source or destination",
	}
)

func init() {
	UnlockCmd.PersistentFlags().StringVar(&side, "side", "", "Specify source or destination side to process")
}
