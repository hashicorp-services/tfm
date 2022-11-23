package copy

import (
	"github.com/spf13/cobra"
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type ClientContexts struct {
	SourceClient                *tfe.Client
	SourceContext               context.Context
	SourceHostname              string
	SourceOrganizationName      string
	SourceToken                 string
	DestinationClient           *tfe.Client
	DestinationContext          context.Context
	DestinationHostname         string
	DestinationOrganizationName string
	DestinationToken            string
}


// `tfe-migrate copy` commands
var CopyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy command",
	Long:  "Copy objects from Source Organization to Destination Organization",
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