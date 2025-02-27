// // Copyright (c) HashiCorp, Inc.
// // SPDX-License-Identifier: MPL-2.0

// package nuke

// import (
// 	"fmt"
// 	"strings"

// 	"github.com/hashicorp-services/tfm/output"
// 	"github.com/hashicorp-services/tfm/tfclient"
// 	tfe "github.com/hashicorp/go-tfe"
// 	"github.com/spf13/cobra"
// )

// var (
// 	o output.Output

// 	// `tfm nuke workspaces` command
// 	workspacesNukeCmd = &cobra.Command{
// 		Use:     "workspaces",
// 		Aliases: []string{"ws"},
// 		Short:   "Workspaces command",
// 		Long:    "Nuke Workspaces in an org",
// 		Run: func(cmd *cobra.Command, args []string) {
// 			nukeWorkspaces(tfclient.GetClientContexts())
// 		},
// 		PostRun: func(cmd *cobra.Command, args []string) {
// 			o.Close()
// 		},
// 	}
// )

// func init() {
// 	// Add commands
// 	NukeCmd.AddCommand(workspacesNukeCmd)
// }

// func nukeWorkspaces(c tfclient.ClientContexts) error {

// 	workspaces := listWorkspaces(c)
// 	workspacesToNuke := []*tfe.Workspace{}

// 	for _, workspace := range workspaces {
// 		if workspace.SourceName == "tfm" {
// 			workspacesToNuke = append(workspacesToNuke, workspace)
// 		}
// 	}

// 	if len(workspacesToNuke) > 0 {
// 		o.AddFormattedMessageCalculated("Found %d Workspaces created by tfm to remove", len(workspacesToNuke))
// 		o.AddTableHeaders("Name", "Description", "ExecutionMode", "VCS Repo", "Locked", "TF Version")

// 		for _, i := range workspacesToNuke {
// 			ws_repo := ""

// 			if i.VCSRepo != nil {
// 				ws_repo = i.VCSRepo.DisplayIdentifier
// 			}
// 			o.AddTableRows(i.Name, i.Description, i.ExecutionMode, ws_repo, i.Locked, i.TerraformVersion)
// 		}
// 		o.Close()

// 		o.AddFormattedMessageCalculatedDanger("Are you sure you want to proceed? %d Workspaces will be deleted!", len(workspacesToNuke))
// 		if confirm() {
// 			for _, i := range workspacesToNuke {
// 				c.DestinationClient.Workspaces.DeleteByID(c.DestinationContext, i.ID)
// 			}

// 			o.AddFormattedMessageCalculatedDanger("%d Workspaces have been nuked!", len(workspacesToNuke))

// 		} else {
// 			o.AddPassUserProvided("Nuke Disarmed!")
// 		}

// 	} else {
// 		fmt.Print("No workspaces created by tfm were found")
// 	}

// 	return nil
// }

// func listWorkspaces(c tfclient.ClientContexts) []*tfe.Workspace {

// 	workspaces := []*tfe.Workspace{}

// 	opts := tfe.WorkspaceListOptions{
// 		ListOptions: tfe.ListOptions{
// 			PageNumber: 1,
// 			PageSize:   100},
// 	}

// 	o.AddMessageUserProvided("Getting list of workspaces from: ", c.DestinationHostname)

// 	for {
// 		items, err := c.DestinationClient.Workspaces.List(c.DestinationContext, c.DestinationOrganizationName, &opts)

// 		if err != nil {
// 			fmt.Println("Error With retrieving Workspaces from ", c.DestinationHostname, " : Error ", err)
// 			// return err
// 		}

// 		workspaces = append(workspaces, items.Items...)

// 		if items.CurrentPage >= items.TotalPages {
// 			break
// 		}
// 		opts.PageNumber = items.NextPage
// 	}

// 	return workspaces
// }

// func confirm() bool {

// 	var input string

// 	fmt.Printf("Do you want to continue with this operation? [y|n]: ")

// 	auto, err := NukeCmd.Flags().GetBool("autoapprove")

// 	if err != nil {
// 		fmt.Println("Error Retrieving autoapprove flag value: ", err)
// 	}

// 	// Check if --autoapprove=false
// 	if !auto {
// 		_, err := fmt.Scanln(&input)
// 		if err != nil {
// 			panic(err)
// 		}
// 	} else {
// 		input = "y"
// 		fmt.Println("y(autoapprove=true)")
// 	}

// 	input = strings.ToLower(input)

// 	if input == "y" || input == "yes" {
// 		return true
// 	}
// 	return false
// }
