package copy

import (
	"fmt"

	"github.com/hashicorp-services/tfm/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
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

// Workflow
// 1. Get source variable sets and properties
// 2. Recreate the variable sets in the destination
// 3. Get the created destinatin variable set ID and name
// 4. Get the variable sets variables in the source variable set
// 5. Recreate the variable sets variables in the destination variable set

// Gets all variable sets in the source org and recreated them in the destination org.
// Does not create variable set variables.
// Also stores variable set ID for later use.

func discoverSrcVariableSets(c tfclient.ClientContexts, output bool) ([]*tfe.VariableSet, error) {
	varSets := []*tfe.VariableSet{}

	opts := tfe.VariableSetListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100,
		},
	}

	for {
		set, err := c.SourceClient.VariableSets.List(c.SourceContext, c.SourceOrganizationName, &opts)
		if err != nil {
			return nil, err
		}

		varSets = append(varSets, set.Items...)

		if output {
			o.AddFormattedMessageCalculated2("Found %d variable sets for org %v", len(varSets), c.SourceOrganizationName)
		}
		if set.CurrentPage >= set.TotalPages {
			break
		}
		opts.PageNumber = set.NextPage

	}
	return varSets, nil
}

func discoverDestVariableSets(c tfclient.ClientContexts, output bool) ([]*tfe.VariableSet, error) {
	varSets := []*tfe.VariableSet{}

	opts := tfe.VariableSetListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100,
		},
	}

	for {
		set, err := c.DestinationClient.VariableSets.List(c.DestinationContext, c.DestinationOrganizationName, &opts)
		if err != nil {
			return nil, err
		}

		varSets = append(varSets, set.Items...)

		if output {
			o.AddFormattedMessageCalculated2("Found %d variable sets for org %v", len(varSets), c.DestinationOrganizationName)
		}
		if set.CurrentPage >= set.TotalPages {
			break
		}
		opts.PageNumber = set.NextPage

	}
	return varSets, nil
}

func createVariableSets(c tfclient.ClientContexts, variableSet *tfe.VariableSet) (string, error) {

	o.AddMessageUserProvided("Copying variable sets for org:", c.SourceOrganizationName)

	//Create the variable sets
	varset, err := c.DestinationClient.VariableSets.Create(c.DestinationContext, c.DestinationOrganizationName, &tfe.VariableSetCreateOptions{
		Type:        "",
		Name:        &variableSet.Name,
		Description: &variableSet.Description,
		Global:      &variableSet.Global,
	})

	if err != nil {
		fmt.Println("Could not create variable set.\n\n Error:", err.Error())
		return "", err
	}

	srcVarSetName := variableSet.Name

	_ = varset

	o.AddDeferredMessageRead("Copied variable set: ", variableSet.Name)

	return srcVarSetName, nil
}

// Get the variable set ID and name.
func getVarSetID(c tfclient.ClientContexts, srcVarSetName string) (string, error) {
	var destVarSetID string

	opts := tfe.VariableSetListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100,
		},
		Include: "Name",
	}

	//get all variable sets
	set, err := c.DestinationClient.VariableSets.List(c.DestinationContext, c.DestinationOrganizationName, &opts)

	if err != nil {
		fmt.Println("Could not list variable sets while trying to get destination variable set ID.\n\n Error:", err.Error())
		return "", err
	}

	if set == nil {
		o.AddFormattedMessageUserProvided2("No variable set named %v exist in org: %v", srcVarSetName, c.DestinationOrganizationName)

	} else {
		for _, set := range set.Items {

			// Save the variable set ID for later.
			destVarSetID := set.ID

			return destVarSetID, err

		}
	}
	return destVarSetID, nil
}

// Gets all variables from the variable set ID provided in the source org and recreates them
// in the destination org and variable set.

func discoverSrcVariableSetVariables(c tfclient.ClientContexts, srcVarSetID string, srcVarSetName string) ([]*tfe.VariableSetVariable, error) {
	variables := []*tfe.VariableSetVariable{}

	opts := tfe.VariableSetVariableListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100,
		},
	}

	for {
		vars, err := c.SourceClient.VariableSetVariables.List(c.SourceContext, srcVarSetID, &opts)
		if err != nil {
			return nil, err
		}

		variables = append(variables, vars.Items...)

		o.AddFormattedMessageCalculated2("Found %d variables in source variable set %v", len(variables), srcVarSetName)

		if vars.CurrentPage >= vars.TotalPages {
			break
		}
		opts.PageNumber = vars.NextPage
	}

	return variables, nil
}

func discoverDestVariableSetVariables(c tfclient.ClientContexts, destVarSetID string, destVarSetName string) ([]*tfe.VariableSetVariable, error) {
	variables := []*tfe.VariableSetVariable{}

	opts := tfe.VariableSetVariableListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100,
		},
	}

	for {
		vars, err := c.DestinationClient.VariableSetVariables.List(c.DestinationContext, destVarSetID, &opts)
		if err != nil {
			return nil, err
		}

		variables = append(variables, vars.Items...)

		o.AddFormattedMessageCalculated2("Found %d variables in destination variable set %v", len(variables), destVarSetName)

		if vars.CurrentPage >= vars.TotalPages {
			break
		}
		opts.PageNumber = vars.NextPage
	}

	return variables, nil
}

func createVariableSetVars(c tfclient.ClientContexts, destVarSetID string, destVarSetName string, variables *tfe.VariableSetVariable) error {
	o.AddMessageUserProvided("Copying variables for variable set:", destVarSetName)

	//Create the variables in the variable set
	vars, err := c.DestinationClient.VariableSetVariables.Create(c.DestinationContext, destVarSetID, &tfe.VariableSetVariableCreateOptions{
		Type:        "",
		Key:         &variables.Key,
		Value:       &variables.Value,
		Description: &variables.Description,
		Category:    &variables.Category,
		HCL:         &variables.HCL,
		Sensitive:   &variables.Sensitive,
	})

	if err != nil {
		fmt.Println("Could not create variable.\n\n Error:", err.Error())
		return err
	}

	_ = vars

	return nil
}

func doesVariableSetExist(srcVarSetName string, destVarSets []*tfe.VariableSet) (string, bool) {
	var destVarSetID string
	for _, s := range destVarSets {
		destVarSetID := s.ID
		if srcVarSetName == s.Name {
			return destVarSetID, true
		}
	}
	return destVarSetID, false
}

func doesVariableSetVarExist(srcVarKey string, destVarName []*tfe.VariableSetVariable) bool {
	for _, v := range destVarName {
		if srcVarKey == v.Key {
			return true
		}
	}
	return false
}

// Creates variable sets and then copies the variables from source set to destination set.
func copyVariableSets(c tfclient.ClientContexts) error {
	srcVarSets, err := discoverSrcVariableSets(c, true)
	if err != nil {
		return errors.Wrap(err, "failed to list variable sets from source")
	}

	destVarSets, err := discoverDestVariableSets(c, true)
	if err != nil {
		return errors.Wrap(err, "failed to list variable sets from source")
	}

	for _, set := range srcVarSets {

		destVarSetID, exists := doesVariableSetExist(set.Name, destVarSets)
		if err != nil {
			return errors.Wrap(err, "failed to get destination variable set.")
		}

		_ = destVarSetID

		if exists {
			o.AddFormattedMessageUserProvided2("Variable set named %v exist in org: %v. Skipping creation.", set.Name, c.DestinationOrganizationName)
		} else {
			srcVarSetName, err := createVariableSets(c, set)
			if err != nil {
				return errors.Wrap(err, "Failed to create variable sets in the destination.")
			}

			_ = srcVarSetName
		}
	}
	copyVarSetVars(c)
	return nil
}

func copyVarSetVars(c tfclient.ClientContexts) error {
	srcVarSets, err := discoverSrcVariableSets(c, false)
	if err != nil {
		return errors.Wrap(err, "failed to list variable sets from source")
	}

	destVarSets, err := discoverDestVariableSets(c, false)
	if err != nil {
		return errors.Wrap(err, "failed to list variable sets from source")
	}

	for _, set := range srcVarSets {

		destVarSetID, exists := doesVariableSetExist(set.Name, destVarSets)
		if err != nil {
			return errors.Wrap(err, "failed to get destination variable set.")
		}

		if exists {
			o.AddFormattedMessageUserProvided2("Variable set named %v exist in org: %v. Checking variables.", set.Name, c.DestinationOrganizationName)

			srcvariables, err := discoverSrcVariableSetVariables(c, set.ID, set.Name)
			if err != nil {
				return errors.Wrap(err, "Failed to get variables for variable set.")
			}

			destvariable, err := discoverDestVariableSetVariables(c, destVarSetID, set.Name)
			if err != nil {
				return errors.Wrap(err, "Failed to get variables for variable set.")
			}

			for _, variable := range srcvariables {

				exists := doesVariableSetVarExist(variable.Key, destvariable)
				if exists {
					o.AddFormattedMessageUserProvided("Variable named %v exists in variable set. Skipping.", variable.Key)

				} else {
					createVariableSetVars(c, destVarSetID, set.Name, variable)
					if err != nil {
						return errors.Wrap(err, "Failed to create variable in variable set.")
					}

				}
			}
		}
	}
	return nil
}
