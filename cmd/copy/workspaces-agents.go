package copy

import (
	"fmt"

	"github.com/hashicorp-services/tfm/tfclient"
	tfe "github.com/hashicorp/go-tfe"
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

	fmt.Printf("Agent Pool ID %v added to to workspace %v", destPoolId, workspace)

	return workspace, nil
}
