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
	state             bool
	vars              bool
	teamaccess        bool
	agents            bool
	vcs 			  bool
	sourcePoolID      string
	destinationPoolID string

	// `tfemigrate copy workspaces` command
	workspacesCopyCmd = &cobra.Command{
		Use:     "workspaces",
		Short:   "Copy Workspaces",
		Aliases: []string{"ws"},
		Long:    "Copy Workspaces from source to destination org",
		//ValidArgs: []string{"state", "vars"},
		//Args:      cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			agentpools, err := helper.ViperStringSliceMap("agentpools-map")
			if err != nil {
				return errors.New("invalid input for 'agentpools'")
			}
			vcsIDs, err := helper.ViperStringSliceMap("vcs-map")
			if err != nil {
				return errors.New("invalid input for vcs")
			}
			switch {
			case state:
				return copyStates(tfclient.GetClientContexts())
			case vars:
				return copyVariables(tfclient.GetClientContexts())
			case teamaccess:
				return copyWsTeamAccess(tfclient.GetClientContexts())
			case agents:
				return createAgentPoolAssignment(tfclient.GetClientContexts(), agentpools) //tfm copy workspaces --agents --source-pool-id x --destination-pool-id
			case vcs:
				return createVCSConfiguration(tfclient.GetClientContexts(), vcsIDs)
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
	workspacesCopyCmd.Flags().BoolVarP(&vars, "vars", "v", false, "Copy workspace variables")
	workspacesCopyCmd.Flags().BoolVarP(&state, "state", "s", false, "Copy workspace states")
	workspacesCopyCmd.Flags().BoolVarP(&teamaccess, "teamaccess", "t", false, "Copy workspace Team Access")
	workspacesCopyCmd.Flags().BoolVarP(&agents, "agents", "g", false, "Assign Agent Pool IDs based on source Pool ID")
	workspacesCopyCmd.Flags().StringSliceP("agentpools-map", "p", []string{}, "Mapping of source agent pool to destination agent pool. Can be supplied multiple times. (optional, i.e. '--agentpools='apool-DgzkahoomwHsBHcJ=apool-vbrJZKLnPy6aLVxE')")
	workspacesCopyCmd.Flags().BoolVarP(&vcs, "vcs", "o", false, "Assign VCS Oauth IDs based on source VCS Oauth ID")
	workspacesCopyCmd.Flags().StringSliceP("vcs-map", "m", []string{}, "Mapping of source vcs oauth id to destination vcs oath id. Can be supplied multiple times. (optional, i.e. '--vcs='oc-UAgBKNE4WNUH4kPM=oc-A324BNKExwefmo13')")
	//workspacesCopyCmd.Flags().StringVarP(&sourcePoolID, "source-pool-id", "m", "", "The source Agent Pool ID (required if agent set)")
	//workspacesCopyCmd.Flags().StringVarP(&destinationPoolID, "destination-pool-id", "n", "", "the destination Agent Pool ID (required if agent set)")
	//workspacesCopyCmd.MarkFlagsRequiredTogether("agents", "source-pool-id", "destination-pool-id")

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

func getSrcWorkspacesCfg(c tfclient.ClientContexts) ([]*tfe.Workspace, error) {
	// Get Workspaces from Config
	srcWorkspacesCfg := viper.GetStringSlice("workspaces")
	var srcWorkspaces []*tfe.Workspace

	o.AddFormattedMessageCalculated("Found %d Workspaces in Configuration", len(srcWorkspacesCfg))
	var err error
	// If not workspaces found in config, default to just assume all workspaces from source will be chosen
	if len(srcWorkspacesCfg) > 0 {
		// use config workspaces
		fmt.Println("Using workspaces config list:", srcWorkspacesCfg)

		//get source workspaces
		srcWorkspaces, err = getSrcWorkspacesFilter(tfclient.GetClientContexts(), srcWorkspacesCfg)
		if err != nil {
			return nil, errors.Wrap(err, "failed to list teams from source")
		}

	} else {
		// Get ALL source workspaces
		srcWorkspaces, err = discoverSrcWorkspaces(tfclient.GetClientContexts())
		if err != nil {
			return nil, errors.Wrap(err, "failed to list teams from source")
		}
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

	return srcWorkspaces, nil

}


func getSrcWorkspacesFilter(c tfclient.ClientContexts, wsList []string) ([]*tfe.Workspace, error) {
	o.AddMessageUserProvided("Getting list of workspaces from: ", c.SourceHostname)
	srcWorkspaces := []*tfe.Workspace{}
	
	fmt.Println("Workspace list from config:", wsList)

	for _, ws := range wsList {

		for {
			opts := tfe.WorkspaceListOptions{
				ListOptions: tfe.ListOptions{
					PageNumber: 1,
					PageSize:   100},
					Search: ws,
			}
	
			items, err := c.SourceClient.Workspaces.List(c.SourceContext, c.SourceOrganizationName, &opts) // This should only return 1 result
			if err != nil {
				return nil, err
			}
	
			srcWorkspaces = append(srcWorkspaces, items.Items...)
		
			if items.CurrentPage >= items.TotalPages {
				break
			}
			opts.PageNumber = items.NextPage

		}
	}

	return srcWorkspaces, nil
}

func getDstWorkspacesFilter(c tfclient.ClientContexts, wsList []string) ([]*tfe.Workspace, error) {
	o.AddMessageUserProvided("Getting list of workspaces from: ", c.DestinationHostname)
	dstWorkspaces := []*tfe.Workspace{}
	
	fmt.Println("Workspace list from config:", wsList)

	for _, ws := range wsList {

		for {
			opts := tfe.WorkspaceListOptions{
				ListOptions: tfe.ListOptions{
					PageNumber: 1,
					PageSize:   100},
					Search: ws,
			}
	
			items, err := c.DestinationClient.Workspaces.List(c.DestinationContext, c.DestinationOrganizationName, &opts) // This should only return 1 result
			if err != nil {
				return nil, err
			}
	
			dstWorkspaces = append(dstWorkspaces, items.Items...)
		
			if items.CurrentPage >= items.TotalPages {
				break
			}
			opts.PageNumber = items.NextPage

		}
	}

	return dstWorkspaces, nil
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

// Takes a workspace name and a slice of workspace as type []*tfe.Workspace and
// returns true if the workspacee name exists within the provided slice of workspaces.
// Used to compare source workspace names to the destination workspace names.
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

	// Get Workspaces from Config OR get ALL workspaces from source
	srcWorkspaces, err := getSrcWorkspacesCfg(c)
	if err != nil {
		return errors.Wrap(err, "failed to list teams from source")
	}

	// Get the destination Workspace properties
	destWorkspaces, err := discoverDestWorkspaces(tfclient.GetClientContexts())
	if err != nil {
		return errors.Wrap(err, "failed to list teams from destination")
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
				//VCSRepo: &tfe.VCSRepoOptions{},
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
