// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package list

import (
	"github.com/hashicorp-services/tfm/cmd/logging"
	"github.com/hashicorp-services/tfm/tfclient"
	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

var (

	// `tfm list teams` command
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
		PostRun: func(cmd *cobra.Command, args []string) {
			o.Close()
		},
	}
)

func init() {
	// Flags().StringP, etc... - the "P" gives us the option for a short hand

	// `tfm copy teams all` command
	//teamsListCmd.Flags().BoolP("all", "a", false, "List all? (optional)")

	// Add commands
	ListCmd.AddCommand(teamsListCmd)

}

func listTeams(c tfclient.ClientContexts) error {
	log := logging.NewLogger("list.teams")

	srcTeams := []*tfe.Team{}

	opts := tfe.TeamListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100},
	}

	if (ListCmd.Flags().Lookup("side").Value.String() == "source") || (!ListCmd.Flags().Lookup("side").Changed) {

		log.Info("listing teams", "org", c.SourceOrganizationName, "host", c.SourceHostname)
		o.AddMessageUserProvided("Getting list of teams from: ", c.SourceHostname)

		for {
			items, err := c.SourceClient.Teams.List(c.SourceContext, c.SourceOrganizationName, &opts)
			if err != nil {
				log.Error("failed to list teams", "org", c.SourceOrganizationName, "error", err)
				return err
			}

			srcTeams = append(srcTeams, items.Items...)
			log.Debug("fetched team page", "page", items.CurrentPage, "total_pages", items.TotalPages, "count", len(srcTeams))

			o.AddFormattedMessageCalculated("Found %d Teams", len(srcTeams))

			if items.CurrentPage >= items.TotalPages {
				break
			}
			opts.PageNumber = items.NextPage

		}
		o.AddTableHeaders("Name")
		log.Info("found teams", "org", c.SourceOrganizationName, "count", len(srcTeams))
		for _, i := range srcTeams {

			o.AddTableRows(i.Name)
		}
	}

	if ListCmd.Flags().Lookup("side").Value.String() == "destination" {
		log.Info("listing teams", "org", c.DestinationOrganizationName, "host", c.DestinationHostname)
		o.AddMessageUserProvided("Getting list of teams from: ", c.DestinationHostname)

		for {
			items, err := c.DestinationClient.Teams.List(c.DestinationContext, c.DestinationOrganizationName, &opts)
			if err != nil {
				log.Error("failed to list destination teams", "org", c.DestinationOrganizationName, "error", err)
				return err
			}

			srcTeams = append(srcTeams, items.Items...)
			log.Debug("fetched destination team page", "page", items.CurrentPage, "total_pages", items.TotalPages, "count", len(srcTeams))

			o.AddFormattedMessageCalculated("Found %d Teams", len(srcTeams))

			if items.CurrentPage >= items.TotalPages {
				break
			}
			opts.PageNumber = items.NextPage

		}
		o.AddTableHeaders("Name")
		log.Info("found teams", "org", c.DestinationOrganizationName, "count", len(srcTeams))
		for _, i := range srcTeams {

			o.AddTableRows(i.Name)
		}
	}

	return nil
}
