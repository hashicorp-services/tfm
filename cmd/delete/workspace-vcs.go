package delete

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

	// `tfemigrate copy workspaces` command
	workspaceVCSDeleteCmd = &cobra.Command{
		Use:   "workspaces-vcs",
		Short: "Delete Workspaces VCS Connection",
		Long:  "Deletes a workspaces VCS Connection. This is used after a migration have been completed to remove the VCS connection on the source workspaces.",
		RunE: func(cmd *cobra.Command, args []string) error {

			// Validate `workspaces-map` if it exists before any other functions can run.
			valid, wsMapCfg, err := validateMap(tfclient.GetClientContexts(), "workspaces-map")
			if err != nil {
				return err
			}

			// Continue the application if `workspaces-map` is not provided. The valid and map output arent needed.
			_ = valid

			return deleteWorkspaceVCSConnection(tfclient.GetClientContexts(), wsMapCfg)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			o.Close()
		},
	}
)

func init() {

	// Add commands
	DeleteCmd.AddCommand(workspaceVCSDeleteCmd)

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

	wsMapCfg, err := helper.ViperStringSliceMap("workspaces-map")
	if err != nil {
		return srcWorkspaces, errors.New("Invalid input for workspaces-map")
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

		o.AddMessageUserProvided("Source Workspaces found in `workspaces-map`:", wsList)

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
		fmt.Println("No workspaces or workspaces-map found in config file (~/.tfm.hcl).\n\nALL WORKSPACES WILL HAVE THEIR VCS CONNECTION REMOVED from ", viper.GetString("src_tfe_hostname"))
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

func deleteWorkspaceVCSConnection(c tfclient.ClientContexts, wsMapCfg map[string]string) error {

	// Get Workspaces from Config OR get ALL workspaces from source
	srcWorkspaces, err := getSrcWorkspacesCfg(c)
	if err != nil {
		return errors.Wrap(err, "Failed to list workspaces from source target")
	}

	for _, srcworkspace := range srcWorkspaces {

		_, err := c.SourceClient.Workspaces.RemoveVCSConnectionByID(c.SourceContext, srcworkspace.ID) 
		if err == nil {
			o.AddPassUserProvided("workspace " + srcworkspace.Name + " VCS Connection has been removed")
		}

	}

	return nil

}