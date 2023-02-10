package copy

import (
	"fmt"

	"github.com/hashicorp-services/tfm/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
)

// All functions related to copying/assigning vcs provider to workspaces

// Update workspace execution mode to agent and assign an agent pool ID to a workspace.
// func configureVCSsettings(c tfclient.ClientContexts, org string, vcsSettings *tfe.VCSRepoOptions, ws string) (*tfe.Workspace, error) {
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

func createVCSConfiguration(c tfclient.ClientContexts, vcsConfig map[string]string) error {

	fmt.Println(vcsConfig)
	o.AddFormattedMessageCalculated("Found %d VCS mappings in Configuration", len(vcsConfig))

	for key, element := range vcsConfig {
		srcvcs := key
		destvcs := element

		// Get the source workspaces properties
		srcWorkspaces, err := getSrcWorkspacesCfg(c)
		if err != nil {
			return errors.Wrap(err, "failed to list Workspaces from source while checking source VCS IDs")
		}

		// For each source workspace with an execution mode of "agent", compare the source agent pool ID to the
		// user provided source pool ID. If they match, update the matching destination workspace with
		// the user provided agent pool ID that exists in the destination.
		for _, ws := range srcWorkspaces {

			if ws.VCSRepo == nil {
				o.AddMessageUserProvided("No VCS ID Assigned to source Workspace: ", ws.Name)
			} else {
					if ws.VCSRepo.OAuthTokenID != srcvcs {
						o.AddFormattedMessageUserProvided2("Workspace %v configured VCS ID does not match provided source ID %v. Skipping.", ws.Name, srcvcs)
					} else {
						o.AddFormattedMessageUserProvided2("Updating destination workspace %v VCS Settings and OauthID %v", ws.Name, destvcs)

						vcsConfig := tfe.VCSRepoOptions{
							Branch:            &ws.VCSRepo.Branch,
							Identifier:        &ws.VCSRepo.Identifier,
							IngressSubmodules: &ws.VCSRepo.IngressSubmodules,
							OAuthTokenID:      &destvcs,
							TagsRegex:         &ws.VCSRepo.TagsRegex,
						}

						configureVCSsettings(c, c.DestinationOrganizationName, vcsConfig, ws.Name)
					}
				} 
			}
		}
	return nil
}

