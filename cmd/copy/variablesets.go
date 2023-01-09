package copy

import (
	"fmt"

	"github.com/hashicorp-services/tfm/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

var (

	// `tfm copy varsets` command
	varSetCopyCmd = &cobra.Command{
		Use:   "varsets",
		Short: "Copy Variable Sets",
		Long:  "Copy Variable Sets from source to destination org",
		RunE: func(cmd *cobra.Command, args []string) error {
			return copyVariableSets(
				tfclient.GetClientContexts())

		},
		PostRun: func(cmd *cobra.Command, args []string) {
			o.Close()
		},
	}
)

func init() {

	// `tfm copy varsets all` command
	varSetCopyCmd.Flags().BoolP("all", "a", false, "Copy all variable sets (optional)")

	// Add commands
	CopyCmd.AddCommand(varSetCopyCmd)

}

func createVariableSets(c tfclient.ClientContexts) error {

	opts := tfe.VariableSetListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100,
		},
	}

	//get all variable sets
	set, err := c.SourceClient.VariableSets.List(c.SourceContext, c.SourceOrganizationName, &opts)

	if err != nil {
		fmt.Println("Could not list variable sets.\n\n Error:", err.Error())
		return err
	}

	if set == nil {
		o.AddMessageUserProvided("No variable sets exist in org:", c.SourceOrganizationName)

	} else {
		o.AddMessageUserProvided("Copying variable sets for org:", c.SourceOrganizationName)
		for _, set := range set.Items {

			//gather variable sets from source org.
			varSet := tfe.VariableSetCreateOptions{
				Type:        "",
				Name:        &set.Name,
				Description: &set.Description,
				Global:      &set.Global,
			}

			//Create the variable sets
			_, err := c.DestinationClient.VariableSets.Create(c.DestinationContext, c.DestinationOrganizationName, &varSet)

			if err != nil {
				fmt.Println("Could not create variable set.\n\n Error:", err.Error())
				return err
			}

			o.AddDeferredMessageRead("Copied variable set: ", set.Name)
		}
	}
	return nil
}

func copyVariableSets(c tfclient.ClientContexts) error {
	createVariableSets(c)
	return nil
}
