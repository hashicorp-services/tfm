package copy

import (
	"fmt"

	"github.com/hashicorp-services/tfm/cmd/helper"
	"github.com/hashicorp-services/tfm/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
)

// Get source workspace team access permissions.
func discoverSrcWsTeamAccess(c tfclient.ClientContexts, wsId string, wsName string) ([]*tfe.TeamAccess, error) {
	o.AddMessageUserProvided("Getting workspace team access permissions from source workspace ", wsName)
	srcTeamAccess := []*tfe.TeamAccess{}

	opts := tfe.TeamAccessListOptions{
		ListOptions: tfe.ListOptions{},
		WorkspaceID: wsId,
	}
	for {
		items, err := c.SourceClient.TeamAccess.List(c.SourceContext, &opts)
		if err != nil {
			return nil, err
		}

		srcTeamAccess = append(srcTeamAccess, items.Items...)

		o.AddFormattedMessageCalculated("Found %d sets of Workspace Team Access Permissions", len(srcTeamAccess))

		if items.CurrentPage >= items.TotalPages {
			break
		}
		opts.PageNumber = items.NextPage

	}

	return srcTeamAccess, nil
}

// Get destination workspace team access permissions.
func discoverDestWsTeamAccess(c tfclient.ClientContexts, wsId string, wsName string) ([]*tfe.TeamAccess, error) {
	o.AddMessageUserProvided("Getting workspace team access permissions from destination workspace ", wsName)
	destTeamAccess := []*tfe.TeamAccess{}

	opts := tfe.TeamAccessListOptions{
		ListOptions: tfe.ListOptions{},
		WorkspaceID: wsId,
	}
	for {
		items, err := c.DestinationClient.TeamAccess.List(c.DestinationContext, &opts)
		if err != nil {
			return nil, err
		}

		destTeamAccess = append(destTeamAccess, items.Items...)

		o.AddFormattedMessageCalculated("Found %d sets of Workspace Team Access Permissions", len(destTeamAccess))

		if items.CurrentPage >= items.TotalPages {
			break
		}
		opts.PageNumber = items.NextPage

	}

	return destTeamAccess, nil
}

// Get a source teams name based on the team ID taken from source workspace team access permissions
func getSrcTeamAccessName(c tfclient.ClientContexts, srcteamId string) (string, error) {
	var srcTeamName string

	t, err := c.SourceClient.Teams.Read(c.SourceContext, srcteamId)
	if err != nil {
		return "", err
	}

	srcTeamName = t.Name

	return srcTeamName, nil
}

// Get the properties of a specigic destination team filtering by name.
func discoverDestTeamsNameFilter(c tfclient.ClientContexts, teamName string) ([]*tfe.Team, error) {
	destTeams := []*tfe.Team{}

	opts := tfe.TeamListOptions{
		ListOptions: tfe.ListOptions{PageNumber: 1, PageSize: 100},
		Include:     []tfe.TeamIncludeOpt{},
		Names:       []string{teamName},
	}
	for {
		items, err := c.DestinationClient.Teams.List(c.DestinationContext, c.DestinationOrganizationName, &opts)
		if err != nil {
			return nil, err
		}

		destTeams = append(destTeams, items.Items...)

		if items.CurrentPage >= items.TotalPages {
			break
		}
		opts.PageNumber = items.NextPage

	}

	return destTeams, nil
}

// Get a specific destination teams Name and ID by filtering by team name.
func getDestTeamAccessNameAndID(c tfclient.ClientContexts, teamName string) (string, string, error) {
	var destTeamName string
	var destTeamId string

	destTeams, err := discoverDestTeamsNameFilter(tfclient.GetClientContexts(), teamName)
	if err != nil {
		fmt.Println("failed to list teams from destination")
	}

	for _, destteams := range destTeams {
		t, err := c.DestinationClient.Teams.Read(c.DestinationContext, destteams.ID)
		if err != nil {
			return "", "", err
		}

		destTeamName = t.Name
		destTeamId = t.ID
	}

	return destTeamName, destTeamId, nil
}

// Check to see if the source and destination team name match.
func doesTeamNameMatch(srcTeamName string, destTeamName string) bool {
	return srcTeamName == destTeamName
}

// Check the destination workspace team access permissions for existing permissions.
func doesTeamAccessPermissionsExist(c tfclient.ClientContexts, teamName string, destTeamAccess []*tfe.TeamAccess, destWorkspaceId string) (bool, error) {
	var destTeamName string

	for _, i := range destTeamAccess {
		t, err := c.DestinationClient.Teams.Read(c.DestinationContext, i.Team.ID)
		if err != nil {
			return true, err
		}
		destTeamName = t.Name

		if teamName == string(destTeamName) {
			return true, err
		}
	}
	return false, nil

}

// Check workspace permissions for custom access type.
func checkCustom(c tfclient.ClientContexts, srcteramaccess *tfe.TeamAccess) bool {
	if srcteramaccess.Access == "custom" {
		return true
	}
	return false
}

// Default workspace access permissions creation. Seperate functions required for custom and default permission creation.
func createTeamAccess(c tfclient.ClientContexts, srcTeamName string, destTeamId string, destWorkspaceId string, srcworkspace *tfe.Workspace, srcteam *tfe.TeamAccess) error {
	o.AddMessageUserProvided("Migrating team access permissions for: ", srcTeamName)

	teamaccess, err := c.DestinationClient.TeamAccess.Add(c.DestinationContext, tfe.TeamAccessAddOptions{
		Type:   "",
		Access: &srcteam.Access,
		Team: &tfe.Team{
			ID: destTeamId,
		},
		Workspace: &tfe.Workspace{
			ID: destWorkspaceId,
		},
	})

	if err != nil {
		return err
	}

	_ = teamaccess

	return nil
}

// Custom workspace access permissions. These can only be edited when Access is 'custom'; otherwise, they are
// read-only and reflect the Access level's implicit permissions.
func createCustomTeamAccess(c tfclient.ClientContexts, srcTeamName string, destTeamId string, destWorkspaceId string, srcworkspace *tfe.Workspace, srcteam *tfe.TeamAccess) error {
	o.AddMessageUserProvided("Migrating team access permissions for: ", srcTeamName)
	teamaccess, err := c.DestinationClient.TeamAccess.Add(c.DestinationContext, tfe.TeamAccessAddOptions{
		Type:             "",
		Access:           &srcteam.Access,
		Runs:             &srcteam.Runs,
		Variables:        &srcteam.Variables,
		StateVersions:    &srcteam.StateVersions,
		SentinelMocks:    &srcteam.SentinelMocks,
		WorkspaceLocking: &srcteam.WorkspaceLocking,
		RunTasks:         &srcteam.RunTasks,
		Team: &tfe.Team{
			ID: destTeamId,
		},
		Workspace: &tfe.Workspace{
			ID: destWorkspaceId,
		},
	})

	if err != nil {
		return err
	}

	_ = teamaccess

	return nil
}

func copyWsTeamAccess(c tfclient.ClientContexts) error {
	// Get the source workspaces properties
	srcWorkspaces, err := getSrcWorkspacesCfg(c)
	if err != nil {
		return errors.Wrap(err, "failed to list Workspaces from source")
	}

	// Get/Check if Workspace map exists
	wsMapCfg, err := helper.ViperStringSliceMap("workspace-map")
	if err != nil {
		fmt.Println("invalid input for workspace-map")
	}

	// Get the destination workspace properties
	destWorkspaces, err := discoverDestWorkspaces(tfclient.GetClientContexts())
	if err != nil {
		return errors.Wrap(err, "failed to list Workspaces from source")
	}

	// For each srcworkspace check to see if a workspace with the same name exists in the destination
	for _, srcworkspace := range srcWorkspaces {
		destWorkSpaceName := srcworkspace.Name

		// Check if Destination Workspace Name to be Change
		if len(wsMapCfg) > 0 {
			destWorkSpaceName = wsMapCfg[srcworkspace.Name]
		}

		if !doesWorkspaceExist(destWorkSpaceName, destWorkspaces) {
			return errors.New("Workspace not found")
		}

		// Get the destination workspace ID
		destWorkspaceId, err := getWorkspaceId(tfclient.GetClientContexts(), destWorkSpaceName)
		if err != nil {
			return errors.Wrap(err, "Failed to get the ID of the destination Workspace that matches the Name of the Source Workspace")
		}

		// Get the source team access permissions for the source workspace
		srcTeamAccess, err := discoverSrcWsTeamAccess(tfclient.GetClientContexts(), srcworkspace.ID, srcworkspace.Name)
		if err != nil {
			return errors.Wrap(err, "failed to list Team Access for source workspace")
		}

		destTeamAccess, err := discoverDestWsTeamAccess(tfclient.GetClientContexts(), destWorkspaceId, destWorkSpaceName)
		if err != nil {
			return errors.Wrap(err, "failed to list Team Access for dest workspace")
		}

		// If The source team access permissions contians teams, get the source team names filtering by Team ID
		for _, srcteam := range srcTeamAccess {
			if len(srcTeamAccess) > 0 {
				srcTeamName, err := getSrcTeamAccessName(tfclient.GetClientContexts(), srcteam.Team.ID)
				if err != nil {
					return errors.Wrap(err, "failed to find source team name")
				}

				// Get the matching destination team names and their IDs
				destTeamNames, destTeamId, err := getDestTeamAccessNameAndID(tfclient.GetClientContexts(), srcTeamName)
				if err != nil {
					return errors.Wrap(err, "failed to find destination team name")
				}

				// Ensure the team names match between the source and destination
				match := doesTeamNameMatch(srcTeamName, destTeamNames)
				if match == true {
					// Loop through the team access for the source workspace
					// For each team access setting, check for an existing access setting in the destination
					exists, err := doesTeamAccessPermissionsExist(tfclient.GetClientContexts(), srcTeamName, destTeamAccess, destWorkspaceId)
					if err != nil {
						return errors.Wrap(err, "failed to get destination permissions")
					}
					if exists {
						o.AddMessageUserProvided("Team access exists in destination workspace, skipping migration for: ", srcTeamName)
					} else {
						custom := checkCustom(c, srcteam)
						if custom {
							createCustomTeamAccess(c, srcTeamName, destTeamId, destWorkspaceId, srcworkspace, srcteam)
						} else {
							createTeamAccess(c, srcTeamName, destTeamId, destWorkspaceId, srcworkspace, srcteam)
						}
					}
				} else {
					fmt.Println("Destination Team ID required to migrate Team Access, but none found")
				}

				if err != nil {
					return errors.Wrap(err, "failed to find destination team name and ID")
				}
			} else {
				fmt.Println("No team access permissions found on source workspace")
			}
		}
	}
	return nil

}
