package list

import (
	"github.com/hashicorp-services/tfe-mig/cmd/helper"
	"github.com/hashicorp-services/tfe-mig/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

var (
	//vcsOutput output.Output

	vcsListCmd = &cobra.Command{
		Use:     "vcs",
		Aliases: []string{"vcs"},
		Short:   "List VCS Providers",
		Long:    "List of VCS Providers.",
		Run: func(cmd *cobra.Command, args []string) {
			vcsListAll(tfclient.GetClientContexts())
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			o.Close()
		},
	}
)

func init() {

	ListCmd.AddCommand(vcsListCmd)
}

func vcsListAllForOrganization(c tfclient.ClientContexts, side string, orgName string) ([]*tfe.OAuthClient, error) {
	var allItems []*tfe.OAuthClient
	opts := tfe.OAuthClientListOptions{
		ListOptions: tfe.ListOptions{PageNumber: 1, PageSize: 100},
	}
	for {

		var items *tfe.OAuthClientList
		var err error
		
		if side == "source" {
			items, err = c.SourceClient.OAuthClients.List(c.SourceContext, orgName, &opts)	
		} else {
			items, err = c.DestinationClient.OAuthClients.List(c.DestinationContext, orgName, &opts)
		}
		if err != nil {
			return nil, err
		}

		allItems = append(allItems, items.Items...)

		if items.CurrentPage >= items.TotalPages {
			break
		}
		opts.PageNumber = items.NextPage
	}


	return allItems, nil
}

// output functions
func vcsListAll(c tfclient.ClientContexts) error {
	o.AddMessageUserProvided("List vcs for all available Organizations in source and destination", "")
	
	sourceOrgs, serr := organizationListAllSource(c)
	if serr != nil {
		helper.LogError(serr, "failed to list organizations")
	}
	destinationOrgs, derr := organizationListAllDestination(c)
	if derr != nil {
		helper.LogError(derr, "failed to list organizations")
	}

	var sourceAllVcsList []*tfe.OAuthClient
	var destinationAllVcsList []*tfe.OAuthClient

	for _, v := range sourceOrgs {
		vcsList, err := vcsListAllForOrganization(c, "source", v.Name)
		if err != nil {
			helper.LogError(err, "failed to list vcs for organization")
		}

		sourceAllVcsList = append(sourceAllVcsList, vcsList...)
	}

	for _, v := range destinationOrgs {
		vcsList, err := vcsListAllForOrganization(c, "destination", v.Name)
		if err != nil {
			helper.LogError(err, "failed to list vcs for organization")
		}

		destinationAllVcsList = append(destinationAllVcsList, vcsList...)
	}

	o.AddFormattedMessageCalculated("Found %d vcs", len(sourceAllVcsList)+len(destinationAllVcsList))

	o.AddTableHeaders("Hostname","Organization", "Name", "Id", "Service Provider", "Service Provider Name", "Created At", "URL")
	for _, i := range sourceAllVcsList {

		vcsName := ""
		if i.Name != nil {
			vcsName = *i.Name
		}

		o.AddTableRows(c.SourceHostname, i.Organization.Name, vcsName, i.ID, i.ServiceProvider, i.ServiceProviderName, i.CreatedAt, i.HTTPURL)
	}

	for _, i := range destinationAllVcsList {

		vcsName := ""
		if i.Name != nil {
			vcsName = *i.Name
		}

		o.AddTableRows(c.DestinationHostname, i.Organization.Name, vcsName, i.ID, i.ServiceProvider, i.ServiceProviderName, i.CreatedAt, i.HTTPURL)
	}

	return nil
}
