// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package list

import (
	"fmt"

	"encoding/json"

	"github.com/hashicorp-services/tfm/tfclient"
	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

var (

	// `tfm list workspaces` command
	workspacesListCmd = &cobra.Command{
		Use:     "workspaces",
		Aliases: []string{"ws"},
		Short:   "Workspaces command",
		Long:    "List Workspaces in an org",
		Run: func(cmd *cobra.Command, args []string) {
			listWorkspaces(tfclient.GetClientContexts(), jsonOut)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			o.Close()
		},
	}
)

func init() {

	// Add commands
	ListCmd.AddCommand(workspacesListCmd)

}

func listWorkspaces(c tfclient.ClientContexts, jsonOut bool) error {

	srcWorkspaces := []*tfe.Workspace{}
	workspaceJSON := make(map[string]interface{}) // Parent JSON object "workspace-names"
	workspaceData := []map[string]string{}        // workspace names slice to go inside parent object "workspace-names"

	opts := tfe.WorkspaceListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100},
	}

	if (ListCmd.Flags().Lookup("side").Value.String() == "source") || (!ListCmd.Flags().Lookup("side").Changed) {

		if jsonOut == false {
			o.AddMessageUserProvided("Getting list of workspaces from: ", c.SourceHostname)
		}

		for {
			items, err := c.SourceClient.Workspaces.List(c.SourceContext, c.SourceOrganizationName, &opts)
			if err != nil {
				fmt.Println("Error With retrieving Workspaces from ", c.SourceHostname, " : Error ", err)
				return err
			}

			srcWorkspaces = append(srcWorkspaces, items.Items...)

			if jsonOut == false {
				o.AddFormattedMessageCalculated("Found %d Workspaces", len(srcWorkspaces))
			}

			if items.CurrentPage >= items.TotalPages {
				break
			}
			opts.PageNumber = items.NextPage

		}

		if jsonOut == false {
			o.AddTableHeaders("Name", "Description", "ExecutionMode", "VCS Repo", "Project ID", "Project Name", "Locked", "TF Version")
		}
		for _, i := range srcWorkspaces {
			ws_repo := ""
			projectID := ""
			projectName := ""

			if i.VCSRepo != nil {
				ws_repo = i.VCSRepo.DisplayIdentifier
			}

			if i.Project != nil {
				projectID = i.Project.ID
				prjN, err := getProjectName(c.SourceClient, c.SourceContext, projectID)
				if err != nil {
					fmt.Println("Error With retrieving Project Name from ", projectID, " : Error ", err)
					return err
				}

				projectName = prjN
			}

			workspaceInfo := map[string]string{
				"name":             i.Name,
				"id":               i.ID,
				"repo":             ws_repo,
				"projectName":      projectName,
				"projectId":        i.Project.ID,
				//"agentPool":        i.AgentPool.ID,
				"terraformVersion": i.TerraformVersion,
				"executionMode":    i.ExecutionMode,
			}

			if jsonOut {
				workspaceData = append(workspaceData, workspaceInfo) // Store workspace name in the slice
			}

			if jsonOut == false {
				o.AddTableRows(i.Name, i.Description, i.ExecutionMode, ws_repo, projectID, projectName, i.Locked, i.TerraformVersion)
			}

		}
		if jsonOut {
			workspaceJSON["workspaces"] = workspaceData // Assign workspace names to the "workspaces" key

			jsonData, err := json.Marshal(workspaceJSON)
			if err != nil {
				fmt.Println("Error marshaling workspaces to JSON:", err)
				return err
			}

			fmt.Println(string(jsonData))
		}

	}

	if ListCmd.Flags().Lookup("side").Value.String() == "destination" {
		if jsonOut == false {
			o.AddMessageUserProvided("Getting list of workspaces from: ", c.DestinationHostname)
		}

		for {
			items, err := c.DestinationClient.Workspaces.List(c.DestinationContext, c.DestinationOrganizationName, &opts)
			if err != nil {
				fmt.Println("Error With retrieving Workspaces from ", c.DestinationHostname, " : Error ", err)
				return err
			}

			srcWorkspaces = append(srcWorkspaces, items.Items...)

			if jsonOut == false {
				o.AddFormattedMessageCalculated("Found %d Workspaces", len(srcWorkspaces))
			}

			if items.CurrentPage >= items.TotalPages {
				break
			}
			opts.PageNumber = items.NextPage

		}

		if jsonOut == false {
			o.AddTableHeaders("Name", "Description", "ExecutionMode", "VCS Repo", "Project ID", "Project Name", "Locked", "TF Version")
		}

		for _, i := range srcWorkspaces {
			ws_repo := ""
			projectID := ""
			projectName := ""

			if i.VCSRepo != nil {
				ws_repo = i.VCSRepo.DisplayIdentifier
			}

			if i.Project != nil {
				projectID = i.Project.ID
				prjN, err := getProjectName(c.DestinationClient, c.DestinationContext, projectID)
				if err != nil {
					fmt.Println("Error With retrieving Project Name from ", projectID, " : Error ", err)
					return err
				}

				projectName = prjN
			}

			workspaceInfo := map[string]string{
				"name":             i.Name,
				"id":               i.ID,
				"repo":             ws_repo,
				"projectName":      projectName,
				"projectId":        i.Project.ID,
				//"agentPool":        i.AgentPool.ID,
				"terraformVersion": i.TerraformVersion,
				"executionMode":    i.ExecutionMode,
			}

			if jsonOut {
				workspaceData = append(workspaceData, workspaceInfo) // Store workspace data in the slice
			}

			if jsonOut == false {
				o.AddTableRows(i.Name, i.Description, i.ExecutionMode, ws_repo, projectID, projectName, i.Locked, i.TerraformVersion)
			}
		}
		if jsonOut {

			workspaceJSON["workspaces"] = workspaceData // Assign workspace names to the "workspaces" key

			jsonData, err := json.Marshal(workspaceJSON)
			if err != nil {
				fmt.Println("Error marshaling workspaces to JSON:", err)
				return err
			}

			fmt.Println(string(jsonData))
		}
	}

	return nil
}
