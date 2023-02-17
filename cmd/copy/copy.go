package copy

import (
	"github.com/spf13/cobra"
)

// `tfm copy` commands
var CopyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy command",
	Long:  "Copy objects from Source Organization to Destination Organization",
}

func init() {

}
