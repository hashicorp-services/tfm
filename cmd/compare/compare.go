package compare

import (
	"github.com/spf13/cobra"
)

var (
	jsonOut bool

	CmpCmd = &cobra.Command{
		Use:     "compare",
		Aliases: []string{"cmp"},
		Short:   "Compare command",
		Long:    "Compare command",
	}
)

func init() {
	CmpCmd.PersistentFlags().BoolVar(&jsonOut, "json", false, "Print the output in JSON format. Only supported with [workspaces, projects]")
}
