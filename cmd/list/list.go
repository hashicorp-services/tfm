package list

import (
	"github.com/spf13/cobra"
)

// `tfm list` commands
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List command",
	Long:  "List objects in an org",
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// discoverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// discoverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
