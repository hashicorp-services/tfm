package delete

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp-services/tfm/output"
	"github.com/hashicorp-services/tfm/tfclient"

	//tfe "github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

var (
	o output.Output

	workspaceId   string
	workspaceName string

	// `tfm delete workspaces` command
	workspaceDeleteCmd = &cobra.Command{
		Use:     "workspace",
		Aliases: []string{"ws"},
		Short:   "Workspace command",
		Long:    "Delete Workspace in an org",
		Run: func(cmd *cobra.Command, args []string) {

			if (workspaceId != "") && (workspaceName != "") {
				o.AddErrorUserProvided("Error: can not supply a workspace ID and a Name at the same time. Only one allowed")
				os.Exit(0)
			}

			if (workspaceId == "") && (workspaceName == "") {
				o.AddErrorUserProvided("Error: Must supply either a Workspace ID or a Workspace Name")
				os.Exit(0)
			}

			deleteWorkspace(tfclient.GetClientContexts())
		},

		PostRun: func(cmd *cobra.Command, args []string) {
			o.Close()
		},
	}
)

func init() {
	workspaceDeleteCmd.Flags().StringVar(&workspaceId, "workspace-id", "", "Specify one single workspace ID to delete")
	workspaceDeleteCmd.Flags().StringVar(&workspaceName, "workspace-name", "", "Specify one single workspace name to delete")

	// Add commands
	DeleteCmd.AddCommand(workspaceDeleteCmd)
}

func deleteWorkspace(c tfclient.ClientContexts) error {
	if workspaceId != "" {
		if DeleteCmd.Flags().Lookup("side").Value.String() == "source" || (!DeleteCmd.Flags().Lookup("side").Changed) {
			o.AddErrorUserProvided("Are you sure you want to proceed? Workspace ID " + workspaceId + " in " + c.SourceHostname + " in org " + c.SourceOrganizationName + " will be deleted")
			if confirm() {
				err := c.SourceClient.Workspaces.DeleteByID(c.SourceContext, workspaceId)
				if err == nil {
					o.AddPassUserProvided("workspace " + workspaceId + " has been deleted")
				} else {
					o.AddErrorUserProvided("There was an error deleting workspace ")
					fmt.Print(err)
				}

			} else {
				o.AddPassUserProvided("workspace not deleted")
			}
		}

		if DeleteCmd.Flags().Lookup("side").Value.String() == "destination" {
			o.AddErrorUserProvided("Are you sure you want to proceed? Workspace ID " + workspaceId + " in " + c.DestinationHostname + " in org " + c.DestinationOrganizationName + " will be deleted")
			if confirm() {
				err := c.DestinationClient.Workspaces.DeleteByID(c.DestinationContext, workspaceId)
				if err == nil {
					o.AddPassUserProvided("workspace " + workspaceId + " has been deleted")
				} else {
					o.AddErrorUserProvided("There was an error deleting workspace ")
					fmt.Print(err)
				}

			} else {
				o.AddPassUserProvided("workspace not deleted")
			}
		}
	}

	if workspaceName != "" {
		if DeleteCmd.Flags().Lookup("side").Value.String() == "source" || (!DeleteCmd.Flags().Lookup("side").Changed) {
			o.AddErrorUserProvided("Are you sure you want to proceed? Workspace Name " + workspaceName + " in " + c.SourceHostname + " in org " + c.SourceOrganizationName + " will be deleted")
			if confirm() {
				err := c.SourceClient.Workspaces.Delete(c.SourceContext, c.SourceOrganizationName, workspaceName)
				if err == nil {
					o.AddPassUserProvided("workspace " + workspaceName + " has been deleted")
				} else {
					o.AddErrorUserProvided("There was an error deleting workspace ")
					fmt.Print(err)
				}

			} else {
				o.AddPassUserProvided("workspace not deleted")
			}
		}

		if DeleteCmd.Flags().Lookup("side").Value.String() == "destination" {
			o.AddErrorUserProvided("Are you sure you want to proceed? Workspace ID " + workspaceName + " in " + c.DestinationHostname + " in org " + c.DestinationOrganizationName + " will be deleted")
			if confirm() {
				err := c.DestinationClient.Workspaces.Delete(c.DestinationContext, c.DestinationOrganizationName, workspaceName)
				if err == nil {
					o.AddPassUserProvided("workspace " + workspaceName + " has been deleted")
				} else {
					o.AddErrorUserProvided("There was an error deleting workspace ")
					fmt.Print(err)
				}

			} else {
				o.AddPassUserProvided("workspace not deleted")
			}
		}
	}

	return nil
}

func confirm() bool {

	var input string

	fmt.Printf("Do you want to continue with this operation? [y|n]: ")

	auto, err := DeleteCmd.Flags().GetBool("autoapprove")

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
