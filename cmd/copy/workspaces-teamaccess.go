package copy

import (
	"fmt"

	"github.com/hashicorp-services/tfe-mig/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
)

func discoverSrcWsTeamAccess(c tfclient.ClientContexts, wsId string) ([]*tfe.TeamAccess, error) {
	o.AddMessageUserProvided("Getting workspace team access permissions from source workspace ", wsId)
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

func discoverDestWsTeamAccess(c tfclient.ClientContexts, wsId string) ([]*tfe.TeamAccess, error) {
	o.AddMessageUserProvided("Getting workspace team access permissions from source workspace ", wsId)
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

func getSrcTeamAccessName(c tfclient.ClientContexts, srcteamId string) (string, error) {
	var srcTeamName string

	//srcTeams, err := discoverSrcTeamsNameFilter(tfclient.GetClientContexts())
	//if err != nil {
	//	fmt.Println("failed to list teams from source")
	//}

	//for _, srcteams := range srcTeams {
	t, err := c.SourceClient.Teams.Read(c.SourceContext, srcteamId)
	if err != nil {
		return "", err
	}

	srcTeamName = t.Name
	//}

	return srcTeamName, nil
}

func discoverSrcTeamsNameFilter(c tfclient.ClientContexts) ([]*tfe.Team, error) {
	o.AddMessageUserProvided("Getting list of teams from: ", c.SourceHostname)
	srcTeams := []*tfe.Team{}

	opts := tfe.TeamListOptions{
		ListOptions: tfe.ListOptions{PageNumber: 1, PageSize: 100},
		Include:     []tfe.TeamIncludeOpt{},
	}
	for {
		items, err := c.SourceClient.Teams.List(c.SourceContext, c.SourceOrganizationName, &opts)
		if err != nil {
			return nil, err
		}

		srcTeams = append(srcTeams, items.Items...)

		o.AddFormattedMessageCalculated("Found %d Teams", len(srcTeams))

		if items.CurrentPage >= items.TotalPages {
			break
		}
		opts.PageNumber = items.NextPage

	}

	return srcTeams, nil
}

func discoverDestTeamsNameFilter(c tfclient.ClientContexts, wsname string) ([]*tfe.Team, error) {
	o.AddMessageUserProvided("Getting list of teams from: ", c.DestinationHostname)
	destTeams := []*tfe.Team{}

	opts := tfe.TeamListOptions{
		ListOptions: tfe.ListOptions{PageNumber: 1, PageSize: 100},
		Include:     []tfe.TeamIncludeOpt{},
		Names:       []string{wsname},
	}
	for {
		items, err := c.DestinationClient.Teams.List(c.DestinationContext, c.DestinationOrganizationName, &opts)
		if err != nil {
			return nil, err
		}

		destTeams = append(destTeams, items.Items...)

		o.AddFormattedMessageCalculated("Found %d Teams", len(destTeams))

		if items.CurrentPage >= items.TotalPages {
			break
		}
		opts.PageNumber = items.NextPage

	}

	return destTeams, nil
}

func getDestTeamAccessNameAndID(c tfclient.ClientContexts, wsname string) (string, string, error) {
	var destTeamName string
	var destTeamId string

	destTeams, err := discoverDestTeamsNameFilter(tfclient.GetClientContexts(), wsname)
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

func doesTeamNameMatch(srcTeamName string, destTeamName string) bool {
	return srcTeamName == destTeamName
}

func doesTeamAccessPermissionsExist(c tfclient.ClientContexts, teamName string, destWorkspaceId string) (bool, error) {
	// Get the source team access permissions for the destination workspace.
	destTeamAccess, err := discoverDestWsTeamAccess(tfclient.GetClientContexts(), destWorkspaceId)
	if err != nil {
		return true, errors.Wrap(err, "failed to list Team Access for dest workspace")
	}

	var destTeamName string

	for _, i := range destTeamAccess {
		t, err := c.DestinationClient.Teams.Read(c.DestinationContext, i.Team.ID)
		if err != nil {
			return true, err
		}
		destTeamName = t.Name
	}

	fmt.Println("dest team access name from array is", destTeamName)

	if teamName == string(destTeamName) {
		return true, err
	}
	return false, err
}

func copyWsTeamAccess(c tfclient.ClientContexts) error {
	// Get the source workspaces properties
	srcWorkspaces, err := discoverSrcWorkspaces(tfclient.GetClientContexts())
	if err != nil {
		return errors.Wrap(err, "failed to list Workspaces from source")
	}

	// Get the destination workspace properties
	destWorkspaces, err := discoverDestWorkspaces(tfclient.GetClientContexts())
	if err != nil {
		return errors.Wrap(err, "failed to list Workspaces from source")
	}

	// For each srcworkspace check to see if a workspace with the same name exists in the destination
	for _, srcworkspace := range srcWorkspaces {

		if !doesWorkspaceExist(srcworkspace.Name, destWorkspaces) {
			return errors.New("Workspace not found")
		}

		// If Destination contians workspace with the same name as the source one, get the destination workspace ID
		destWorkspaceId, err := getWorkspaceId(tfclient.GetClientContexts(), srcworkspace.Name)
		if err != nil {
			return errors.Wrap(err, "Failed to get the ID of the destination Workspace that matches the Name of the Source Workspace")
		}

		// Get the source team access permissions for the source workspace
		srcTeamAccess, err := discoverSrcWsTeamAccess(tfclient.GetClientContexts(), srcworkspace.ID)
		if err != nil {
			return errors.Wrap(err, "failed to list Team Access for source workspace")
		}

		// If The source team access permissions contians teams, get the source team names
		for _, srcteamname := range srcTeamAccess {
			if len(srcTeamAccess) > 0 {
				srcTeamNames, err := getSrcTeamAccessName(tfclient.GetClientContexts(), srcteamname.Team.ID)
				if err != nil {
					return errors.Wrap(err, "failed to find source team name")
				}

				fmt.Println("srcTeamNames is ", srcTeamNames)

				// Get the matching destination team names and their IDs
				fmt.Println("team", srcTeamNames)
				destTeamNames, destTeamId, err := getDestTeamAccessNameAndID(tfclient.GetClientContexts(), srcTeamNames)
				if err != nil {
					return errors.Wrap(err, "failed to find destination team name")
				}

				fmt.Println("DestTeamNames is ", destTeamNames, "for SrcTeamName", srcteamname)

				// Ensure the team names match between the source and destination
				match := doesTeamNameMatch(srcTeamNames, destTeamNames)
				fmt.Println(match)
				if match == true {
					if err != nil {
						return errors.Wrap(err, "failed to find destination team ID")
					}
					println("destTeamId is", destTeamId)
					// Need to make a conditional for if "custom" permissions are set
					// Loop through the team access for the source workspace
					for _, teamaccess := range srcTeamAccess {
						// For each team access setting, check for an existing access setting in the destination
						exists, err := doesTeamAccessPermissionsExist(tfclient.GetClientContexts(), srcTeamNames, destWorkspaceId)
						if err != nil {
							return errors.Wrap(err, "failed to get destination permissions")
						}
						println("exists =", exists)
						if exists {
							fmt.Println("Exists in destination will not migrate", teamaccess.Team.Name)
						} else {
							teamaccess, err := c.DestinationClient.TeamAccess.Add(c.DestinationContext, tfe.TeamAccessAddOptions{
								Type:   "",
								Access: &teamaccess.Access,
								//Runs:             &teamaccess.Runs,
								//Variables:        &teamaccess.Variables,
								//StateVersions:    &teamaccess.StateVersions,
								//SentinelMocks:    &teamaccess.SentinelMocks,
								//WorkspaceLocking: &teamaccess.WorkspaceLocking,
								//RunTasks:         &teamaccess.RunTasks,
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
							println("test2")
							o.AddDeferredMessageRead("Migrated Team Access ", teamaccess.ID)
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
