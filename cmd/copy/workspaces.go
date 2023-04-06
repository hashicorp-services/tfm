package copy

import (
	"fmt"
	"os"

	"github.com/hashicorp-services/tfm/cmd/helper"
	"github.com/hashicorp-services/tfm/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	state      bool
	vars       bool
	teamaccess bool
	agents     bool
	vcs        bool
	ssh        bool

	// `tfemigrate copy workspaces` command
	workspacesCopyCmd = &cobra.Command{
		Use:     "workspaces",
		Short:   "Copy Workspaces",
		Aliases: []string{"ws"},
		Long:    "Copy Workspaces from source to destination org",
		//ValidArgs: []string{"state", "vars"},
		//Args:      cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			// Validate `workspace-map` if it exists before any other functions can run.
			valid, wsMapCfg, err := validateMap(tfclient.GetClientContexts(), "workspace-map")
			if err != nil {
				return err
			}

			// Continue the application if `workspace-map` is not provided. The valid and map output arent needed.
			_ = valid

			switch {
			case state:
				return copyStates(tfclient.GetClientContexts())
			case vars:
				return copyVariables(tfclient.GetClientContexts())
			case teamaccess:
				return copyWsTeamAccess(tfclient.GetClientContexts())

			// A map is required if the --agents flag is provided. Check for a valid map.
			case agents:
				valid, apoolIDs, err := validateMap(tfclient.GetClientContexts(), "agents-map")
				if err != nil {
					return err
				}

				if !valid {
					os.Exit(0)
				} else {
					return createAgentPoolAssignment(tfclient.GetClientContexts(), apoolIDs)
				}

			// A map is required if the --vcs flag is provided. Check for a valid map.
			case vcs:
				valid, vcsIDs, err := validateMap(tfclient.GetClientContexts(), "vcs-map")
				if err != nil {
					return err
				}

				if !valid {
					os.Exit(0)
				} else {
					return createVCSConfiguration(tfclient.GetClientContexts(), vcsIDs)
				}

			// A map is required if the --ssh flag is provided. Check for a valid map.
			case ssh:
				valid, sshIDs, err := validateMap(tfclient.GetClientContexts(), "ssh-map")
				if err != nil {
					return err
				}
				if !valid {
					os.Exit(0)
				} else {
					return createSSHConfiguration(tfclient.GetClientContexts(), sshIDs)
				}
			}

			return copyWorkspaces(
				tfclient.GetClientContexts(), wsMapCfg)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			o.Close()
		},
	}
)

func init() {

	// `tfemigrate copy workspaces --workspace-id [WORKSPACEID]`
	workspacesCopyCmd.Flags().String("workspace-id", "", "Specify one single workspace ID to copy to destination")
	workspacesCopyCmd.Flags().BoolVarP(&vars, "vars", "", false, "Copy workspace variables")
	workspacesCopyCmd.Flags().BoolVarP(&state, "state", "", false, "Copy workspace states")
	workspacesCopyCmd.Flags().BoolVarP(&teamaccess, "teamaccess", "", false, "Copy workspace Team Access")
	workspacesCopyCmd.Flags().BoolVarP(&agents, "agents", "", false, "Mapping of source Agent Pool IDs to destination Agent Pool IDs in config file")
	workspacesCopyCmd.Flags().BoolVarP(&vcs, "vcs", "", false, "Mapping of source vcs Oauth ID to destination vcs Oath in config file")
	workspacesCopyCmd.Flags().BoolVarP(&ssh, "ssh", "", false, "Mapping of source ssh id to destination ssh id in config file")

	// Add commands
	CopyCmd.AddCommand(workspacesCopyCmd)

}

// Gets all workspaces from the source target
func discoverSrcWorkspaces(c tfclient.ClientContexts) ([]*tfe.Workspace, error) {
	o.AddMessageUserProvided("Getting list of Workspaces from: ", c.SourceHostname)
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

// Gets all workspaces defined in the configuration file `workspaces` or `workspaces-map` lists from the source target
func getSrcWorkspacesCfg(c tfclient.ClientContexts) ([]*tfe.Workspace, error) {

	var srcWorkspaces []*tfe.Workspace

	// Get source Workspace list from config list `workspaces` if it exists
	srcWorkspacesCfg := viper.GetStringSlice("workspaces")

	wsMapCfg, err := helper.ViperStringSliceMap("workspace-map")
	if err != nil {
		return srcWorkspaces, errors.New("Invalid input for workspace-map")
	}

	if len(srcWorkspacesCfg) > 0 {
		o.AddFormattedMessageCalculated("Found %d workspaces in `workspaces` list", len(srcWorkspacesCfg))
	}

	// If no workspaces found in config (list or map), default to just assume all workspaces from source will be chosen
	if len(wsMapCfg) > 0 {

		// use config workspaces from map
		var wsList []string

		for key := range wsMapCfg {
			wsList = append(wsList, key)
		}

		o.AddMessageUserProvided("Source Workspaces found in `workspace-map`:", wsList)

		// Set source workspaces
		srcWorkspaces, err = getSrcWorkspacesFilter(tfclient.GetClientContexts(), wsList)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to list Workspaces in map from source")
		}

	} else if len(srcWorkspacesCfg) > 0 {
		// use config workspaces from list

		fmt.Println("Using Workspaces config list:", srcWorkspacesCfg)

		//get source workspaces
		srcWorkspaces, err = getSrcWorkspacesFilter(tfclient.GetClientContexts(), srcWorkspacesCfg)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to list Workspaces from source")
		}

	} else {
		// Get ALL source workspaces
		srcWorkspaces, err = discoverSrcWorkspaces(tfclient.GetClientContexts())
		if err != nil {
			return nil, errors.Wrap(err, "Failed to list Workspaces from source")
		}
	}

	// Check Workspaces exist in source from config
	for _, s := range srcWorkspacesCfg {
		//fmt.Println("\nFound Workspace ", s, "in config, check if it exists in", viper.GetString("src_tfe_hostname"))
		exists := doesWorkspaceExist(s, srcWorkspaces)
		if !exists {
			fmt.Printf("Defined Workspace in config %s DOES NOT exist in %s. \n Please validate your configuration.", s, viper.GetString("src_tfe_hostname"))
			break
		}
	}

	return srcWorkspaces, nil

}

func getSrcWorkspacesFilter(c tfclient.ClientContexts, wsList []string) ([]*tfe.Workspace, error) {
	o.AddMessageUserProvided("Getting list of Workspaces from: ", c.SourceHostname)
	srcWorkspaces := []*tfe.Workspace{}

	//fmt.Println("Workspace list from config:", wsList)

	for _, ws := range wsList {

		for {
			opts := tfe.WorkspaceListOptions{
				ListOptions: tfe.ListOptions{
					PageNumber: 1,
					PageSize:   100},
				Search: ws,
			}

			items, err := c.SourceClient.Workspaces.List(c.SourceContext, c.SourceOrganizationName, &opts) // This should only return 1 result

			indexMatch := 0

			// If multiple workspaces named similar, find exact match
			if len(items.Items) > 1 {
				for _, result := range items.Items {
					if ws == result.Name {
						// Finding matching workspace name
						break
					}
					indexMatch++
				}
			}

			if err != nil {
				return nil, err
			}

			srcWorkspaces = append(srcWorkspaces, items.Items[indexMatch])

			if items.CurrentPage >= items.TotalPages {
				break
			}
			opts.PageNumber = items.NextPage

		}
	}

	return srcWorkspaces, nil
}

func getDstWorkspacesFilter(c tfclient.ClientContexts, wsList []string) ([]*tfe.Workspace, error) {
	o.AddMessageUserProvided("Getting list of Workspaces from: ", c.DestinationHostname)
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

			// This should only return 1 result
			items, err := c.DestinationClient.Workspaces.List(c.DestinationContext, c.DestinationOrganizationName, &opts)
			if err != nil {
				return nil, err
			}

			indexMatch := 0

			// If multiple workspaces named similar, find exact match
			if len(items.Items) > 1 {
				for _, result := range items.Items {
					if ws == result.Name {
						// Finding matching workspace name
						break
					}
					indexMatch++
				}
			}

			dstWorkspaces = append(dstWorkspaces, items.Items[indexMatch])

			if items.CurrentPage >= items.TotalPages {
				break
			}
			opts.PageNumber = items.NextPage

		}
	}

	return dstWorkspaces, nil
}

func discoverDestWorkspaces(c tfclient.ClientContexts, output bool) ([]*tfe.Workspace, error) {
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

		if output {
			o.AddFormattedMessageCalculated("Found %d Workspaces in destination target", len(destWorkspaces))
		}

		if items.CurrentPage >= items.TotalPages {
			break
		}
		opts.PageNumber = items.NextPage

	}

	return destWorkspaces, nil
}

// Takes a workspace name and a slice of workspace as type []*tfe.Workspace and
// returns true if the workspace name exists within the provided slice of workspaces.
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
func copyWorkspaces(c tfclient.ClientContexts, wsMapCfg map[string]string) error {

	// Get Workspaces from Config OR get ALL workspaces from source
	srcWorkspaces, err := getSrcWorkspacesCfg(c)
	if err != nil {
		return errors.Wrap(err, "Failed to list workspaces from source target")
	}

	// Get the destination Workspace properties
	destWorkspaces, err := discoverDestWorkspaces(tfclient.GetClientContexts(), false)
	if err != nil {
		return errors.Wrap(err, "Failed to list workspaces from destination target")
	}

	var project tfe.Project

	// Check if Project ID is set
	if viper.GetString("projectID") != "" {
		project.ID = viper.GetString("dst_tfc_project_id")
		fmt.Println("Project ID Set")
		fmt.Printf(" Project ID, %s, Project Name, %s", project.ID, project.Name)
		// os.Exit(0)
	} else {

		// get Default Project ID
		project.ID, err = getDstDefaultProjectID(c)

		if err != nil {
			fmt.Println("Error Retrieving Destination Project ID. Destination TFE/TFC API may not be supported. Please check credentials")
			os.Exit(0)
		}
	}

	// Loop each workspace in the srcWorkspaces slice, check for the workspace existence in the destination,
	// and if a workspace exists in the destination, then do nothing, else create workspace in destination.
	for _, srcworkspace := range srcWorkspaces {
		destWorkSpaceName := srcworkspace.Name

		// Copy tags over
		var tag []*tfe.Tag
		workspaceSource := "tfm"

		for _, t := range srcworkspace.TagNames {
			tag = append(tag, &tfe.Tag{Name: t})
		}

		// Check if the destination Workspace name differs from the source name
		if len(wsMapCfg) > 0 {
			o.AddMessageUserProvided3("Source Workspace:", srcworkspace.Name, "\nDestination Workspace:", wsMapCfg[srcworkspace.Name])
			destWorkSpaceName = wsMapCfg[srcworkspace.Name]
		}

		exists := doesWorkspaceExist(destWorkSpaceName, destWorkspaces)

		if exists {
			o.AddMessageUserProvided2(destWorkSpaceName, "exists in destination will not migrate", srcworkspace.Name)
		} else {
			srcworkspace, err := c.DestinationClient.Workspaces.Create(c.DestinationContext, c.DestinationOrganizationName, tfe.WorkspaceCreateOptions{
				Type: "",
				// AgentPoolID:        new(string), covered with `assignAgentPool` function
				AllowDestroyPlan:   &srcworkspace.AllowDestroyPlan,
				AssessmentsEnabled: &srcworkspace.AssessmentsEnabled,
				AutoApply:          &srcworkspace.AutoApply,
				Description:        &srcworkspace.Description,
				//ExecutionMode:       &srcworkspace.ExecutionMode,
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
				Project:          &project,
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

func getDstDefaultProjectID(c tfclient.ClientContexts) (string, error) {

	dstProjects := []*tfe.Project{}

	opts := tfe.ProjectListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100},
	}

	for {
		items, err := c.DestinationClient.Projects.List(c.DestinationContext, c.DestinationOrganizationName, &opts)
		if err != nil {
			return "", err
		}
		dstProjects = append(dstProjects, items.Items...)

		if items.CurrentPage >= items.TotalPages {
			break
		}
		opts.PageNumber = items.NextPage

	}

	for _, i := range dstProjects {

		if i.Name == "Default Project" {
			return i.ID, nil
		}
	}

	return "", nil
}
