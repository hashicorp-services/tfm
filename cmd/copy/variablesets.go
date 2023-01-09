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

// Gets all variable sets in the source org and recreated them in the destination org.
// Does not create variable set variables.
// Also stores variable set ID for later use.
func createVariableSets(c tfclient.ClientContexts) (string, error) {
	var varSetID string

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
		return "", err
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
				return "", err
			}

			// Save the variable set ID for later.
			varSetID := set.ID

			return varSetID, err

			o.AddDeferredMessageRead("Copied variable set: ", set.Name)
		}
	}
	return varSetID, nil
}

// Gets all variables from the variable set ID provided in the source org and recreates them
// in the destination org and variable set.
func createVariableSetVars(c tfclient.ClientContexts, varSetID string) error {
	opts := tfe.VariableSetVariableListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100,
		},
	}

	//Get all variables in the variable set
	vars, err := c.SourceClient.VariableSetVariables.List(c.SourceContext, varSetID, &opts)

	if err != nil {
		fmt.Println("Could not list variables.\n\n Error:", err.Error())
		return err
	}

	if vars == nil {
		o.AddMessageUserProvided("No variables exist for variable set:", varSetName)

	} else {
		o.AddMessageUserProvided("Copying variable sets for org:", c.SourceOrganizationName)
		for _, vars := range vars.Items {

			//gather variable sets from source org.
			varSet := tfe.VariableSetVariableCreateOptions{
				Type:        "",
				Key:         new(string),
				Value:       new(string),
				Description: new(string),
				Category:    &"",
				HCL:         new(bool),
				Sensitive:   new(bool),
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

// Creates variable sets and then copies the variables from source set to destination set.

// Copy variable sets - no variables
func copyVariableSets(c tfclient.ClientContexts) error {
	createVariableSets(c)
	return nil
}
