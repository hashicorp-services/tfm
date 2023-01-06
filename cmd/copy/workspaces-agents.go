package copy

import (
	"github.com/hashicorp-services/tfm/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
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

func createAgentPoolAssignment(c tfclient.ClientContexts) error {

	// Get Workspaces from Config
	srcAgentPoolsCfg := viper.GetStringSlice("srcagentpool")
	destAgentPoolsCfg := viper.GetStringSlice("destagentpool")

	o.AddFormattedMessageCalculated("Found %d source agent pools in Configuration", len(srcAgentPoolsCfg))
	o.AddFormattedMessageCalculated("Found %d destination agent pools in Configuration", len(destAgentPoolsCfg))

	// Get the source workspaces properties
	srcWorkspaces, err := discoverSrcWorkspaces(c)
	if err != nil {
		return errors.Wrap(err, "failed to list Workspaces from source while checking source agent pool IDs")
	}

	// For each source workspace with an execution mode of "agent", compare the source agent pool ID to the
	// user provided source pool ID. If they match, update the matching destination workspace with
	// the user provided agent pool ID that exists in the destination.
	for _, ws := range srcWorkspaces {
		isagent := checkExecution(c, ws)

		if !isagent {
			o.AddMessageUserProvided("No Agent Pool Assigned to source Workspace: ", ws.Name)
		} else {
			if ws.AgentPool != nil {
				if ws.AgentPool.ID != srcAgentPoolsCfg[0] {
					o.AddMessageUserProvided("Assigned Agent Pool ID does not match the user provided ID: ", srcAgentPoolsCfg[0])
				} else {
					o.AddFormattedMessageUserProvided2("Updating destination workspace %v execution mode to type agent and assigning pool ID %v", ws.Name, destAgentPoolsCfg[0])
					assignAgentPool(c, c.DestinationOrganizationName, destAgentPoolsCfg[0], ws.Name)
				}
			} else {
				o.AddMessageUserProvided("No Agent Pool Assigned to source Workspace: ", ws.Name)
			}

		}
	}
	return nil
}
