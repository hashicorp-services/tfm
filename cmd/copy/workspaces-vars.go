// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package copy

import (
	"fmt"

	"github.com/hashicorp-services/tfm/cmd/helper"
	"github.com/hashicorp-services/tfm/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
)

// Copys variables from a source workspace to a destination workspace
func variableCopy(c tfclient.ClientContexts, sourceWorkspaceID string, destinationWorkspaceID string, skipSensitive bool) error {

	variableListOpts := tfe.VariableListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100,
		},
	}

	//Get all variables in source the workspace
	srcWsVars, err := c.SourceClient.Variables.List(c.SourceContext, sourceWorkspaceID, &variableListOpts)
	if err != nil {
		fmt.Println("Could not list soruce Workspace variables.\n\n Error:", err.Error())
		return err
	}

	//Get all variables in destination the workspace
	destWsVars, err := c.DestinationClient.Variables.List(c.DestinationContext, destinationWorkspaceID, &variableListOpts)
	if err != nil {
		fmt.Println("Could not list destination Workspace variables.\n\n Error:", err.Error())
		return err
	}

	// For each variable in the source worksapce
	for _, workspaceVar := range srcWsVars.Items {
		destVarName := workspaceVar.Key

		//gather variables properties from source workspace. Variables marked as sensitive will be set to "" in the destination unless skipped
		variableOpts := tfe.VariableCreateOptions{
			Type:        "",
			Key:         &workspaceVar.Key,
			Value:       &workspaceVar.Value,
			Description: &workspaceVar.Description,
			Category:    &workspaceVar.Category,
			HCL:         &workspaceVar.HCL,
			Sensitive:   &workspaceVar.Sensitive,
		}

		// Check for the existence of the variable with the same Key name in the destination
		exists := doesVarExist(destVarName, destWsVars)

		// If the variable exists in the destination, do nothing and inform the user
		if exists {
			o.AddMessageUserProvided("Exists in destination will not migrate", destVarName)

		} else if skipSensitive {
			//Create the variable in the destination workspace but skip any sensitive ones
			if workspaceVar.Sensitive {
				o.AddMessageUserProvided(destVarName, "is sensitive and will not be copied")
			} else {
				//Create the variable in the destination workspace
				o.AddMessageUserProvided("Copying", destVarName)
				_, err := c.DestinationClient.Variables.Create(c.DestinationContext, destinationWorkspaceID, variableOpts)
				if err != nil {
					fmt.Println("Could not create Workspace variable.\n\n Error:", err.Error())
					return err
				}
			}
		} else {
			//Create the variable in the destination workspace
			o.AddMessageUserProvided("Copying", destVarName)
			_, err := c.DestinationClient.Variables.Create(c.DestinationContext, destinationWorkspaceID, variableOpts)
			if err != nil {
				fmt.Println("Could not create Workspace variable.\n\n Error:", err.Error())
				return err
			}
		}
	}

	return nil
}

// Compares the source variable key to existing destination variable keys
func doesVarExist(workspaceName string, v *tfe.VariableList) bool {
	for _, w := range v.Items {
		if workspaceName == w.Key {
			return true
		}
	}
	return false
}

// Main function used for --vars flag
func copyVariables(c tfclient.ClientContexts, skipSecure bool) error {

	// Get the source workspaces from the config file or ALL workspaces if non provided in the config file
	srcWorkspaces, err := getSrcWorkspacesCfg(c)
	if err != nil {
		return errors.Wrap(err, "Failed to list Workspaces from source")
	}

	// Get/Check if Workspace map exists
	wsMapCfg, err := helper.ViperStringSliceMap("workspaces-map")
	if err != nil {
		fmt.Println("Invalid input for workspaces-map")
	}

	// Get the destination workspaces
	destWorkspaces, err := discoverDestWorkspaces(tfclient.GetClientContexts(), true)
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

		// Check if the workspace name prefix and suffix are set
		if len(wsNamePrefix) > 0 || len(wsNameSuffix) > 0 {
			srcworkspaceSlice := []*tfe.Workspace{{Name: destWorkSpaceName}}
			newDestWorkspaceName := standardizeNamingConvention(srcworkspaceSlice, wsNamePrefix, wsNameSuffix)
			destWorkSpaceName = newDestWorkspaceName[0].Name
		}

		// Ensure the destination workspace exists in the destination target
		exists := doesWorkspaceExist(destWorkSpaceName, destWorkspaces)

		// If the destination workspace exists, get the ID of the workspace
		if exists {
			destWorkspaceId, err := getWorkspaceId(tfclient.GetClientContexts(), destWorkSpaceName)
			if err != nil {
				return errors.Wrap(err, "Failed to get the ID of the destination Workspace.")
			}

			fmt.Printf("Source ws %v has a matching ws %v in destination with ID %v. Comparing and copying existing variables...\n", srcworkspace.Name, destWorkSpaceName, destWorkspaceId)

			// Copy Variables from Source to Destination Workspace
			variableCopy(c, srcworkspace.ID, destWorkspaceId, skipSensitive)

			// Unlock the workspace
			unlockWorkspace(tfclient.GetClientContexts(), destWorkspaceId)
		} else {
			fmt.Printf("Source workspace named %v does not exist in destination. No variables to migrate\n", srcworkspace.Name)
		}
	}
	return nil
}
