package list

import (
	"github.com/hashicorp-services/tfm/tfclient"
	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

var (

	// `tfm list project` command
	projectListCmd = &cobra.Command{
		Use:     "projects",
		Aliases: []string{"prj"},
		Short:   "Projects command",
		Long:    "List Projects in an org",
		Run: func(cmd *cobra.Command, args []string) {
			listProjects(tfclient.GetClientContexts())
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			o.Close()
		},
	}
)

func init() {

	// Add commands
	ListCmd.AddCommand(projectListCmd)

}

func listProjects(c tfclient.ClientContexts) error {

	srcProjects := []*tfe.Project{}

	opts := tfe.ProjectListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100},
	}

	if (ListCmd.Flags().Lookup("side").Value.String() == "source") || (!ListCmd.Flags().Lookup("side").Changed) {

		o.AddMessageUserProvided("Getting list of projects from: ", c.SourceHostname)

		for {
			items, err := c.SourceClient.Projects.List(c.SourceContext, c.SourceOrganizationName, &opts)
			if err != nil {
				return err
			}

			srcProjects = append(srcProjects, items.Items...)

			o.AddFormattedMessageCalculated("Found %d Projects", len(srcProjects))

			if items.CurrentPage >= items.TotalPages {
				break
			}
			opts.PageNumber = items.NextPage

		}
		o.AddTableHeaders("Name", "ID")
		for _, i := range srcProjects {

			o.AddTableRows(i.Name, i.ID)
		}
	}

	if ListCmd.Flags().Lookup("side").Value.String() == "destination" {
		o.AddMessageUserProvided("Getting list of projects from: ", c.DestinationHostname)

		for {
			items, err := c.DestinationClient.Projects.List(c.DestinationContext, c.DestinationOrganizationName, &opts)
			if err != nil {
				return err
			}

			srcProjects = append(srcProjects, items.Items...)

			o.AddFormattedMessageCalculated("Found %d Projects", len(srcProjects))

			if items.CurrentPage >= items.TotalPages {
				break
			}
			opts.PageNumber = items.NextPage

		}
		o.AddTableHeaders("Name", "ID")
		for _, i := range srcProjects {

			o.AddTableRows(i.Name, i.ID)
		}
	}

	return nil
}
