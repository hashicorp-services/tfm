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
		Long:    "List of VCS Providers. Will default to source if no side is specified",
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

// helper functions
func vcsListAllForOrganization(c tfclient.ClientContexts, orgName string) ([]*tfe.OAuthClient, error) {
	var allItems []*tfe.OAuthClient
	opts := tfe.OAuthClientListOptions{
		ListOptions: tfe.ListOptions{PageNumber: 1, PageSize: 100},
	}
	for {
		var items *tfe.OAuthClientList
		var err error
	
		if (ListCmd.Flags().Lookup("side").Value.String() == "source") || (!ListCmd.Flags().Lookup("side").Changed)  {
			items, err = c.SourceClient.OAuthClients.List(c.SourceContext, orgName, &opts)
		} 
		
		if ListCmd.Flags().Lookup("side").Value.String() == "destination" {
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

func organizationListAll(c tfclient.ClientContexts) ([]*tfe.Organization, error) {
	allItems := []*tfe.Organization{}
	opts := tfe.OrganizationListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100},
	}
	for {
		var items *tfe.OrganizationList
		var err error

		

		if (ListCmd.Flags().Lookup("side").Value.String() == "source") || (!ListCmd.Flags().Lookup("side").Changed) {

			items, err = c.SourceClient.Organizations.List(c.SourceContext, &opts)
			if err != nil {
				return nil, err
			}
	
			allItems = append(allItems, items.Items...)
			if items.CurrentPage >= items.TotalPages {
				break
			}
			opts.PageNumber = items.NextPage
		} 
		if ListCmd.Flags().Lookup("side").Value.String() == "destination" {
			items, err = c.DestinationClient.Organizations.List(c.DestinationContext, &opts)
			if err != nil {
				return nil, err
			}
	
			allItems = append(allItems, items.Items...)
			if items.CurrentPage >= items.TotalPages {
				break
			}
			opts.PageNumber = items.NextPage
		} 
	}

	return allItems, nil
}

// output functions
func vcsListAll(c tfclient.ClientContexts) error {
	o.AddMessageUserProvided("List vcs for all available Organizations", "")

	allOrgs, err := organizationListAll(c)
	if err != nil {
		helper.LogError(err, "failed to list organizations")
	}

	var allVcsList []*tfe.OAuthClient

	for _, v := range allOrgs {
		vcsList, err := vcsListAllForOrganization(c, v.Name)
		if err != nil {
			helper.LogError(err, "failed to list vcs for organization")
		}

		allVcsList = append(allVcsList, vcsList...)
	}

	o.AddFormattedMessageCalculated("Found %d vcs", len(allVcsList))

	o.AddTableHeaders("Organization", "Name", "Id", "Service Provider", "Service Provider Name", "Created At", "URL")
	for _, i := range allVcsList {

		vcsName := ""
		if i.Name != nil {
			vcsName = *i.Name
		}

		o.AddTableRows(i.Organization.Name, vcsName, i.ID, i.ServiceProvider, i.ServiceProviderName, i.CreatedAt, i.HTTPURL)
	}

	return nil
}