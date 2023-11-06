// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package unlock

import (
	"fmt"
	"github.com/hashicorp-services/tfm/cmd/helper"
	"github.com/hashicorp-services/tfm/output"
	"github.com/hashicorp-services/tfm/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	o output.Output

	// `tfm list workspaces` command
	workspacesUnlockCmd = &cobra.Command{
		Use:     "workspaces",
		Aliases: []string{"ws"},
		Short:   "Workspaces command",
		Long:    "Unlock Workspaces in an org",
		Run: func(cmd *cobra.Command, args []string) {
			workspaceUnlock(tfclient.GetClientContexts())
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			o.Close()
		},
	}
)

func init() {
	// Add commands
	UnlockCmd.AddCommand(workspacesUnlockCmd)
}

// All functions related to unlocking a workspace
func workspaceUnlock(c tfclient.ClientContexts) error {

	if (UnlockCmd.Flags().Lookup("side").Value.String() == "source") || (!UnlockCmd.Flags().Lookup("side").Changed) {
		o.AddMessageUserProvided("Unlocking configured workspaces on:", c.SourceHostname)

		// Get the source workspaces properties
		srcWorkspaces, err := getSrcWorkspacesCfg(c)
		if err != nil {
			return errors.Wrap(err, "failed to list Workspaces from source while checking lock status")
		}

		for _, ws := range srcWorkspaces {

			wsProperties, err := c.SourceClient.Workspaces.ReadByID(c.SourceContext, ws.ID)
			if err != nil {
				return err
			}

			if wsProperties.Locked {
				fmt.Println("Unlocking Workspace:", ws.Name)
				unlockStats, unlockErr := c.SourceClient.Workspaces.Unlock(c.SourceContext, ws.ID)

				if unlockErr != nil {
					return unlockErr
				}

				_ = unlockStats
			} else {
				fmt.Println("Workspace is already unlocked:", ws.Name)
			}
		}
	}

	if UnlockCmd.Flags().Lookup("side").Value.String() == "destination" {
		o.AddMessageUserProvided("Unlocking all configured workspaces on:", c.DestinationHostname)

		// Get the destination workspaces properties
		dstWorkspaces, err := getDstWorkspacesCfg(c)
		if err != nil {
			return errors.Wrap(err, "failed to list Workspaces from destination while checking lock status")
		}

		for _, ws := range dstWorkspaces {

			wsProperties, err := c.DestinationClient.Workspaces.ReadByID(c.DestinationContext, ws.ID)
			if err != nil {
				return err
			}

			if wsProperties.Locked {
				fmt.Println("Unlocking Workspace:", ws.Name)
				unlockStats, unlockErr := c.DestinationClient.Workspaces.Unlock(c.DestinationContext, ws.ID)

				if unlockErr != nil {
					return unlockErr
				}

				_ = unlockStats
			} else {
				fmt.Println("Workspace is already unlocked:", ws.Name)
			}
		}
	}

	return nil
}

func getSrcWorkspacesFilter(c tfclient.ClientContexts, wsList []string) ([]*tfe.Workspace, error) {
	// o.AddMessageUserProvided("Getting list of Workspaces from:", c.SourceHostname)
	
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

func discoverSrcWorkspaces(c tfclient.ClientContexts, output bool) ([]*tfe.Workspace, error) {
	// o.AddMessageUserProvided("\nGetting list of Workspaces from:", c.SourceHostname)
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

		if output {
			o.AddFormattedMessageCalculated("\nFound %d Workspaces", len(srcWorkspaces))
		}

		if items.CurrentPage >= items.TotalPages {
			break
		}
		opts.PageNumber = items.NextPage

	}

	return srcWorkspaces, nil
}

func getDstWorkspacesFilter(c tfclient.ClientContexts, wsList []string) ([]*tfe.Workspace, error) {
	// o.AddMessageUserProvided("Getting list of Workspaces from:", c.DestinationHostname)
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

func discoverDstWorkspaces(c tfclient.ClientContexts, output bool) ([]*tfe.Workspace, error) {
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

func doesWorkspaceExist(workspaceName string, ws []*tfe.Workspace) bool {
	for _, w := range ws {
		if workspaceName == w.Name {
			return true
		}
	}
	return false
}

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
	if len(srcWorkspacesCfg) > 0 && len(wsMapCfg) > 0 {
		o.AddErrorUserProvided("'workspaces' list and 'workpaces-map' cannot be defined at the same time.")
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
		o.AddMessageUserProvided2("\nWarning:\n\n", "ALL WORKSPACES WILL BE UNLOCKED in", viper.GetString("src_tfe_hostname"))

		srcWorkspaces, err = discoverSrcWorkspaces(tfclient.GetClientContexts(), false)

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

func getDstWorkspacesCfg(c tfclient.ClientContexts) ([]*tfe.Workspace, error) {

	var dstWorkspaces []*tfe.Workspace

	// Get Destination Workspace list from config list `workspaces` if it exists
	dstWorkspacesCfg := viper.GetStringSlice("workspaces")

	wsMapCfg, err := helper.ViperStringSliceMap("workspaces-map")
	if err != nil {
		return dstWorkspaces, errors.New("Invalid input for workspaces-map")
	}

	if len(dstWorkspacesCfg) > 0 {
		o.AddFormattedMessageCalculated("Found %d workspaces in `workspaces` list", len(dstWorkspacesCfg))
	}

	// If no workspaces found in config (list or map), default to just assume all workspaces from Destination will be chosen
	if len(dstWorkspacesCfg) > 0 && len(wsMapCfg) > 0 {
		o.AddErrorUserProvided("'workspaces' list and 'workpaces-map' cannot be defined at the same time.")
		os.Exit(0)

	} else if len(wsMapCfg) > 0 {

		// use config workspaces from map
		var wsList []string

		for key := range wsMapCfg {
			wsList = append(wsList, wsMapCfg[key])
		}
		o.AddMessageUserProvided("Destination Workspaces found in `workspaces-map`:", wsList)

		// Set destination workspaces
		dstWorkspaces, err = getDstWorkspacesFilter(tfclient.GetClientContexts(), wsList)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to list Workspaces in map from destination")
		}

	} else if len(dstWorkspacesCfg) > 0 {
		// use config workspaces from list

		fmt.Println("Using Workspaces config list:", dstWorkspacesCfg)

		//get destination workspaces
		dstWorkspaces, err = getDstWorkspacesFilter(tfclient.GetClientContexts(), dstWorkspacesCfg)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to list Workspaces from destination")
		}

	} else {
		// Get ALL destination workspaces
		o.AddMessageUserProvided2("\nWarning:\n\n", "ALL WORKSPACES WILL BE UNLOCKED in", viper.GetString("src_tfe_hostname"))

		dstWorkspaces, err = discoverDstWorkspaces(tfclient.GetClientContexts(), false)

		if err != nil {
			return nil, errors.Wrap(err, "Failed to list Workspaces from destination")
		}
	}

	// Check Workspaces exist in source from config
	for _, s := range dstWorkspacesCfg {
		//fmt.Println("\nFound Workspace ", s, "in config, check if it exists in", viper.GetString("src_tfe_hostname"))
		exists := doesWorkspaceExist(s, dstWorkspaces)
		if !exists {
			fmt.Printf("Defined Workspace in config %s DOES NOT exist in %s. \n Please validate your configuration.", s, viper.GetString("src_tfe_hostname"))
			break
		}
	}

	return dstWorkspaces, nil
}
