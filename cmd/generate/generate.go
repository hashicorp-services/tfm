// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package generate

import (
	"github.com/spf13/cobra"
)

var (
	GenerateCmd = &cobra.Command{
		Use:   "generate",
		Short: "generate command for generating .tfm.hcl config template",
		Long:  "generate a .tfm.hcl file template ",
	}
)

func init() {

}
