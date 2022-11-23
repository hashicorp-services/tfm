package list

import (
	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
	"github.com/hashicorp-services/tfe-mig/tfclient"

)

var (

	// `tfe-migrate list teams` command
	teamsListCmd = &cobra.Command{
		Use:   "teams",
		Short: "Teams command",
		Long:  "Act upon Teams in an org",
		// RunE: func(cmd *cobra.Command, args []string) error {
		// 	return listTeams(
		// 		tfeclient.GetClientContexts())

		// },
		Run: func(cmd *cobra.Command, args []string) {
			// return orgShow(
			// 	viper.GetString("name"))
			listTeams(tfclient.GetClientContexts())
		},
	}
)

func init() {
	// Flags().StringP, etc... - the "P" gives us the option for a short hand

	// `tfe-migrate copy teams all` command
	teamsListCmd.Flags().BoolP("all", "a", false, "List all? (optional)")

	// Add commands
	ListCmd.AddCommand(teamsListCmd)

}

func listTeams(c tfclient.ClientContexts) error {
	o.AddMessageUserProvided("Getting list of teams from: ", c.SourceHostname)
	srcTeams := []*tfe.Team{}

	opts := tfe.TeamListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100},
	}
	for {
		items, err := c.SourceClient.Teams.List(c.SourceContext, c.SourceOrganizationName, &opts)
		if err != nil {
			return err
		}

		srcTeams = append(srcTeams, items.Items...)

		o.AddFormattedMessageCalculated("Found %d Teams", len(srcTeams))

		if items.CurrentPage >= items.TotalPages {
			break
		}
		opts.PageNumber = items.NextPage

	}
	o.AddTableHeaders("Name")
	for _, i := range srcTeams {

		o.AddTableRows(i.Name)
	}

	return nil
}
