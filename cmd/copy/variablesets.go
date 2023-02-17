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

			// Validate the map if it exists
			valid, varsets, err := validateMap(tfclient.GetClientContexts(), "varsets-map")
			if err != nil {
				return err
			}

			// If the map does not exist, or is not valid, copy all variable sets from the source target
			if !valid {
				return copyVariableSetsAll(
					tfclient.GetClientContexts())

				// If the map exists and is valid, copy only the variable sets specified with the new desired name specified on the right side of the `varsets-map`
				// from the configuration file
			} else {
				return copyVariableSetsCfg(
					tfclient.GetClientContexts(), varsets)
			}

		},
		PostRun: func(cmd *cobra.Command, args []string) {
			o.Close()
		},
	}
)

func init() {

	// Add commands
	CopyCmd.AddCommand(varSetCopyCmd)

}

// Workflow
// 1. Get source variable sets and properties
// 2. Recreate the variable sets in the destination
// 3. Get the created destinatin variable set ID and name
// 4. Get the variable sets variables in the source variable set
// 5. Recreate the variable sets variables in the destination variable set

// Get source target variable sets
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
			o.AddFormattedMessageCalculated2("Found %d variable sets for source org %v", len(varSets), c.SourceOrganizationName)
		}
		if set.CurrentPage >= set.TotalPages {
			break
		}
		opts.PageNumber = set.NextPage

	}
	return varSets, nil
}

// Get destination target variable sets
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
			o.AddFormattedMessageCalculated2("Found %d variable sets for destination org %v", len(varSets), c.DestinationOrganizationName)
		}
		if set.CurrentPage >= set.TotalPages {
			break
		}
		opts.PageNumber = set.NextPage

	}
	return varSets, nil
}

// Function that creates variable sets
func createVariableSets(c tfclient.ClientContexts, variableSet *tfe.VariableSet, destSetNameCfg string, useCfg bool) (string, error) {

	// If no config file `varsets-map` is specified
	if !useCfg {
		o.AddFormattedMessageUserProvided2("Copying variable set %v from source org %v", variableSet.Name, c.SourceOrganizationName)

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

	} else {
		o.AddFormattedMessageUserProvided3("Copying variable set %v from source org %v with new name %v", variableSet.Name, c.SourceOrganizationName, destSetNameCfg)

		//Create the variable sets
		varset, err := c.DestinationClient.VariableSets.Create(c.DestinationContext, c.DestinationOrganizationName, &tfe.VariableSetCreateOptions{
			Type:        "",
			Name:        &destSetNameCfg,
			Description: &variableSet.Description,
			Global:      &variableSet.Global,
		})

		if err != nil {
			fmt.Println("Could not create variable set.\n\n Error:", err.Error())
			return "", err
		}

		_ = varset

		o.AddDeferredMessageRead("Copied variable set: ", destSetNameCfg)

		return destSetNameCfg, nil
	}
}

// Get the variable set ID
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

// Gets all variables from the source variable set ID provided
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

// Gets all variables from the destination variable set ID provided
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

// Function that creates variables in the destination target variable set.
func createVariableSetVars(c tfclient.ClientContexts, destVarSetID string, destVarSetName string, variables *tfe.VariableSetVariable) error {
	o.AddFormattedMessageUserProvided2("Copying variable %v for variable set %v", variables.Key, destVarSetName)

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

// Check the destination variable set existence
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

// Check the destination variable set variable for existence
func doesVariableSetVarExist(srcVarKey string, destVarName []*tfe.VariableSetVariable) bool {
	for _, v := range destVarName {
		if srcVarKey == v.Key {
			return true
		}
	}
	return false
}

// Main function for `--varsets` flag when no map is specified in the config file
// Creates variable sets and then copies the variables from source set to destination set.
// If the user specifies a "varsets-map" list in the config file, only those variable sets will
// be copied. If they do not, all variable sets will be copied.
func copyVariableSetsAll(c tfclient.ClientContexts) error {

	// Get all source target var sets
	srcVarSets, err := discoverSrcVariableSets(c, true)
	if err != nil {
		return errors.Wrap(err, "failed to list variable sets from source")
	}

	// Get all destination target varsets
	destVarSets, err := discoverDestVariableSets(c, true)
	if err != nil {
		return errors.Wrap(err, "failed to list variable sets from destination")
	}

	// for each variable set in source
	for _, set := range srcVarSets {

		// Check for the existence of the variable set in the destination
		// Also get the ID if it does exist
		// Names must match between source and destination variable set
		destVarSetID, exists := doesVariableSetExist(set.Name, destVarSets)
		if err != nil {
			return errors.Wrap(err, "failed to get destination variable set.")
		}

		_ = destVarSetID

		// If it exists, inform the user.
		if exists {
			o.AddFormattedMessageUserProvided2("Variable set named %v exist in destination org: %v. Skipping creation.", set.Name, c.DestinationOrganizationName)
		} else {

			// Create a copy of the variable set in the destination
			srcVarSetName, err := createVariableSets(c, set, "", false)
			if err != nil {
				return errors.Wrap(err, "Failed to create variable sets in the destination org.")
			}

			_ = srcVarSetName
		}
	}

	// Copy all of the variables for the variable sets
	copyVarSetVars(c)

	return nil
}

// Copys the variable set variables for all variable sets
func copyVarSetVars(c tfclient.ClientContexts) error {
	srcVarSets, err := discoverSrcVariableSets(c, false)
	if err != nil {
		return errors.Wrap(err, "failed to list variable sets from source")
	}

	destVarSets, err := discoverDestVariableSets(c, false)
	if err != nil {
		return errors.Wrap(err, "failed to list variable sets from destination")
	}

	for _, set := range srcVarSets {

		destVarSetID, exists := doesVariableSetExist(set.Name, destVarSets)
		if err != nil {
			return errors.Wrap(err, "failed to get destination variable set ID.")
		}

		if exists {
			o.AddFormattedMessageUserProvided2("Variable set named %v exist in org: %v. Checking variables.", set.Name, c.DestinationOrganizationName)

			srcvariables, err := discoverSrcVariableSetVariables(c, set.ID, set.Name)
			if err != nil {
				return errors.Wrap(err, "Failed to get variables for source variable set.")
			}

			destvariable, err := discoverDestVariableSetVariables(c, destVarSetID, set.Name)
			if err != nil {
				return errors.Wrap(err, "Failed to get variables for destination variable set.")
			}

			for _, variable := range srcvariables {

				exists := doesVariableSetVarExist(variable.Key, destvariable)
				if exists {
					o.AddFormattedMessageUserProvided("Variable named %v exists in destination variable set. Skipping.", variable.Key)

				} else {
					createVariableSetVars(c, destVarSetID, set.Name, variable)
					if err != nil {
						return errors.Wrap(err, "Failed to create variable in variable destination set.")
					}
				}
			}
		}
	}
	return nil
}

// Main function for `--varsets` flag when a map is specified in the config file
// If variablesets-map is defined in the config file this function will be used. Because we cannot
// list variable sets filtered by name, a seperate function is required.
// Takes the varsets-map list from the config file after it has been converted to a map by viper.
func copyVariableSetsCfg(c tfclient.ClientContexts, varsets map[string]string) error {
	srcVarSets, err := discoverSrcVariableSets(c, true)
	if err != nil {
		return errors.Wrap(err, "failed to list variable sets from source")
	}

	destVarSets, err := discoverDestVariableSets(c, true)
	if err != nil {
		return errors.Wrap(err, "failed to list variable sets from destination")
	}

	// for each key:element pair in the varsets map, assign the key and element to vars for readability
	for key, element := range varsets {
		srcsetname := key
		destsetname := element

		// Check for the existence of all the provided source variable sets
		srcVarSetID, exists := doesVariableSetExist(srcsetname, srcVarSets)
		if err != nil {
			return errors.Wrap(err, "failed to get source variable set.")
		}

		_ = srcVarSetID

		// If the variable set doesn't exist in the source, inform the user.
		if !exists {
			o.AddFormattedMessageUserProvided2("Variable Set named %v does not exist in source org %v. Skipping.", srcsetname, c.SourceOrganizationName)
		} else {

			// For each variable set in the source
			for _, set := range srcVarSets {

				// Does the source variable set name match the one provided in the config file?
				// If not, do nothing.
				if set.Name != srcsetname {
					continue

				} else {
					destVarSetID, exists := doesVariableSetExist(destsetname, destVarSets)
					if err != nil {
						return errors.Wrap(err, "failed to get destination variable set.")
					}

					_ = destVarSetID

					if exists {
						o.AddFormattedMessageUserProvided2("Variable set named %v exist in destination org: %v. Skipping creation.", set.Name, c.DestinationOrganizationName)
					} else {
						srcVarSetName, err := createVariableSets(c, set, destsetname, true)
						if err != nil {
							return errors.Wrap(err, "Failed to create variable sets in the destination org.")
						}

						_ = srcVarSetName
					}
				}
			}
		}
	}
	copyVarSetVarsCfg(c, varsets)

	return nil
}

// Only copy the variables for the variable sets defined by the user in the configuration.
func copyVarSetVarsCfg(c tfclient.ClientContexts, varsets map[string]string) error {
	srcVarSets, err := discoverSrcVariableSets(c, false)
	if err != nil {
		return errors.Wrap(err, "failed to list variable sets from source")
	}

	destVarSets, err := discoverDestVariableSets(c, false)
	if err != nil {
		return errors.Wrap(err, "failed to list variable sets from destination")
	}

	// for each key:element pair in the varsets map, assign the key and element to vars for readability
	for key, element := range varsets {
		srcsetname := key
		destsetname := element

		// Check for the existence of all the provided source variable sets
		srcVarSetID, exists := doesVariableSetExist(srcsetname, srcVarSets)
		if err != nil {
			return errors.Wrap(err, "failed to get source variable set.")
		}

		_ = srcVarSetID

		// If the variable set doesn't exist in the source, inform the user.
		if !exists {
			o.AddFormattedMessageUserProvided2("Variable Set named %v does not exist in source org %v. Skipping.", srcsetname, c.SourceOrganizationName)
		} else {
			for _, set := range srcVarSets {

				// Does the source variable set name match the one provided in the config file?
				// If not, do nothing.
				if set.Name != srcsetname {
					continue

				} else {
					destVarSetID, exists := doesVariableSetExist(destsetname, destVarSets)
					if err != nil {
						return errors.Wrap(err, "failed to get destination variable set ID.")
					}
					if exists {
						o.AddFormattedMessageUserProvided2("Variable set named %v exist in org: %v. Checking variables.", set.Name, c.DestinationOrganizationName)

						srcvariables, err := discoverSrcVariableSetVariables(c, set.ID, set.Name)
						if err != nil {
							return errors.Wrap(err, "Failed to get variables for source variable set.")
						}

						destvariable, err := discoverDestVariableSetVariables(c, destVarSetID, destsetname)
						if err != nil {
							return errors.Wrap(err, "Failed to get variables for destination variable set.")
						}

						for _, variable := range srcvariables {

							exists := doesVariableSetVarExist(variable.Key, destvariable)
							if exists {
								o.AddFormattedMessageUserProvided("Variable named %v exists in destination variable set. Skipping.", variable.Key)

							} else {
								createVariableSetVars(c, destVarSetID, destsetname, variable)
								if err != nil {
									return errors.Wrap(err, "Failed to create variable in destination variable set.")
								}
							}
						}
					}
				}
			}
		}
	}
	return nil
}
