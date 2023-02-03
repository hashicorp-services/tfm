package list

import (
	"fmt"

	"github.com/hashicorp-services/tfm/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

var (

	// `tfm list ssh` command
	sshListCmd = &cobra.Command{
		Use:   "ssh",
		Short: "ssh-keys command",
		Long:  "Lists the ssh-keys for an org",
		// RunE: func(cmd *cobra.Command, args []string) error {
		// 	return listTeams(
		// 		tfeclient.GetClientContexts())

		// },
		Run: func(cmd *cobra.Command, args []string) {
			// return orgShow(
			// 	viper.GetString("name"))
			listSrcSSHKeys(tfclient.GetClientContexts())
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			o.Close()
		},
	}
)

func init() {
	// Flags().StringP, etc... - the "P" gives us the option for a short hand

	// `tfm list ssh all` command
	//sshListCmd.Flags().BoolP("all", "a", false, "List all? (optional)")

	// Add commands
	ListCmd.AddCommand(sshListCmd)

}

func listSrcSSHKeys(c tfclient.ClientContexts) error {

	keys := []*tfe.SSHKey{}

	opts := tfe.SSHKeyListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100,
		},
	}

	if (ListCmd.Flags().Lookup("side").Value.String() == "source") || (!ListCmd.Flags().Lookup("side").Changed) {

		o.AddMessageUserProvided("Getting list of SSH keys from: ", c.SourceHostname)

		for {
			k, err := c.SourceClient.SSHKeys.List(c.SourceContext, c.SourceOrganizationName, &opts)
			if err != nil {
				return err
			}

			fmt.Println()
			keys = append(keys, k.Items...)

			o.AddFormattedMessageCalculated("Found %d SSH keys", len(keys))

			if k.CurrentPage >= k.TotalPages {
				break
			}
			opts.PageNumber = k.NextPage

		}
		o.AddTableHeaders("Key Name", "Key ID")
		for _, i := range keys {

			o.AddTableRows(i.Name, i.ID)

		}
	}

	if ListCmd.Flags().Lookup("side").Value.String() == "destination" {

		o.AddMessageUserProvided("Getting list of SSH keys from: ", c.DestinationHostname)

		for {
			k, err := c.DestinationClient.SSHKeys.List(c.DestinationContext, c.DestinationOrganizationName, &opts)
			if err != nil {
				return err
			}

			fmt.Println()
			keys = append(keys, k.Items...)

			o.AddFormattedMessageCalculated("Found %d SSH keys", len(keys))

			if k.CurrentPage >= k.TotalPages {
				break
			}
			opts.PageNumber = k.NextPage

		}
		o.AddTableHeaders("Key Name", "Key ID")
		for _, i := range keys {

			o.AddTableRows(i.Name, i.ID)

		}
	}

	return nil
}
