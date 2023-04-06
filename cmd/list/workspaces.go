package list

import (
	"github.com/hashicorp-services/tfm/tfclient"
	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

var (

	// `tfm list teams` command
	workspacesListCmd = &cobra.Command{
		Use:   "workspaces",
		Short: "workspaces command",
		Long:  "Act upon workspaces in an org",
		// RunE: func(cmd *cobra.Command, args []string) error {
		// 	return listworkspaces(
		// 		tfeclient.GetClientContexts())

		// },
		Run: func(cmd *cobra.Command, args []string) {
			// return orgShow(
			// 	viper.GetString("name"))
			listworkspaces(tfclient.GetClientContexts())
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			o.Close()
		},
	}
)

func init() {
	// Flags().StringP, etc... - the "P" gives us the option for a short hand

	// `tfm copy workspaces all` command
	// workspacesListCmd.Flags().BoolP("all", "a", false, "List all? (optional)")

	// Add commands
	ListCmd.AddCommand(workspacesListCmd)

}

func listworkspaces(c tfclient.ClientContexts) error {

	srcWorkspaces := []*tfe.WorkspaceList{}

	opts := tfe.WorkspaceListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100},
	}

	if (ListCmd.Flags().Lookup("side").Value.String() == "source") || (!ListCmd.Flags().Lookup("side").Changed) {

		o.AddMessageUserProvided("Getting list of workspaces from: ", c.SourceHostname)

		for {
			items, err := c.SourceClient.Workspaces.List(c.SourceContext, c.SourceOrganizationName, &opts)
			if err != nil {
				return err
			}

			srcWorkspaces = append(srcWorkspaces, items.Items...)

			o.AddFormattedMessageCalculated("Found %d Workspaces", len(srcWorkspaces))

			if items.CurrentPage >= items.TotalPages {
				break
			}
			opts.PageNumber = items.NextPage

		}
		o.AddTableHeaders("Name")
		for _, i := range srcWorkspaces {

			o.AddTableRows(i.Name)
		}
	}

	if ListCmd.Flags().Lookup("side").Value.String() == "destination" {
		o.AddMessageUserProvided("Getting list of teams from: ", c.DestinationHostname)

		for {
			items, err := c.DestinationClient.Workspaces.List(c.DestinationContext, c.DestinationOrganizationName, &opts)
			if err != nil {
				return err
			}

			srcWorkspaces = append(srcWorkspaces, items.Items...)

			o.AddFormattedMessageCalculated("Found %d Workspaces", len(srcWorkspaces))

			if items.CurrentPage >= items.TotalPages {
				break
			}
			opts.PageNumber = items.NextPage

		}
		o.AddTableHeaders("Name")
		for _, i := range srcWorkspaces {

			o.AddTableRows(i.Name)
		}
	}

	return nil
}
