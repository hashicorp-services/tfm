package copy

import (
	"fmt"

	"github.com/hashicorp-services/tfe-mig/output"
	"github.com/hashicorp-services/tfe-mig/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	o output.Output

	// `tfemig copy teams` command
	teamCopyCmd = &cobra.Command{
		Use:   "teams",
		Short: "Copy Teams",
		Long:  "Copy Teams from source to destination org",
		RunE: func(cmd *cobra.Command, args []string) error {
			return copyTeams(
				tfclient.GetClientContexts())

		},
		PostRun: func(cmd *cobra.Command, args []string) {
			o.Close()
		},
	}
)

func init() {
	// Flags().StringP, etc... - the "P" gives us the option for a short hand

	// `tfe-migrate copy teams all` command
	teamCopyCmd.Flags().BoolP("all", "a", false, "List all? (optional)")

	// Add commands
	CopyCmd.AddCommand(teamCopyCmd)

}

func discoverSrcTeams(c tfclient.ClientContexts) ([]*tfe.Team, error) {
	o.AddMessageUserProvided("Getting list of teams from: ", c.SourceHostname)
	srcTeams := []*tfe.Team{}

	opts := tfe.TeamListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100},
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

func discoverDestTeams(c tfclient.ClientContexts) ([]*tfe.Team, error) {
	o.AddMessageUserProvided("Getting list of teams from: ", c.DestinationHostname)
	destTeams := []*tfe.Team{}

	opts := tfe.TeamListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100},
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

// Takes a team name and a slice of teams as type []*tfe.Team and
// returns true if the team name exists within the provided slice of teams.
// Used to compare source team names to the destination team names.
func doesTeamExist(teamName string, teams []*tfe.Team) bool {
	for _, t := range teams {
		if teamName == t.Name {
			return true
		}
	}
	return false
}

// Gets all source team names and all destination team names and recreates
// the source teams in the destination if the team name does not exist in the destination.
func copyTeams(c tfclient.ClientContexts) error {
	// Get the source teams properties
	srcTeams, err := discoverSrcTeams(tfclient.GetClientContexts())
	if err != nil {
		return errors.Wrap(err, "failed to list teams from source")
	}

	// Get the destination teams properties
	destTeams, err := discoverDestTeams(tfclient.GetClientContexts())
	if err != nil {
		return errors.Wrap(err, "failed to list teams from destination")
	}

	// Loop each team in the srcTeams slice, check for the team existence in the destination,
	// and if a team exists in the destination, then do nothing, else create team in destination.
	for _, srcteam := range srcTeams {
		exists := doesTeamExist(srcteam.Name, destTeams)
		if exists {
			fmt.Println("Exists in destination will not migrate", srcteam.Name)
		} else {
			srcteam, err := c.DestinationClient.Teams.Create(c.DestinationContext, c.DestinationOrganizationName, tfe.TeamCreateOptions{
				Type:      "",
				Name:      &srcteam.Name,
				SSOTeamID: &srcteam.SSOTeamID,
				OrganizationAccess: &tfe.OrganizationAccessOptions{
					ManagePolicies:        &srcteam.OrganizationAccess.ManagePolicies,
					ManagePolicyOverrides: &srcteam.OrganizationAccess.ManagePolicyOverrides,
					ManageWorkspaces:      &srcteam.OrganizationAccess.ManageWorkspaces,
					ManageVCSSettings:     &srcteam.OrganizationAccess.ManageVCSSettings,
					ManageProviders:       &srcteam.OrganizationAccess.ManageProviders,
					ManageModules:         &srcteam.OrganizationAccess.ManageModules,
					ManageRunTasks:        &srcteam.OrganizationAccess.ManageRunTasks,
				},
				Visibility: &srcteam.Visibility,
			})
			if err != nil {
				return err
			}
			o.AddDeferredMessageRead("Migrated", srcteam.Name)
		}
	}
	return nil
}
