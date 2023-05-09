// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package list

import (
	"fmt"

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
			listWorkspaces(tfclient.GetClientContexts())
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

func listWorkspaces(c tfclient.ClientContexts) error {

	srcWorkspaces := []*tfe.Workspace{}

	opts := tfe.WorkspaceListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100},
	}

	if (ListCmd.Flags().Lookup("side").Value.String() == "source") || (!ListCmd.Flags().Lookup("side").Changed) {

		o.AddMessageUserProvided("Getting list of workspaces from: ", c.SourceHostname)

		for {
			items, err := c.SourceClient.Workspaces.List(c.SourceContext, c.SourceOrganizationName, &opts)
			if err != nil {
				fmt.Println("Error With retrieving Workspaces from ", c.SourceHostname, " : Error ", err)
				return err
			}

			srcWorkspaces = append(srcWorkspaces, items.Items...)

			o.AddFormattedMessageCalculated("Found %d Workspaces", len(srcWorkspaces))

			if items.CurrentPage >= items.TotalPages {
				break
			}
			opts.PageNumber = items.NextPage

		}
		o.AddTableHeaders("Name", "Description", "ExecutionMode", "VCS Repo", "Project ID", "Project Name", "Locked", "TF Version")
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

			o.AddTableRows(i.Name, i.Description, i.ExecutionMode, ws_repo, projectID, projectName, i.Locked, i.TerraformVersion)
		}
	}

	if ListCmd.Flags().Lookup("side").Value.String() == "destination" {
		o.AddMessageUserProvided("Getting list of workspaces from: ", c.DestinationHostname)

		for {
			items, err := c.DestinationClient.Workspaces.List(c.DestinationContext, c.DestinationOrganizationName, &opts)
			if err != nil {
				fmt.Println("Error With retrieving Workspaces from ", c.DestinationHostname, " : Error ", err)
				return err
			}

			srcWorkspaces = append(srcWorkspaces, items.Items...)

			o.AddFormattedMessageCalculated("Found %d Workspaces", len(srcWorkspaces))

			if items.CurrentPage >= items.TotalPages {
				break
			}
			opts.PageNumber = items.NextPage

		}
		o.AddTableHeaders("Name", "Description", "ExecutionMode", "VCS Repo", "Project ID", "Project Name", "Locked", "TF Version")
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

			o.AddTableRows(i.Name, i.Description, i.ExecutionMode, ws_repo, projectID, projectName, i.Locked, i.TerraformVersion)
		}
	}

	return nil
}
