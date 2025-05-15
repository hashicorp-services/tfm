// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package copy

import (
	"fmt"

	"strings"

	"github.com/hashicorp-services/tfm/output"
	"github.com/hashicorp-services/tfm/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	o            output.Output
	orgToProject bool

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

	// Add commands
	CopyCmd.AddCommand(teamCopyCmd)
	teamCopyCmd.Flags().BoolVarP(&orgToProject, "org-to-project", "o", false, "Migrate from organization to project access")

}

// Get all source target Teams
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

// Get all destination target Teams
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
// Used to compare source team names to the existing destination team names.
func doesTeamExist(teamName string, teams []*tfe.Team) bool {
	// Convert the teamName to lowercase (or uppercase if you prefer) for case-insensitive comparison
	teamName = strings.ToLower(teamName)

	for _, t := range teams {
		// Convert the team name in the teams slice to lowercase (or uppercase) for comparison
		if teamName == strings.ToLower(t.Name) {
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
			if orgToProject {
				fmt.Println("Migrating from organization to project access", srcteam.Name)

				// Take teams from source and create them in destination.
				// if orgToProject is true, it will
				team, err := c.DestinationClient.Teams.Create(c.DestinationContext, c.DestinationOrganizationName, tfe.TeamCreateOptions{
					Type:      "",
					Name:      &srcteam.Name,
					SSOTeamID: &srcteam.SSOTeamID,
					OrganizationAccess: &tfe.OrganizationAccessOptions{
						ManagePolicies:        tfe.Bool(false),
						ManagePolicyOverrides: tfe.Bool(false),
						ManageWorkspaces:      tfe.Bool(false),
						ManageVCSSettings:     tfe.Bool(false),
						ManageProviders:       tfe.Bool(false),
						ManageModules:         tfe.Bool(false),
						ManageRunTasks:        tfe.Bool(false),
						// release v202302-1
						ManageProjects: tfe.Bool(false),
						ReadWorkspaces: tfe.Bool(false),
						ReadProjects:   tfe.Bool(false),
						// release 202303-1
						ManageMembership: tfe.Bool(false),
					},
					Visibility: &srcteam.Visibility,
				})

				if err != nil {
					return err
				}

				o.AddDeferredMessageRead("Created team in destination organization", team.Name)
				o.AddDeferredMessageRead("New ID", team.ID)

			} else {
				fmt.Println("Migrating", srcteam.Name)
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
						// release v202302-1
						ManageProjects: &srcteam.OrganizationAccess.ManageProjects,
						ReadWorkspaces: &srcteam.OrganizationAccess.ReadWorkspaces,
						ReadProjects:   &srcteam.OrganizationAccess.ReadProjects,
						// release 202303-1
						ManageMembership: &srcteam.OrganizationAccess.ManageMembership,
					},
					Visibility: &srcteam.Visibility,
				})
				if err != nil {
					return err
				}
				o.AddDeferredMessageRead("Migrated", srcteam.Name)
			}
		}
	}
	return nil
}
