package delete

import (
	"github.com/spf13/cobra"
)

var (
	side    string

	DeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "delete command",
		Long:  "delete objects in an org. DANGER this will delete things!",
	}
)

func init() {

	DeleteCmd.PersistentFlags().StringVar(&side, "side", "", "Specify source or destination side to process")

}
