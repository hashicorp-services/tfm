package copy

import (
	"fmt"

	"github.com/hashicorp-services/tfm/cmd/helper"
	"github.com/hashicorp-services/tfm/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
)

// Copys variables from a source workspace to a destination workspace
func variableCopy(c tfclient.ClientContexts, sourceWorkspaceID string, destinationWorkspaceID string) error {

	variableListOpts := tfe.VariableListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100,
		},
	}

	//Get all variables in source the workspace
	srcWsVars, err := c.SourceClient.Variables.List(c.SourceContext, sourceWorkspaceID, &variableListOpts)
	if err != nil {
		fmt.Println("Could not list soruce workspace variables.\n\n Error:", err.Error())
		return err
	}

	//Get all variables in destination the workspace
	destWsVars, err := c.DestinationClient.Variables.List(c.DestinationContext, destinationWorkspaceID, &variableListOpts)
	if err != nil {
		fmt.Println("Could not list destination workspace variables.\n\n Error:", err.Error())
		return err
	}

	// For each variable in the source worksapce
	for _, workspaceVar := range srcWsVars.Items {
		destVarName := workspaceVar.Key

		//gather variables properties from source workspace. Variables marked as sensitive will be set to "" in the destination
		variableOpts := tfe.VariableCreateOptions{
			Type:        "",
			Key:         &workspaceVar.Key,
			Value:       &workspaceVar.Value,
			Description: &workspaceVar.Description,
			Category:    &workspaceVar.Category,
			HCL:         &workspaceVar.HCL,
			Sensitive:   &workspaceVar.Sensitive,
		}

		exists := doesVarExist(destVarName, destWsVars)

		if exists {
			o.AddMessageUserProvided("Exists in destination will not migrate", destVarName)
		} else {
			//Create the variable in the destination workspace
			_, err := c.DestinationClient.Variables.Create(c.DestinationContext, destinationWorkspaceID, variableOpts)
			if err != nil {
				fmt.Println("Could not create Workspace variable.\n\n Error:", err.Error())
				return err
			}
		}
	}

	return nil
}

func doesVarExist(workspaceName string, v *tfe.VariableList) bool {
	for _, w := range v.Items {
		if workspaceName == w.Key {
			return true
		}
	}
	return false
}

// Main function used for --vars flag
func copyVariables(c tfclient.ClientContexts) error {

	// Get the source workspaces from the config file or ALL workspaces if non provided in the config file
	srcWorkspaces, err := getSrcWorkspacesCfg(c)
	if err != nil {
		return errors.Wrap(err, "failed to list Workspaces from source")
	}

	// Get/Check if Workspace map exists
	wsMapCfg, err := helper.ViperStringSliceMap("workspace-map")
	if err != nil {
		fmt.Println("invalid input for workspace-map")
	}

	// Get the destination workspaces
	destWorkspaces, err := discoverDestWorkspaces(tfclient.GetClientContexts())
	if err != nil {
		return errors.Wrap(err, "failed to list Workspaces from source")
	}

	// For each workspace
	for _, srcworkspace := range srcWorkspaces {
		destWorkSpaceName := srcworkspace.Name

		// Check if the destination Workspace name differs from the source name
		if len(wsMapCfg) > 0 {
			destWorkSpaceName = wsMapCfg[srcworkspace.Name]
		}

		// Ensure the destination workspace exists in the destination target
		exists := doesWorkspaceExist(destWorkSpaceName, destWorkspaces)

		// If the destination workspace exists, get the ID of the workspace
		if exists {
			destWorkspaceId, err := getWorkspaceId(tfclient.GetClientContexts(), destWorkSpaceName)
			if err != nil {
				return errors.Wrap(err, "Failed to get the ID of the destination Workspace.")
			}

			fmt.Printf("Source ws %v has a matching ws %v in destination with ID %v. Comparing existing variables...\n", srcworkspace.Name, destWorkSpaceName, destWorkspaceId)

			// Copy Variables from Source to Destination Workspace
			variableCopy(c, srcworkspace.ID, destWorkspaceId)

			// Unlock the workspace
			unlockWorkspace(tfclient.GetClientContexts(), destWorkspaceId)
		} else {
			fmt.Printf("Source workspace named %v does not exist in destination. No variables to migrate\n", srcworkspace.Name)
		}
	}
	return nil
}
