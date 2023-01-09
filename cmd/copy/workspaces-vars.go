package copy

import (
	"fmt"

	"github.com/hashicorp-services/tfm/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
)

func variableCopy(c tfclient.ClientContexts, sourceWorkspaceID string, destinationWorkspaceID string) error {

	variableListOpts := tfe.VariableListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100,
		},
	}

	//get all variables in workspace
	workspaceVars, err := c.SourceClient.Variables.List(c.SourceContext, sourceWorkspaceID, &variableListOpts)

	if err != nil {
		fmt.Println("Could not list workspace variables.\n\n Error:", err.Error())
		return err
	}

	for _, workspaceVar := range workspaceVars.Items {

		//gather variables from source workspace. Variables marked as sensitive will be set to "" in the destination
		variableOpts := tfe.VariableCreateOptions{
			Type:        "",
			Key:         &workspaceVar.Key,
			Value:       &workspaceVar.Value,
			Description: &workspaceVar.Description,
			Category:    &workspaceVar.Category,
			HCL:         &workspaceVar.HCL,
			Sensitive:   &workspaceVar.Sensitive,
		}

		//Create the variable
		_, err := c.DestinationClient.Variables.Create(c.DestinationContext, destinationWorkspaceID, variableOpts)

		if err != nil {
			fmt.Println("Could not create Workspace variable.\n\n Error:", err.Error())
			return err
		}
	}

	return nil
}

func copyVariables(c tfclient.ClientContexts) error {
	// Get the source workspaces properties
	srcWorkspaces, err := getSrcWorkspacesCfg(c)
	if err != nil {
		return errors.Wrap(err, "failed to list Workspaces from source")
	}

	destWorkspaces, err := discoverDestWorkspaces(tfclient.GetClientContexts())
	if err != nil {
		return errors.Wrap(err, "failed to list Workspaces from source")
	}

	for _, srcworkspace := range srcWorkspaces {
		exists := doesWorkspaceExist(srcworkspace.Name, destWorkspaces)

		if exists {
			destWorkspaceId, err := getWorkspaceId(tfclient.GetClientContexts(), srcworkspace.Name)
			if err != nil {
				return errors.Wrap(err, "Failed to get the ID of the destination Workspace that matches the Name of the Source Workspace")
			}

			fmt.Printf("Source ws %v has a matching ws %v in destination with ID %v. Comparing existing variables...\n", srcworkspace.Name, srcworkspace.Name, destWorkspaceId)

			// Copy Variables from Source to Destination Workspace
			variableCopy(c, srcworkspace.ID, destWorkspaceId)

			unlockWorkspace(tfclient.GetClientContexts(), destWorkspaceId)
		} else {
			fmt.Printf("Source workspace named %v does not exist in destination. No variables to migrate\n", srcworkspace.Name)
		}
	}
	return nil
}
