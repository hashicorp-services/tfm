// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lock

import (
	"github.com/spf13/cobra"
)

// `tfm lock` commands
var (
	side    string

	LockCmd = &cobra.Command{
		Use:   "lock",
		Short: "Lock",
		Long:  "Locks objects in source or destination",
	}
)

func init() {
	LockCmd.PersistentFlags().StringVar(&side, "side", "", "Specify source or destination side to process")
}
