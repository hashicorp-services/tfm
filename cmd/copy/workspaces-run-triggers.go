// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package copy

import (
	"github.com/hashicorp-services/tfm/cmd/helper"
	"github.com/hashicorp-services/tfm/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
)

func copyRunTriggers(c tfclient.ClientContexts) error {
	// Get the source workspaces from the config file or ALL workspaces if non provided in the config file
	srcWorkspaces, err := getSrcWorkspacesCfg(c)
	if err != nil {
		return errors.Wrap(err, "Failed to list Workspaces from source")
	}

	// Get/Check if Workspace map exists
	wsMapCfg, err := helper.ViperStringSliceMap("workspaces-map")
	if err != nil {
		o.AddErrorUserProvided("Invalid input for workspaces-map")
	}

	// Get the destination workspaces
	destWorkspaces, err := discoverDestWorkspaces(tfclient.GetClientContexts(), true)
	if err != nil {
		return errors.Wrap(err, "failed to list Workspaces from destination")
	}

	// For each workspace
	for _, srcWorkspace := range srcWorkspaces {
		destWorkSpaceName := srcWorkspace.Name

		// Check if the destination Workspace name differs from the source name
		if len(wsMapCfg) > 0 {
			destWorkSpaceName = wsMapCfg[srcWorkspace.Name]
		}

		destWorkSpaceID, err := getWorkspaceId(tfclient.GetClientContexts(), destWorkSpaceName)
		if err != nil {
			return errors.Wrap(err, "Failed to get the ID of the destination Workspace that matches the Name of the Source Workspace: "+srcWorkspace.Name)
		}

		// Ensure the destination workspace exists in the destination target
		exists := doesWorkspaceExist(destWorkSpaceName, destWorkspaces)

		// If the destination workspace exists
		if exists {

			opts := tfe.WorkspaceUpdateOptions{
				Type:                "",
				AutoApplyRunTrigger: &srcWorkspace.AutoApplyRunTrigger,
			}

			// well we want to first set the destination workspace auto apply run trigger setting to match the source workspace // asdfasdfasdfasdf
			o.AddFormattedMessageUserProvided2("Setting %v workspace's Run Trigger Auto Apply setting to %v", destWorkSpaceName, srcWorkspace.AutoApplyRunTrigger)
			_, err := c.DestinationClient.Workspaces.Update(c.DestinationContext, c.DestinationOrganizationName, destWorkSpaceName, opts)
			if err != nil {
				return errors.Wrap(err, "failed to set workspace "+destWorkSpaceName+" auto apply run trigger")
			}

			// List all run triggers for the source workspace
			inboundOpts := tfe.RunTriggerListOptions{
				ListOptions: tfe.ListOptions{
					PageNumber: 1,
					PageSize:   100},
				RunTriggerType: "inbound",
			}

			inboundRunTriggers, err := c.SourceClient.RunTriggers.List(c.SourceContext, srcWorkspace.ID, &inboundOpts)
			if err != nil {
				return errors.Wrap(err, "Failed to list run triggers for source workspace: "+srcWorkspace.Name)
			}

			// ok got a list of run triggers for the source workspace
			for _, runTrigger := range inboundRunTriggers.Items {

				runTriggerWorkspaceName := runTrigger.SourceableName

				// need to see if this workspace exists in the destination
				exists := doesWorkspaceExist(runTriggerWorkspaceName, destWorkspaces)

				if exists {

					// ok proven that it does exist, now to get the workspace id
					destRunTriggerWorkspaceID, err := getWorkspaceId(tfclient.GetClientContexts(), runTriggerWorkspaceName)
					if err != nil {
						return errors.Wrap(err, "Failed to get the ID of the destination Workspace that matches the Name of the Source Run Trigger Workspace: "+runTriggerWorkspaceName)
					}

					// Check if the run trigger already exists in the destination workspace
					existingRunTriggers, err := c.DestinationClient.RunTriggers.List(c.DestinationContext, destWorkSpaceID, &inboundOpts)
					if err != nil {
						return errors.Wrap(err, "Failed to list run triggers for destination workspace: "+destWorkSpaceName)
					}

					runTriggerExists := false
					for _, existingRunTrigger := range existingRunTriggers.Items {
						if existingRunTrigger.SourceableChoice.Workspace.ID == destRunTriggerWorkspaceID {
							runTriggerExists = true
							break
						}
					}

					if runTriggerExists {
						o.AddFormattedMessageUserProvided2("Run trigger already exists for workspace %v connecting from %v", destWorkSpaceName, runTriggerWorkspaceName)
						continue
					}

					// now we want to create the run trigger in the destination workspace
					runTriggerOpts := tfe.RunTriggerCreateOptions{
						Sourceable: &tfe.Workspace{ID: destRunTriggerWorkspaceID},
					}
					createdRunTrigger, err := c.DestinationClient.RunTriggers.Create(c.DestinationContext, destWorkSpaceID, runTriggerOpts)
					_ = createdRunTrigger

					if err != nil {
						return errors.Wrap(err, "Failed to create run trigger in destination workspace: "+destWorkSpaceName)
					}
					o.AddFormattedMessageCalculated2("Created run trigger for workspace %v to %v", destWorkSpaceName, runTriggerWorkspaceName)
				} else { 
					o.AddFormattedMessageCalculated("Workspace named %v does not exist in destination. Not able to setup Run Trigger", runTriggerWorkspaceName)
				}
			}
		} 
	}
	return nil
}
