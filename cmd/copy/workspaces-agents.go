package copy

import (

	"github.com/hashicorp-services/tfm/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
)

// All functions related to copying/assigning agent pools to workspaces

// Check workspace properties for execution type.
func checkExecution(c tfclient.ClientContexts, ws *tfe.Workspace) bool {
	if ws.ExecutionMode == "agent" {
		return true
	}
	return false
}

// Update workspace execution mode to agent and assign an agent pool ID to a workspace.
func assignAgentPool(c tfclient.ClientContexts, org string, destPoolId string, ws string) (*tfe.Workspace, error) {

	executionMode := "agent"

	opts := tfe.WorkspaceUpdateOptions{
		Type:          "",
		AgentPoolID:   &destPoolId,
		ExecutionMode: &executionMode,
	}

	workspace, err := c.DestinationClient.Workspaces.Update(c.DestinationContext, c.DestinationOrganizationName, ws, opts)
	if err != nil {
		return nil, err
	}

	return workspace, nil
}

func createAgentPoolAssignment(c tfclient.ClientContexts, agentpools map[string]string) error {

	for key, element := range agentpools {
		srcpool := key
		destpool := element

		// Get the source workspaces properties
		srcWorkspaces, err := getSrcWorkspacesCfg(c)
		if err != nil {
			return errors.Wrap(err, "failed to list Workspaces from source while checking source agent pool IDs")
		}

		// Get/Check if Workspace map exists
		wsMapCfg, err := helper.ViperStringSliceMap("workspace-map")
		if err != nil {
			fmt.Println("invalid input for workspace-map")
		}

		// For each source workspace with an execution mode of "agent", compare the source agent pool ID to the
		// user provided source pool ID. If they match, update the matching destination workspace with
		// the user provided agent pool ID that exists in the destination.
		for _, ws := range srcWorkspaces {
			isagent := checkExecution(c, ws)
			destWorkSpaceName := ws.Name

			// Check if Destination Workspace Name to be Change
			if len(wsMapCfg) > 0 {
				destWorkSpaceName = wsMapCfg[ws.Name]
			}

			if !isagent {
				o.AddMessageUserProvided("No Agent Pool Assigned to source Workspace: ", ws.Name)
			} else {
				if ws.AgentPool != nil {
					if ws.AgentPool.ID != srcpool {
						o.AddFormattedMessageUserProvided2("Workspace %v assigned agent pool ID does not match provided source ID %v. Skipping.", ws.Name, srcpool)
					} else {
						o.AddFormattedMessageUserProvided2("Updating destination workspace %v execution mode to type agent and assigning pool ID %v", destWorkSpaceName, destpool)
						assignAgentPool(c, c.DestinationOrganizationName, destpool, destWorkSpaceName)
					}
				} else {
					o.AddMessageUserProvided("No Agent Pool Assigned to source Workspace: ", ws.Name)
				}

			}
		}
	}
	return nil
}
