package copy

import (
	"fmt"

	"github.com/hashicorp-services/tfe-mig/cmd/helper"
	"github.com/hashicorp-services/tfe-mig/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (

	// `tfemigrate copy workspaces` command
	workspacesCopyCmd = &cobra.Command{
		Use:     "workspaces",
		Short:   "Copy Workspaces",
		Aliases: []string{"ws"},
		Long:    "Copy Workspaces from source to destination org",
		RunE: func(cmd *cobra.Command, args []string) error {
			return copyWorkspaces(
				tfclient.GetClientContexts())

		},
		PostRun: func(cmd *cobra.Command, args []string) {
			o.Close()
		},
	}
)

func init() {
	// `tfemigrate copy workspaces --all` command
	workspacesCopyCmd.Flags().BoolP("all", "a", false, "Copy all Workspaces found in configuration file")

	// `tfemigrate copy workspaces --workspace-id [WORKSPACEID]`
	workspacesCopyCmd.Flags().String("workspace-id", "", "Specify one single workspace ID to copy to destination")
	workspacesCopyCmd.Flags().BoolP("vars", "", false, "Copy workspace variables")
	workspacesCopyCmd.Flags().BoolP("state", "",false,  "Copy workspace states")

	// Add commands
	CopyCmd.AddCommand(workspacesCopyCmd)

}

// Gets all source workspaces and ensure destination workspaces exist and recreates
// the workspace in the destination if the workspace does not exist in the destination.
func copyWorkspaces(c tfclient.ClientContexts) error {

	// Get list of workspaces from configuration file
	// This function will only work with a configuration file as we expect the migration to be automated in a pipeline
	// thus repeatable as migration of workspaces occur.

	// Check List of Workspaces from Config
	srcWorkspaces := viper.GetStringSlice("workspaces")

	o.AddFormattedMessageCalculated("Found %d Workspaces in Configuration", len(srcWorkspaces))

	fmt.Println("\n\nworkspaces from config are: ", srcWorkspaces)
	fmt.Printf("\n\nworkspaces type are: %T", srcWorkspaces)

	for i, s := range srcWorkspaces {
		fmt.Println("Workspace ", i, ":", s)
	}

	// Check Workspaces exist in source
	workspaceExists(tfclient.GetClientContexts())

	// // Get the source teams properties
	// srcWorkspaces, err := discoverSrcTeams(tfclient.GetClientContexts())
	// if err != nil {
	// 	return errors.Wrap(err, "failed to list teams from source")
	// }

	// // Get the destination teams properties
	// destWorkspaces, err := discoverDestTeams(tfclient.GetClientContexts())
	// if err != nil {
	// 	return errors.Wrap(err, "failed to list teams from destination")
	// }

	// // Loop each team in the srcTeams slice, check for the team existence in the destination,
	// // and if a team exists in the destination, then do nothing, else create team in destination.
	// for _, srcteam := range srcTeams {
	// 	exists := doesTeamExist(srcteam.Name, destTeams)
	// 	if exists {
	// 		fmt.Println("Exists in destination will not migrate", srcteam.Name)
	// 	} else {
	// 		srcteam, err := c.DestinationClient.Teams.Create(c.DestinationContext, c.DestinationOrganizationName, tfe.TeamCreateOptions{
	// 			Type:      "",
	// 			Name:      &srcteam.Name,
	// 			SSOTeamID: &srcteam.SSOTeamID,
	// 			OrganizationAccess: &tfe.OrganizationAccessOptions{
	// 				ManagePolicies:        &srcteam.OrganizationAccess.ManagePolicies,
	// 				ManagePolicyOverrides: &srcteam.OrganizationAccess.ManagePolicyOverrides,
	// 				ManageWorkspaces:      &srcteam.OrganizationAccess.ManageWorkspaces,
	// 				ManageVCSSettings:     &srcteam.OrganizationAccess.ManageVCSSettings,
	// 				ManageProviders:       &srcteam.OrganizationAccess.ManageProviders,
	// 				ManageModules:         &srcteam.OrganizationAccess.ManageModules,
	// 				ManageRunTasks:        &srcteam.OrganizationAccess.ManageRunTasks,
	// 			},
	// 			Visibility: &srcteam.Visibility,
	// 		})
	// 		if err != nil {
	// 			return err
	// 		}
	// 		o.AddDeferredMessageRead("Migrated", srcteam.Name)
	// 	}
	// }
	return nil
}

func workspaceExists(c tfclient.ClientContexts) error {
	allItems := []*tfe.Workspace{}

	opts := tfe.WorkspaceListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100},
	}

	for {
		items, err := c.SourceClient.Workspaces.List(c.SourceContext, c.SourceOrganizationName, &opts)
		if err != nil {
			helper.LogError(err, "failed to list workspace")
		}

		allItems = append(allItems, items.Items...)

		o.AddFormattedMessageCalculated("Found %d Workspaces", len(allItems))

		if items.CurrentPage >= items.TotalPages {
			break
		}
		opts.PageNumber = items.NextPage
	}
	o.AddTableHeaders("Name")
	for _, i := range allItems {

		o.AddTableRows(i.Name)
	}

	return nil

}
