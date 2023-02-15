package copy

import (
	"fmt"
	"os"

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
	workspacesCopyCmd.Flags().BoolVarP(&agents, "agents", "", false, "Mapping of source Agent Pool IDs to destination Agent Pool IDs in config file")
	workspacesCopyCmd.Flags().BoolVarP(&vcs, "vcs", "", false, "Mapping of source vcs Oauth ID to destination vcs Oath in config file")
	workspacesCopyCmd.Flags().BoolVarP(&ssh, "ssh", "", false, "Mapping of source ssh id to destination ssh id in config file")

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
			return nil, errors.Wrap(err, "failed to list workspaces from source")
		}

	} else {
		// Get ALL source workspaces
		srcWorkspaces, err = discoverSrcWorkspaces(tfclient.GetClientContexts())
		if err != nil {
			return nil, errors.Wrap(err, "failed to list workspaces from source")
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
		return errors.Wrap(err, "failed to list workspaces from source")
	}

	// Get the destination Workspace properties
	destWorkspaces, err := discoverDestWorkspaces(tfclient.GetClientContexts())
	if err != nil {
		return errors.Wrap(err, "failed to list workspaces from destination")
	}

	// Loop each workspace in the srcWorkspaces slice, check for the workspace existence in the destination,
	// and if a workspace exists in the destination, then do nothing, else create workspace in destination.
	for _, srcworkspace := range srcWorkspaces {
		exists := doesWorkspaceExist(srcworkspace.Name, destWorkspaces)

		// Copy tags over
		var tag []*tfe.Tag
		workspaceSource := "tfm"

		for _, t := range srcworkspace.TagNames {
			tag = append(tag, &tfe.Tag{Name: t})
		}

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
				Name:               &srcworkspace.Name,
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
