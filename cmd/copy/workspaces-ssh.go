package copy

import (
	"fmt"

	"github.com/hashicorp-services/tfm/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
)

// All functions related to copying/assigning ssh provider to workspaces

// Update destination workspace with ssh-key id.
func configureSSHsettings(c tfclient.ClientContexts, org string, sshId string, ws string) (*tfe.Workspace, error) {

	workspaceSSHOptions := tfe.WorkspaceAssignSSHKeyOptions{
		Type:    "",
		SSHKeyID: &sshId,
	}

	workspace, err := c.DestinationClient.Workspaces.AssignSSHKey(c.DestinationContext, ws, workspaceSSHOptions)
	if err != nil {
		return nil, err
	}

	return workspace, nil
}

func createSSHConfiguration(c tfclient.ClientContexts, sshConfig map[string]string) error {

	fmt.Println(sshConfig)
	o.AddFormattedMessageCalculated("Found %d ssh mappings in Configuration", len(sshConfig))

	for key, element := range sshConfig {
		srcSsh := key
		destSsh := element

		// Get the source workspaces properties
		srcWorkspaces, err := discoverSrcWorkspaces(c)
		if err != nil {
			return errors.Wrap(err, "failed to list Workspaces from source while checking source VCS IDs")
		}

		// For each source workspace with a configured ssh key compare the source SSH ID to the
		// user provided SSH ID. If they match, update the matching destination workspace with
		// the user provided SSH ID that exists in the destination.
		for _, ws := range srcWorkspaces {

			if ws.SSHKey.ID != "" {
				o.AddMessageUserProvided("No SSH ID Assigned to source Workspace: ", ws.Name)
			} else {
				if ws.SSHKey.ID != "" {
					if ws.SSHKey.ID != srcSsh {
						o.AddFormattedMessageUserProvided2("Workspace %v configured SSH ID does not match provided source ID %v. Skipping.", ws.Name, srcSsh)
					} else {
						o.AddFormattedMessageUserProvided2("Updating destination workspace %v VCS Settings and OauthID %v", ws.Name, destSsh)

						configureSSHsettings(c, c.DestinationOrganizationName, destSsh, ws.ID)
					}
				} else {
					o.AddMessageUserProvided("No VCS configured to source Workspace: ", ws.Name)
				}
			}
		}
	}
	return nil
}
