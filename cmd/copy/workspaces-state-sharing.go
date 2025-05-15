// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package copy

import (
	"github.com/hashicorp-services/tfm/cmd/helper"
	"github.com/hashicorp-services/tfm/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
)

// All functions related to remote state sharing on workspaces
func copyRemoteStateSharing(c tfclient.ClientContexts, consolidateGlobal bool) error {

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

		// Check if the destination Workspace name differs from the source name
		if len(wsMapCfg) > 0 {
			destWorkSpaceName = wsMapCfg[srcWorkspace.Name]
		}

		// Ensure the destination workspace exists in the destination target
		exists := doesWorkspaceExist(destWorkSpaceName, destWorkspaces)

		// If the destination workspace exists and workspace consolidation is not enabled
		if exists {

			if srcWorkspace.GlobalRemoteState && !consolidateGlobal {

				opts := tfe.WorkspaceUpdateOptions{
					Type:              "",
					GlobalRemoteState: tfe.Bool(true),
				}

				o.AddFormattedMessageUserProvided("Setting %v workspace's remote state sharing setting to org wide in destination.", destWorkSpaceName)
				_, err := c.DestinationClient.Workspaces.Update(c.DestinationContext, c.DestinationOrganizationName, destWorkSpaceName, opts)
				if err != nil {
					return errors.Wrap(err, "failed to set workspace "+destWorkSpaceName+"to be shared globally")
				}
			}

			// If the workspace is not to be shared globally
			if !srcWorkspace.GlobalRemoteState {

				opts := tfe.WorkspaceUpdateOptions{
					Type:              "",
					GlobalRemoteState: tfe.Bool(false),
				}
				// Set workspace to not be shared globally
				o.AddFormattedMessageUserProvided("Setting %v workspace's remote state sharing setting to selected workspaces in destination.", destWorkSpaceName)

				c.DestinationClient.Workspaces.Update(c.DestinationContext, c.DestinationOrganizationName, destWorkSpaceName, opts)

				// Gather remote state consumers on the workspace
				remoteStateOpts := tfe.RemoteStateConsumersListOptions{}
				consumers, err := c.SourceClient.Workspaces.ListRemoteStateConsumers(c.SourceContext, srcWorkspace.ID, &remoteStateOpts)

				if err != nil {
					return errors.Wrap(err, "failed to list remote state consumers")
				}

				if len(consumers.Items) > 0 {
					// Get the list of workspaces in the destination
					destinationWorkspaceList, err := lookupWorkspaces(c, "destination")

					if wsMapCfg == nil {

						// Compare the workspaces in the source and destination and return a slice of matching workspaces
						workspaceMapping := CompareWorkspacesByName(consumers.Items, destinationWorkspaceList)

						o.AddFormattedMessageUserProvided2("Adding %v remote state consumers to destination workspace %v.", len(workspaceMapping), destWorkSpaceName)

						addRemoteStateConsumerOpts := tfe.WorkspaceAddRemoteStateConsumersOptions{
							Workspaces: workspaceMapping,
						}

						err = c.DestinationClient.Workspaces.AddRemoteStateConsumers(c.DestinationContext, destWorkSpaceID, addRemoteStateConsumerOpts)
						if err != nil {
							return errors.Wrap(err, "failed to add remote state consumers")
						}

					} else if len(wsMapCfg) > 0 {

						workspaceMapping := CompareWorkspacesByNameMap(consumers.Items, destinationWorkspaceList, wsMapCfg)

						o.AddFormattedMessageUserProvided2("Adding %v remote state consumers to destination workspace %v.", len(workspaceMapping), destWorkSpaceName)

						addRemoteStateConsumerOpts := tfe.WorkspaceAddRemoteStateConsumersOptions{
							Workspaces: workspaceMapping,
						}

						err = c.DestinationClient.Workspaces.AddRemoteStateConsumers(c.DestinationContext, destWorkSpaceID, addRemoteStateConsumerOpts)
						if err != nil {
							return errors.Wrap(err, "failed to add remote state consumers")
						}
					}

				} else {
					o.AddFormattedMessageCalculated("No remote state consumers to add to destination workspace %v.", destWorkSpaceName)
				}
			}

			// In a organization consolidation scenario where the source workspace is set to share remote state globally, in the destination we need to set the workspace to not share remote state globally and only share the workspace with any matching workspaces from the source organization
			if srcWorkspace.GlobalRemoteState && consolidateGlobal {

				opts := tfe.WorkspaceUpdateOptions{
					Type:              "",
					GlobalRemoteState: tfe.Bool(false),
				}

				// Set workspace to not be shared globally
				o.AddFormattedMessageUserProvided("Setting %v workspace's remote state sharing setting to selected workspaces in destination.", destWorkSpaceName)
				c.DestinationClient.Workspaces.Update(c.DestinationContext, c.DestinationOrganizationName, destWorkSpaceName, opts)

				// Gather remote state consumers on the workspace
				remoteStateOpts := tfe.RemoteStateConsumersListOptions{}
				consumers, err := c.SourceClient.Workspaces.ListRemoteStateConsumers(c.SourceContext, srcWorkspace.ID, &remoteStateOpts)
				if len(consumers.Items) > 0 {

					if err != nil {
						return errors.Wrap(err, "failed to list remote state consumers")
					}

					// Get the list of workspaces in the destination
					destinationWorkspaceList, err := lookupWorkspaces(c, "destination")
					
					if err != nil {
						return errors.Wrap(err, "failed to list Workspaces from destination")
					}

					if wsMapCfg == nil {

						// Compare the workspaces in the source and destination and return a slice of matching workspaces
						workspaceMapping := CompareWorkspacesByName(consumers.Items, destinationWorkspaceList)

						o.AddFormattedMessageUserProvided2("Adding %v remote state consumers to destination workspace %v.", len(workspaceMapping), destWorkSpaceName)

						addRemoteStateConsumerOpts := tfe.WorkspaceAddRemoteStateConsumersOptions{
							Workspaces: workspaceMapping,
						}

						err = c.DestinationClient.Workspaces.AddRemoteStateConsumers(c.DestinationContext, destWorkSpaceID, addRemoteStateConsumerOpts)
						if err != nil {
							return errors.Wrap(err, "failed to add remote state consumers")
						}
					} else if len(wsMapCfg) > 0 {

						// Compare the workspaces in the source and destination and return a slice of matching workspaces
						workspaceMapping := CompareWorkspacesByNameMap(consumers.Items, destinationWorkspaceList, wsMapCfg)

						o.AddFormattedMessageUserProvided2("Adding %v remote state consumers to destination workspace %v.", len(workspaceMapping), destWorkSpaceName)

						addRemoteStateConsumerOpts := tfe.WorkspaceAddRemoteStateConsumersOptions{
							Workspaces: workspaceMapping,
						}

						err = c.DestinationClient.Workspaces.AddRemoteStateConsumers(c.DestinationContext, destWorkSpaceID, addRemoteStateConsumerOpts)
						if err != nil {
							return errors.Wrap(err, "failed to add remote state consumers")
						}
					}
				}
			}
		} else {
			o.AddFormattedMessageCalculated("Source workspace named %v does not exist in destination. No Remote State sharing to configure\n", destWorkSpaceName)
		}
	}
	return nil

}

// function to list all workspaces in the source or destination.
func lookupWorkspaces(c tfclient.ClientContexts, side string) ([]*tfe.Workspace, error) {

	workspaces := []*tfe.Workspace{}

	opts := tfe.WorkspaceListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100},
	}

	if side == "source" {
		for {
			items, err := c.SourceClient.Workspaces.List(c.SourceContext, c.SourceOrganizationName, &opts)
			if err != nil {
				return nil, errors.Wrap(err, "Error With retrieving Workspaces from "+c.SourceHostname)
			}
			workspaces = append(workspaces, items.Items...)
			if items.CurrentPage >= items.TotalPages {
				break
			}
			opts.PageNumber = items.NextPage
		}

		return workspaces, nil
	}

	if side == "destination" {
		for {
			items, err := c.DestinationClient.Workspaces.List(c.DestinationContext, c.DestinationOrganizationName, &opts)
			if err != nil {
				return nil, errors.Wrap(err, "Error With retrieving Workspaces from "+c.DestinationHostname)
			}
			workspaces = append(workspaces, items.Items...)
			if items.CurrentPage >= items.TotalPages {
				break
			}
			opts.PageNumber = items.NextPage
		}

		return workspaces, nil
	}
	return nil, nil
}

// CompareWorkspacesByName compares two slices of workspaces and returns a slice of workspaces with matching names.
func CompareWorkspacesByName(workspaceConsumers, destinationWorkspaces []*tfe.Workspace) []*tfe.Workspace {
	nameSet := make(map[string]*tfe.Workspace, len(workspaceConsumers))
	for _, workspace := range workspaceConsumers {
		nameSet[workspace.Name] = workspace
	}

	var matchingWorkspaces []*tfe.Workspace
	for _, workspace := range destinationWorkspaces {
		if _, exists := nameSet[workspace.Name]; exists {
			matchingWorkspaces = append(matchingWorkspaces, workspace)
		}
	}

	return matchingWorkspaces
}

func CompareWorkspacesByNameMap(workspaceConsumers, destinationWorkspaces []*tfe.Workspace, wsMapCfg map[string]string) []*tfe.Workspace {

	// Build a set of destination workspace names mapped from workspaceConsumers using wsMapCfg
	mappedNames := make(map[string]struct{})
	for _, consumer := range workspaceConsumers {
		if newWorkspaceName, ok := wsMapCfg[consumer.Name]; ok {
			mappedNames[newWorkspaceName] = struct{}{}
		}
	}

	var matchingWorkspaces []*tfe.Workspace
	for _, workspace := range destinationWorkspaces {
		if _, exists := mappedNames[workspace.Name]; exists {
			matchingWorkspaces = append(matchingWorkspaces, workspace)
		}
	}

	return matchingWorkspaces
}
