package nuke

import (
	"github.com/spf13/cobra"
)

var (
	NukeCmd = &cobra.Command{
		Use:   "nuke",
		Short: "nuke command",
		Long:  "nuke objects in an org. DANGER this will delete things!",
		Hidden: true,
	}
)

func init() {}
