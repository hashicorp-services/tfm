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

// All functions related to copying/assigning vcs provider(s) to workspaces

// Update the destination workspace VCS setttings
func configureVCSsettings(c tfclient.ClientContexts, org string, vcsOptions tfe.VCSRepoOptions, ws string) (*tfe.Workspace, error) {

	workspaceOptions := tfe.WorkspaceUpdateOptions{
		Type:    "",
		VCSRepo: &vcsOptions,
	}

	workspace, err := c.DestinationClient.Workspaces.Update(c.DestinationContext, c.DestinationOrganizationName, ws, workspaceOptions)
	if err != nil {
		return nil, err
	}

	return workspace, nil
}

// Main function for --vcs flag
func createVCSConfiguration(c tfclient.ClientContexts, vcsConfig map[string]string) error {

	// for each `source-ot-ID=dest-ot-ID` string in the map, define the source oauth-ID and the target oauth-ID
	for key, element := range vcsConfig {
		srcvcs := key
		destvcs := element

		// Get the source workspaces properties
		srcWorkspaces, err := getSrcWorkspacesCfg(c)
		if err != nil {
			return errors.Wrap(err, "Failed to list Workspaces from source while checking source VCS IDs")
		}

		// Get/Check if Workspace map exists
		wsMapCfg, err := helper.ViperStringSliceMap("workspaces-map")
		if err != nil {
			fmt.Println("Invalid input for workspaces-map")
		}

		// For each source workspace with a VCS connection, compare the source oauth ID to the
		// user provided oauth ID. If they match, update the destination workspace with
		// the user provided oauth ID that exists in the destination.
		for _, ws := range srcWorkspaces {
			destWorkSpaceName := ws.Name

			// Check if the destination Workspace name differs from the source name
			if len(wsMapCfg) > 0 {
				destWorkSpaceName = wsMapCfg[ws.Name]
			}

			// If the source workspace has no VCS assigned, do nothing and inform the user
			if ws.VCSRepo == nil {
				o.AddMessageUserProvided("No VCS ID Assigned to source Workspace: ", ws.Name)
			} else {

				// If the source Workspace assigned VCS does not match the one provided by the user on the left side of the `vcs-map`, do nothing and inform the user
				if ws.VCSRepo.OAuthTokenID != srcvcs {
					o.AddFormattedMessageUserProvided2("Workspace %v configured VCS ID does not match provided source ID %v. Skipping.", ws.Name, srcvcs)

					// If the source Workspace assigned VCS matches the one provided by the user on the left side of the `vcs-map`, update the destination Workspace
					// with the VCS provided by the user on the right side of the `vcs-map`
				} else {
					o.AddFormattedMessageUserProvided2("Updating destination Workspace %v VCS Settings and OauthID %v", destWorkSpaceName, destvcs)

					vcsConfig := tfe.VCSRepoOptions{
						Branch:            &ws.VCSRepo.Branch,
						Identifier:        &ws.VCSRepo.Identifier,
						IngressSubmodules: &ws.VCSRepo.IngressSubmodules,
						OAuthTokenID:      &destvcs,
						TagsRegex:         &ws.VCSRepo.TagsRegex,
					}

					configureVCSsettings(c, c.DestinationOrganizationName, vcsConfig, destWorkSpaceName)
				}
			}
		}
	}
	return nil
}
