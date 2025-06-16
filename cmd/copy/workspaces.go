// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package copy

import (
	"fmt"
	"os"
	"strings"
	"time"

	"slices"

	"github.com/hashicorp-services/tfm/cmd/helper"
	"github.com/hashicorp-services/tfm/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	state         bool
	vars          bool
	skipSensitive bool
	// skipEmpty          bool // Skip empty workspaces, may reimplement later.
	teamaccess         bool
	agents             bool
	vcs                bool
	ssh                bool
	remoteStateSharing bool
	consolidateGlobal  bool
	last               int
	lock               bool
	unlock             bool
	runTriggers        bool
	createDstProject   bool
	planOnly           bool
	wsNamePrefix       string
	wsNameSuffix       string

	// `tfemigrate copy workspaces` command
	workspacesCopyCmd = &cobra.Command{
		Use:     "workspaces",
		Short:   "Copy Workspaces",
		Aliases: []string{"ws"},
		Long:    "Copy Workspaces from source to destination org",
		//ValidArgs: []string{"state", "vars"},
		//Args:      cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			// Validate `workspaces-map` if it exists before any other functions can run.
			valid, wsMapCfg, err := validateMap(tfclient.GetClientContexts(), "workspaces-map")
			if err != nil {
				return err
			}

			// Continue the application if `workspaces-map` is not provided. The valid and map output arent needed.
			_ = valid

			switch {
			case state:
				return copyStates(tfclient.GetClientContexts(), last)
			case vars:
				return copyVariables(tfclient.GetClientContexts(), skipSensitive)
			case teamaccess:
				return copyWsTeamAccess(tfclient.GetClientContexts())

			// A map is required if the --agents flag is provided. Check for a valid map.
			case agents:

				agentPoolID := viper.GetString("agent-assignment-id")
				validMap, agentPoolIDs, err := validateMap(tfclient.GetClientContexts(), "agents-map")
				if err != nil {
					return err
				}

				if len(agentPoolID) > 0 && len(agentPoolIDs) > 0 {
					o.AddErrorUserProvided("'agents-map' and 'agent-assignment-id' cannot be defined at the same time.")
					os.Exit(0)
				} else if len(agentPoolID) > 0 {
					return createAgentPoolAssignmentSingle(tfclient.GetClientContexts(), agentPoolID)

				} else if len(agentPoolIDs) > 0 {
					if !validMap {
						os.Exit(0)
					} else {
						return createAgentPoolAssignmentMap(tfclient.GetClientContexts(), agentPoolIDs)
					}
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

			case remoteStateSharing:
				return copyRemoteStateSharing(tfclient.GetClientContexts(), consolidateGlobal)

			case runTriggers:
				return copyRunTriggers(tfclient.GetClientContexts())
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
	workspacesCopyCmd.PersistentFlags().StringVar(&wsNameSuffix, "add-suffix", "", "(optional) Only performed on destination workspaces, adds a suffix, if missing.")
	workspacesCopyCmd.PersistentFlags().StringVar(&wsNamePrefix, "add-prefix", "", "(optional) Only performed on destination workspaces, adds a prefix, if missing.")
	workspacesCopyCmd.Flags().BoolVarP(&vars, "vars", "", false, "Copy workspace variables")
	workspacesCopyCmd.Flags().BoolVarP(&skipSensitive, "skip-sensitive-vars", "", false, "Skip copying sensitive variables. Must be used with --vars flag")
	// workspacesCopyCmd.Flags().BoolVarP(&skipEmpty, "skip-empty", "", false, "Skip empty workspaces.") // May reimplement later.
	workspacesCopyCmd.Flags().BoolVarP(&createDstProject, "create-dst-project", "", false, "Creates destination project, if not existing. Defaults to source organization name.")
	workspacesCopyCmd.Flags().BoolVarP(&planOnly, "plan-only", "", false, "Only plan the copy operation without making changes")

	workspacesCopyCmd.Flags().BoolVarP(&state, "state", "", false, "Copy workspace states")
	workspacesCopyCmd.Flags().IntVarP(&last, "last", "l", last, "Copy the last X number of state files only.")
	// SetInterspersed prevents cobra from parsing arguments that appear after flags
	workspacesCopyCmd.Flags().SetInterspersed(false)

	// Prevents users from using the --last flag in an unwanted way
	workspacesCopyCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if last > 0 && !state {
			return errors.New("--last flag is only valid after the --state flag is set")
		}

		return nil

	}
	workspacesCopyCmd.Flags().BoolVarP(&teamaccess, "teamaccess", "", false, "Copy workspace Team Access")
	workspacesCopyCmd.Flags().BoolVarP(&agents, "agents", "", false, "Mapping of source Agent Pool IDs to destination Agent Pool IDs in config file")
	workspacesCopyCmd.Flags().BoolVarP(&vcs, "vcs", "", false, "Mapping of source vcs Oauth ID or GitHub App ID to destination vcs Oauth or GitHub App ID in config file")
	workspacesCopyCmd.Flags().BoolVarP(&ssh, "ssh", "", false, "Mapping of source ssh id to destination ssh id in config file")
	workspacesCopyCmd.Flags().BoolVarP(&lock, "lock", "", false, "Lock all source workspaces")
	workspacesCopyCmd.Flags().BoolVarP(&unlock, "unlock", "", false, "Unlock all source workspaces")
	workspacesCopyCmd.Flags().BoolVarP(&remoteStateSharing, "remote-state-sharing", "", false, "Copy remote state sharing settings")
	workspacesCopyCmd.Flags().BoolVarP(&consolidateGlobal, "consolidate-global", "", false, "Consolidate global remote state sharing settings. Must be used with --remote-state-sharing flag")
	workspacesCopyCmd.Flags().BoolVarP(&runTriggers, "run-triggers", "", false, "Copy workspace run triggers")

	// Add commands
	CopyCmd.AddCommand(workspacesCopyCmd)
}

// Gets all workspaces from the source target
func discoverSrcWorkspaces(c tfclient.ClientContexts) ([]*tfe.Workspace, error) {
	o.AddMessageUserProvided("\nGetting list of Workspaces from: ", c.SourceHostname)
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

		o.AddFormattedMessageCalculated("\nFound %d Workspaces", len(srcWorkspaces))

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

	srcExcludeWorkspaces := viper.GetStringSlice("exclude-workspaces")

	wsMapCfg, err := helper.ViperStringSliceMap("workspaces-map")
	if err != nil {
		return srcWorkspaces, errors.New("Invalid input for workspaces-map")
	}

	if len(srcWorkspacesCfg) > 0 {
		o.AddFormattedMessageCalculated("Found %d workspaces in `workspaces` list", len(srcWorkspacesCfg))
	}

	// If no workspaces found in config (list or map), default to just assume all workspaces from source will be chosen
	if ((len(srcExcludeWorkspaces)) > 0 && len(srcWorkspacesCfg) > 0) ||
		((len(srcExcludeWorkspaces)) > 0 && len(wsMapCfg) > 0) ||
		(len(srcWorkspacesCfg) > 0 && len(wsMapCfg) > 0) {
		o.AddErrorUserProvided("In config: only one of 'exclude-workspaces', 'workspaces', or 'workspaces-map' can be defined at the same time.")
		os.Exit(0)

	} else if len(wsMapCfg) > 0 {

		// use config workspaces from map
		var wsList []string

		for key := range wsMapCfg {
			wsList = append(wsList, key)
		}

		o.AddMessageUserProvided("Source Workspaces found in `workspaces-map`:", wsList)

		// Set source workspaces
		srcWorkspaces, err = getSrcWorkspacesFilter(tfclient.GetClientContexts(), wsList)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to list Workspaces in map from source")
		}

	} else if len(srcExcludeWorkspaces) > 0 {
		o.AddMessageUserProvided("Excluding workspaces from config list:", srcExcludeWorkspaces)

		srcWorkspaces, err := discoverSrcWorkspaces(c)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to list Workspaces from source")
		}
		for _, ws := range srcExcludeWorkspaces {
			for i, w := range srcWorkspaces {
				if ws == w.Name {
					srcWorkspaces = slices.Delete(srcWorkspaces, i, i+1)
					o.AddMessageUserProvided("Excluding Workspace:", ws)
					break
				}
			}
		}

		o.AddFormattedMessageCalculated("Will migrate %d total workspaces", len(srcWorkspaces))
		return srcWorkspaces, nil

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
		o.AddMessageUserProvided2("\nWarning:\n\n", "ALL WORKSPACES WILL BE MIGRATED from", viper.GetString("src_tfe_hostname"))

		srcWorkspaces, err = discoverSrcWorkspaces(tfclient.GetClientContexts())
		if !confirm() {
			fmt.Println("\n\n**** Canceling tfm run **** ")
			os.Exit(1)
		}

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
	// Updated the message to make it more clear in the output
	o.AddMessageUserProvided("\nDiscovering workspaces in destination: ", c.DestinationHostname+"\n")
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
	projectDescription := "Generated during tfm copy"

	// Check if Project ID is set
	if viper.GetString("dst_tfc_project_id") != "" {

		project.ID = viper.GetString("dst_tfc_project_id")
		o.AddMessageUserProvided("Destination Project ID is Set: ", project.ID)

	} else if createDstProject {
		// Create a new project in the destination using the source organization name

		projectName := c.SourceOrganizationName

		// Check if the project name is already in use

		existingPrjID, err := checkDstProjectExists(c, projectName)

		if err != nil {
			return errors.Wrap(err, "Failed to check if project exists in destination")
		}
		if existingPrjID != "" {
			o.AddMessageUserProvided("Project already exists in destination: ", projectName)
			project.ID = existingPrjID
		} else {

			if !planOnly {

				fmt.Println("\n**** Creating Project ****\n ")
				projectPtr, err := c.DestinationClient.Projects.Create(c.DestinationContext, c.DestinationOrganizationName, tfe.ProjectCreateOptions{
					Name:        projectName,
					Description: &projectDescription,
				})

				if err != nil {
					return errors.Wrap(err, "Failed to create project in destination")
				}

				project = *projectPtr
				o.AddMessageUserProvided("Created new project in destination: ", project.Name)
				o.AddMessageUserProvided("Project ID: ", project.ID)
			} else {
				o.AddMessageUserProvided("Project will be created in destination: ", projectName)
			}
		}

	} else {

		// get Default Project ID
		project.ID, err = getDstDefaultProjectID(c)

		if err != nil {
			fmt.Println("Error Retrieving Destination Project ID. Destination TFE/TFC API may not be supported. Please check credentials")
			os.Exit(0)
		}
	}

	// Check if the workspace name prefix and suffix are set
	if len(wsNamePrefix) > 0 || len(wsNameSuffix) > 0 {
		// Standardize the naming convention for the workspaces
		srcWorkspaces = standardizeNamingConvention(srcWorkspaces, wsNamePrefix, wsNameSuffix)
	}

	if planOnly {
		fmt.Println("\n**** Plan Only Run ****\n ")
	} else {
		fmt.Println("\n**** Performing Workspace Copy ****\n ")
	}

	// Loop each workspace in the srcWorkspaces slice, check for the workspace existence in the destination,
	// and if a workspace exists in the destination, then do nothing, else create workspace in destination.
	for _, srcworkspace := range srcWorkspaces {
		destWorkSpaceName := srcworkspace.Name

		// Copy tags over
		var tag []*tfe.Tag
		// workspaceSource := "tfm"

		for _, t := range srcworkspace.TagNames {
			tag = append(tag, &tfe.Tag{Name: t})
		}

		// Check if the destination Workspace name differs from the source name
		if len(wsMapCfg) > 0 {
			o.AddMessageUserProvided3("Source Workspace:", srcworkspace.Name, "\nDestination Workspace:", wsMapCfg[srcworkspace.Name])
			destWorkSpaceName = wsMapCfg[srcworkspace.Name]
		}

		exists := doesWorkspaceExist(destWorkSpaceName, destWorkspaces)
		length := len(srcworkspace.Name)

		if length > 64 {
			return errors.New("Workspace name is too long. Max length is 64 characters.")
		}

		c.SourceClient.Workspaces.Read(c.SourceContext, c.SourceOrganizationName, srcworkspace.Name)
		if err != nil {
			return errors.Wrap(err, "Failed to read workspace from source")
		}

		if exists {
			// Check if the destination workspace name differs from the source name
			// Added info to clarify Destination workspace

			o.AddMessageUserProvided2(destWorkSpaceName, "exists in destination will not migrate", srcworkspace.Name)

		} else if planOnly {

			o.AddMessageUserProvided("Following workspace will be migrated: ", destWorkSpaceName)

		} else {

			migratedWorkspace, err := c.DestinationClient.Workspaces.Create(c.DestinationContext, c.DestinationOrganizationName, tfe.WorkspaceCreateOptions{
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
				Name:                       &destWorkSpaceName,
				QueueAllRuns:               &srcworkspace.QueueAllRuns,
				SpeculativeEnabled:         &srcworkspace.SpeculativeEnabled,
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
			o.AddDeferredMessageRead("Migrated", migratedWorkspace.Name)

			// Add tag to source workspace from migration.

			date := time.Now().Format("2006-01-02")

			// Validate required fields before calling AddTagBindings
			if srcworkspace.ID == "" || c.SourceContext == nil {
				return fmt.Errorf("invalid source workspace ID or context")
			}

			// Tag Bindings only work with v202502-1 and newer versions of API
			if c.SourceClient.RemoteTFEVersion() >= "v202502-1" && c.DestinationClient.RemoteTFEVersion() >= "v202502-1" {
				_, err = c.SourceClient.Workspaces.AddTagBindings(c.SourceContext, srcworkspace.ID, tfe.WorkspaceAddTagBindingsOptions{
					TagBindings: []*tfe.TagBinding{
						{Key: "migrated:", Value: "true"},
						{Key: "migration-date:", Value: date},
						{Key: "migrate-destination-workspace:", Value: destWorkSpaceName},
						{Key: "migrate-destination-workspace-id:", Value: migratedWorkspace.ID},
						{Key: "migrate-destination-hostname:", Value: c.DestinationHostname},
					},
				})
				if err != nil {
					return fmt.Errorf("failed to add tags on source workspace: %w", err)
				}

				_, err = c.DestinationClient.Workspaces.AddTagBindings(c.DestinationContext, migratedWorkspace.ID, tfe.WorkspaceAddTagBindingsOptions{
					TagBindings: []*tfe.TagBinding{
						{Key: "migrated:", Value: "true"},
						{Key: "migration-date:", Value: date},
						{Key: "migrate-source-workspace:", Value: srcworkspace.Name},
						{Key: "migrate-source-workspace-id:", Value: srcworkspace.ID},
						{Key: "migrate-source-hostname:", Value: c.SourceHostname},
					},
				})
				if err != nil {
					return fmt.Errorf("failed to add tags on destination workspace - : %w", err)
				}
			}
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

func checkDstProjectExists(c tfclient.ClientContexts, dstProjectName string) (string, error) {

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

		if i.Name == dstProjectName {
			return i.ID, nil
		}
	}

	return "", nil
}

func confirm() bool {

	var input string

	fmt.Printf("Do you want to continue with this operation? [y|n]: ")

	auto, err := CopyCmd.Flags().GetBool("autoapprove")

	if err != nil {
		fmt.Println("Error Retrieving autoapprove flag value: ", err)
	}

	// Check if --autoapprove=false
	if !auto {
		_, err := fmt.Scanln(&input)
		if err != nil {
			panic(err)
		}
	} else {
		input = "y"
		fmt.Println("y(autoapprove=true)")
	}

	input = strings.ToLower(input)

	if input == "y" || input == "yes" {
		return true
	}
	return false
}

func standardizeNamingConvention(workspaceList []*tfe.Workspace, prefix string, suffix string) []*tfe.Workspace {
	workspaceListUpdated := workspaceList

	fmt.Print("\n**** Standardizing workspace names with prefix and suffix ****\n\n")

	for _, ws := range workspaceList {

		if !strings.Contains(ws.Name, prefix) || !strings.Contains(ws.Name, suffix) {
			o.AddMessageUserProvided("Renaming with specified prefix/suffix: ", ws.Name)

			if !strings.Contains(ws.Name, prefix+"-") {
				ws.Name = prefix + "-" + ws.Name
			}

			if !strings.Contains(ws.Name, "-"+suffix) {
				ws.Name = ws.Name + "-" + suffix
			}

			if strings.Contains(ws.Name, prefix) && strings.Contains(ws.Name, suffix) {
				continue
			}

			return workspaceListUpdated
		}
	}
	return workspaceListUpdated
}
