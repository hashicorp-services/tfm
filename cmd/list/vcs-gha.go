// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package list

import (
	"github.com/hashicorp-services/tfm/cmd/helper"
	"github.com/hashicorp-services/tfm/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

var (
	ghaVcsListCmd = &cobra.Command{
		Use:     "vcs-gha",
		Aliases: []string{"vcs-gha"},
		Short:   "List GHA VCS Providers",
		Long:    "List of GitHub App VCS Providers. Will default to source if no side is specified",
		Run: func(cmd *cobra.Command, args []string) {
			ghaVcsList(tfclient.GetClientContexts())
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			o.Close()
		},
	}
)

func init() {
	ListCmd.AddCommand(ghaVcsListCmd)
}

// helper functions
func ghaVcsListAllForOrganization(c tfclient.ClientContexts) ([]*tfe.GHAInstallation, error) {
	var allGHAItems []*tfe.GHAInstallation
	optsGHA := tfe.GHAInstallationListOptions{
		ListOptions: tfe.ListOptions{PageNumber: 1, PageSize: 100},
	}
	for {
		var ghaItems *tfe.GHAInstallationList
		var err error

		if (ListCmd.Flags().Lookup("side").Value.String() == "source") || (!ListCmd.Flags().Lookup("side").Changed) {
			ghaItems, err = c.SourceClient.GHAInstallations.List(c.SourceContext, &optsGHA)
		}

		if ListCmd.Flags().Lookup("side").Value.String() == "destination" {
			ghaItems, err = c.DestinationClient.GHAInstallations.List(c.DestinationContext, &optsGHA)
		}
		if err != nil {
			return nil, err
		}

		allGHAItems = append(allGHAItems, ghaItems.Items...)

		if ghaItems.CurrentPage >= ghaItems.TotalPages {
			break
		}
		optsGHA.PageNumber = ghaItems.NextPage
	}

	return allGHAItems, nil
}

// output functions
func ghaVcsList(c tfclient.ClientContexts) error {
	o.AddMessageUserProvided("List vcs for configured Organizations", "")

	var orgGhaVcsList []*tfe.GHAInstallation
	var err error

	if (ListCmd.Flags().Lookup("side").Value.String() == "source") || (!ListCmd.Flags().Lookup("side").Changed) {
		orgGhaVcsList, err = ghaVcsListAllForOrganization(c)
	}

	if ListCmd.Flags().Lookup("side").Value.String() == "destination" {
		orgGhaVcsList, err = ghaVcsListAllForOrganization(c)
	}

	if err != nil {
		helper.LogError(err, "failed to list vcs for organization")
	}

	o.AddFormattedMessageCalculated("Found %d vcs", len(orgGhaVcsList))

	o.AddTableHeaders("Name", "Installation ID", "ID")
	for _, i := range orgGhaVcsList {

		// The ID and Installation ID are flipped as they are flipped in the TFE/HCP TF UI, so we are matching the UI instead of the API/SDK
		o.AddTableRows(*i.Name, *i.ID, *i.InstallationID) 
	}

	return nil
}
