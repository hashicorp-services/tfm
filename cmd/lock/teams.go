// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lock

import (
	"github.com/hashicorp-services/tfm/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// `tfm lock teams` command
	teamsLockCmd = &cobra.Command{
		Use:   "teams",
		Short: "Lock Teams",
		Long:  "Set teams - excluding owner - permissions to read in organization",
		Run: func(cmd *cobra.Command, args []string) {
			lockTeams(tfclient.GetClientContexts())
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			o.Close()
		},
	}
)

func init() {
	// Add commands
	LockCmd.AddCommand(teamsLockCmd)
}

func lockTeams(c tfclient.ClientContexts) error {

	o.AddMessageUserProvided("Locking teams in: ", c.SourceHostname)

	optsList := tfe.TeamListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100},
	}

	orgAccess := tfe.OrganizationAccessOptions{
		ManagePolicies:           tfe.Bool(false),
		ManagePolicyOverrides:    tfe.Bool(false),
		ManageWorkspaces:         tfe.Bool(false),
		ManageVCSSettings:        tfe.Bool(false),
		ManageProviders:          tfe.Bool(false),
		ManageModules:            tfe.Bool(false),
		ManageRunTasks:           tfe.Bool(false),
		ManageProjects:           tfe.Bool(false),
		ReadWorkspaces:           tfe.Bool(true),
		ReadProjects:             tfe.Bool(true),
		ManageMembership:         tfe.Bool(false),
		ManageTeams:              tfe.Bool(false),
		ManageOrganizationAccess: tfe.Bool(false),
		AccessSecretTeams:        tfe.Bool(false),
		ManageAgentPools:         tfe.Bool(false),
	}

	srcTeams, err := c.SourceClient.Teams.List(c.SourceContext, c.SourceOrganizationName, &optsList)

	if err != nil {
		return errors.Wrap(err, "failed to list teams from source while checking lock status")
	}

	for _, team := range srcTeams.Items {
		if team.Name == "owners" {
			o.AddMessageUserProvided("Skipping team: ", team.Name)
			continue
		}
		o.AddMessageUserProvided("Locking team: ", team.Name)
		_, err := c.SourceClient.Teams.Update(c.SourceContext, team.ID, tfe.TeamUpdateOptions{
			Type:               "",
			OrganizationAccess: &orgAccess,
		})

		if err != nil {
			return errors.Wrap(err, "failed to list Workspaces from source while checking lock status")
		}
		o.AddMessageUserProvided("Locked team: ", team.Name)
	}

	return nil
}
