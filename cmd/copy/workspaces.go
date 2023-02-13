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
	vcs               bool
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
				return errors.New("invalid input for 'agentpools-map'")
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
	//workspacesCopyCmd.Flags().BoolP("all", "a", false, "Copy all Workspaces found in configuration file")

	// `tfemigrate copy workspaces --workspace-id [WORKSPACEID]`
	workspacesCopyCmd.Flags().String("workspace-id", "", "Specify one single workspace ID to copy to destination")
	workspacesCopyCmd.Flags().BoolVarP(&vars, "vars", "", false, "Copy workspace variables")
	workspacesCopyCmd.Flags().BoolVarP(&state, "state", "", false, "Copy workspace states")
	workspacesCopyCmd.Flags().BoolVarP(&teamaccess, "teamaccess", "", false, "Copy workspace Team Access")
	workspacesCopyCmd.Flags().BoolVarP(&agents, "agents", "", false, "Assign Agent Pool IDs based on source Pool ID")
	workspacesCopyCmd.Flags().StringSliceP("agentpools-map", "", []string{}, "Mapping of source agent pool to destination agent pool. Can be supplied multiple times. (optional, i.e. '--agentpools='apool-DgzkahoomwHsBHcJ=apool-vbrJZKLnPy6aLVxE')")
	workspacesCopyCmd.Flags().BoolVarP(&vcs, "vcs", "", false, "Assign VCS Oauth IDs based on source VCS Oauth ID")
	workspacesCopyCmd.Flags().StringSliceP("vcs-map", "", []string{}, "Mapping of source vcs oauth id to destination vcs oath id. Can be supplied multiple times. (optional, i.e. '--vcs='oc-UAgBKNE4WNUH4kPM=oc-A324BNKExwefmo13')")
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

// This struct is meant to be used with Workspace maps
// This is house the source and destination name of a workspace
type workspaceMigrate struct {
	srcWorkspaceName string
	dstWorkspaceName string
}

func getSrcWorkspacesCfg(c tfclient.ClientContexts) ([]*tfe.Workspace, error) {
	var srcWorkspaces []*tfe.Workspace

	// Get Workspace List from Config
	srcWorkspacesCfg := viper.GetStringSlice("workspaces")

	// Get workspace Map from Config
	fmt.Println("Reading WS Map")
	wsMapCfg, err := helper.ViperStringSliceMap("workspace-map")
	if err != nil {
		return srcWorkspaces, errors.New("invalid input for workspace-map")
	}
	fmt.Println(" WS Map Config Is: ")
	fmt.Println(wsMapCfg)

	o.AddFormattedMessageCalculated("Found %d Workspaces in a Map in Configuration", len(wsMapCfg))

	o.AddFormattedMessageCalculated("Found %d Workspaces List in Configuration", len(srcWorkspacesCfg))

	// If no workspaces found in config (list or map), default to just assume all workspaces from source will be chosen
	if len(wsMapCfg) > 0 {
		// use config workspaces from map
		fmt.Println("Using workspaces config map:", wsMapCfg)
		var wsList []string

		for key := range wsMapCfg {
			wsList = append(wsList, key)
			fmt.Println("key:", key)
		}

		fmt.Println("Source WS List from Map:", wsList)

		// Set source workspaces
		srcWorkspaces, err = getSrcWorkspacesFilter(tfclient.GetClientContexts(), wsList)
		if err != nil {
			return nil, errors.Wrap(err, "failed to list teams from source")
		}

	} else if len(srcWorkspacesCfg) > 0 {
		// use config workspaces from list
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
		fmt.Println("\nFound Workspace ", s, "in config, check if it exists in", viper.GetString("sourceHostname"))
		exists := doesWorkspaceExist(s, srcWorkspaces)
		if !exists {
			fmt.Printf("Defined Workspace in Config %s DOES NOT exist in %s. \n Please validate your configuration.", s, viper.GetString("sourceHostname"))
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

	// Get/Check if Workspace map exists
	wsMapCfg, err := helper.ViperStringSliceMap("workspace-map")
	if err != nil {
		fmt.Println("invalid input for workspace-map")
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
		destWorkSpaceName := srcworkspace.Name

		// Copy tags over
		var tag []*tfe.Tag
		workspaceSource := "tfm"

		for _, t := range srcworkspace.TagNames {
			tag = append(tag, &tfe.Tag{Name: t})
		}

		// Check if Destination Workspace Name to be Change
		if len(wsMapCfg) > 0 {
			fmt.Println("Using WS Map:", wsMapCfg)
			fmt.Println("Source Workspace:", srcworkspace.Name, "\nDestination Workspace:", wsMapCfg[srcworkspace.Name])
			destWorkSpaceName = wsMapCfg[srcworkspace.Name]
		}

		exists := doesWorkspaceExist(destWorkSpaceName, destWorkspaces)

		if exists {
			fmt.Println("Exists in destination will not migrate", srcworkspace.Name)
		} else {
			srcworkspace, err := c.DestinationClient.Workspaces.Create(c.DestinationContext, c.DestinationOrganizationName, tfe.WorkspaceCreateOptions{
				Type: "",
				// AgentPoolID:        new(string), covered with `assignAgentPool` function
				AllowDestroyPlan:    &srcworkspace.AllowDestroyPlan,
				AssessmentsEnabled:  &srcworkspace.AssessmentsEnabled,
				AutoApply:           &srcworkspace.AutoApply,
				Description:         &srcworkspace.Description,
				ExecutionMode:       &srcworkspace.ExecutionMode,
				FileTriggersEnabled: &srcworkspace.FileTriggersEnabled,
				GlobalRemoteState:   &srcworkspace.GlobalRemoteState,
				// MigrationEnvironment:       new(string), legacy usage only will not add
				Name:               &destWorkSpaceName,
				QueueAllRuns:       &srcworkspace.QueueAllRuns,
				SpeculativeEnabled: &srcworkspace.SpeculativeEnabled,
				SourceName:         &workspaceSource, // beta may remove
				// SourceURL:                  &srcworkspace.SourceURL, // beta
				StructuredRunOutputEnabled: &srcworkspace.StructuredRunOutputEnabled,
				TerraformVersion:           &srcworkspace.TerraformVersion,
				TriggerPrefixes:            srcworkspace.TriggerPrefixes,
				TriggerPatterns:            srcworkspace.TriggerPatterns,
				//VCSRepo: &tfe.VCSRepoOptions{}, covered with `configureVCSsettings` function`
				WorkingDirectory: &srcworkspace.WorkingDirectory,
				Tags:             tag,
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
