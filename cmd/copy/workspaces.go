package copy

import (
	"fmt"

	"github.com/hashicorp-services/tfm/cmd/helper"
	"github.com/hashicorp-services/tfm/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	state bool
	//vars  bool
	teamaccess bool

	// `tfemigrate copy workspaces` command
	workspacesCopyCmd = &cobra.Command{
		Use:     "workspaces",
		Short:   "Copy Workspaces",
		Aliases: []string{"ws"},
		Long:    "Copy Workspaces from source to destination org",
		//ValidArgs: []string{"state", "vars"},
		//Args:      cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch {
			case state:
				return copyStates(tfclient.GetClientContexts())
			//case vars:
			//return variableCopy(tfclient.GetClientContexts())
			//}
			case teamaccess:
				return copyWsTeamAccess(tfclient.GetClientContexts())
			}
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
	workspacesCopyCmd.Flags().BoolVarP(&state, "state", "s", false, "Copy workspace states")
	//workspacesCopyCmd.Flags().BoolVarP(&vars, "vars", "v", false, "Copy workspace vars")
	workspacesCopyCmd.Flags().BoolVarP(&teamaccess, "teamaccess", "t", false, "Copy workspace Team Access")

	// Add commands
	CopyCmd.AddCommand(workspacesCopyCmd)

}

func discoverSrcWorkspaces(c tfclient.ClientContexts) ([]*tfe.Workspace, error) {
	o.AddMessageUserProvided("Getting list of workspaces from: ", c.SourceHostname)
	srcWorkspaces := []*tfe.Workspace{}

	opts := tfe.WorkspaceListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100},
	}
	for {
		items, err := c.SourceClient.Workspaces.List(c.SourceContext, c.SourceOrganizationName, &opts)
		if err != nil {
			return nil, err
		}

		srcWorkspaces = append(srcWorkspaces, items.Items...)

		o.AddFormattedMessageCalculated("Found %d Workspaces", len(srcWorkspaces))

		if items.CurrentPage >= items.TotalPages {
			break
		}
		opts.PageNumber = items.NextPage

	}

	return srcWorkspaces, nil
}

func discoverDestWorkspaces(c tfclient.ClientContexts) ([]*tfe.Workspace, error) {
	o.AddMessageUserProvided("Getting list of workspaces from: ", c.DestinationHostname)
	destWorkspaces := []*tfe.Workspace{}

	opts := tfe.WorkspaceListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100},
	}
	for {
		items, err := c.DestinationClient.Workspaces.List(c.DestinationContext, c.DestinationOrganizationName, &opts)
		if err != nil {
			return nil, err
		}

		destWorkspaces = append(destWorkspaces, items.Items...)

		o.AddFormattedMessageCalculated("Found %d Workspaces", len(destWorkspaces))

		if items.CurrentPage >= items.TotalPages {
			break
		}
		opts.PageNumber = items.NextPage

	}

	return destWorkspaces, nil
}

// Takes a team name and a slice of teams as type []*tfe.Team and
// returns true if the team name exists within the provided slice of teams.
// Used to compare source team names to the destination team names.
func doesWorkspaceExist(workspaceName string, ws []*tfe.Workspace) bool {
	for _, w := range ws {
		if workspaceName == w.Name {
			return true
		}
	}
	return false
}

// Gets all source workspaces and ensure destination workspaces exist and recreates
// the workspace in the destination if the workspace does not exist in the destination.
func copyWorkspaces(c tfclient.ClientContexts) error {

	// Get list of workspaces from configuration file
	// This function will only work with a configuration file as we expect the migration to be automated in a pipeline
	// thus repeatable as migration of workspaces occur.

	// Get Workspaces from Config
	srcWorkspacesCfg := viper.GetStringSlice("workspaces")

	o.AddFormattedMessageCalculated("Found %d Workspaces in Configuration", len(srcWorkspacesCfg))

	// Get the source workspaces properties
	srcWorkspaces, err := discoverSrcWorkspaces(tfclient.GetClientContexts())
	if err != nil {
		return errors.Wrap(err, "failed to list teams from source")
	}

	// Get the destination Workspace properties
	destWorkspaces, err := discoverDestWorkspaces(tfclient.GetClientContexts())
	if err != nil {
		return errors.Wrap(err, "failed to list teams from destination")
	}

	// Check Workspaces exist in source from config
	for _, s := range srcWorkspacesCfg {
		fmt.Println("\nFound Workspace in config:", s, " exists in", viper.GetString("sourceHostname"))
		exists := doesWorkspaceExist(s, srcWorkspaces)
		if !exists {
			fmt.Printf("Defined Workspace in Config %s does not exist in %s", s, viper.GetString("sourceHostname"))
			break
		}
	}

	// Loop each team in the srcWorkspaces slice, check for the workspace existence in the destination,
	// and if a workspace exists in the destination, then do nothing, else create workspace in destination.
	// Most values will be
	for _, srcworkspace := range srcWorkspaces {
		exists := doesWorkspaceExist(srcworkspace.Name, destWorkspaces)

		// Copy tags over
		var tag []*tfe.Tag

		for _, t := range srcworkspace.TagNames {
			tag = append(tag, &tfe.Tag{Name: t})
		}

		if exists {
			fmt.Println("Exists in destination will not migrate", srcworkspace.Name)
		} else {
			srcworkspace, err := c.DestinationClient.Workspaces.Create(c.DestinationContext, c.DestinationOrganizationName, tfe.WorkspaceCreateOptions{
				Type: "",
				// AgentPoolID:        new(string),
				AllowDestroyPlan: &srcworkspace.AllowDestroyPlan,
				// AssessmentsEnabled: new(bool),
				// AutoApply:          new(bool),
				Description:   &srcworkspace.Description,
				ExecutionMode: &srcworkspace.ExecutionMode,
				// FileTriggersEnabled:        new(bool),
				// GlobalRemoteState:          new(bool),
				// MigrationEnvironment:       new(string),
				Name: &srcworkspace.Name,
				// QueueAllRuns:               new(bool),
				// SpeculativeEnabled:         new(bool),
				// SourceName:                 new(string),
				// SourceURL:                  &srcworkspace.SourceURL,
				StructuredRunOutputEnabled: &srcworkspace.StructuredRunOutputEnabled,
				TerraformVersion:           &srcworkspace.TerraformVersion,
				// TriggerPrefixes:            []string{},
				// TriggerPatterns:            []string{},
				VCSRepo: &tfe.VCSRepoOptions{},
				// WorkingDirectory: new(string),
				Tags: tag,
			})
			if err != nil {
				fmt.Println("Could not create Workspace.\n\n Error:", err.Error())
				return err
			}
			o.AddDeferredMessageRead("Migrated", srcworkspace.Name)
		}
	}
	return nil
}

func workspaceExists(c tfclient.ClientContexts, ws []string) error {
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
