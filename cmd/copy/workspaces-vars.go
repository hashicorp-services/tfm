
package copy

import (
	"fmt"
	"github.com/hashicorp-services/tfe-mig/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	//"github.com/spf13/cobra"
)

/*
var (

	// `tfemigrate copy workspaces` command
	variableCopyCmd = &cobra.Command{
		Use:     "var",
		Aliases: []string{"var"},
		Short:   "Copy Variables",
		Long:    "Copy Variables from source to destination workspace",
		Run: func(cmd *cobra.Command, args []string) {
			variableCopyHelper(tfclient.GetClientContexts())

		},
	}
)

func init() {
	CopyCmd.AddCommand(variableCopyCmd)
}

func variableCopyHelper(c tfclient.ClientContexts) {
	fmt.Print(variableCopy(c, "ws-VPw9i3fJhi2K5EqM", "ws-4gn2wuKYDJPohw5u"))
}

*/
// above is not going to be used except for testing. Function below will be utilized by workspace copy cmd
func variableCopy(c tfclient.ClientContexts, sourceWorkspaceID string, destinationWorkspaceID string) error {

	variableListOpts := tfe.VariableListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100,
		},
	}

	//get all variables in workspace
	workspaceVars, err := c.SourceClient.Variables.List(c.SourceContext, sourceWorkspaceID, &variableListOpts )

	if err != nil {
		fmt.Println("Could not list workspace variables.\n\n Error:", err.Error())
		return err
	}

	for _, workspaceVar := range workspaceVars.Items {
		
		//gather variables from source workspace. Variables marked as sensitive will be set to "" in the destination
		variableOpts := tfe.VariableCreateOptions{ 
			Type: "",
			Key: &workspaceVar.Key,
			Value: &workspaceVar.Value,
			Description: &workspaceVar.Description,
			Category: &workspaceVar.Category,
			HCL: &workspaceVar.HCL,
			Sensitive: &workspaceVar.Sensitive,
		}
		
		//Create the variable
		_, err := c.DestinationClient.Variables.Create(c.DestinationContext, destinationWorkspaceID, variableOpts)
		
		if err != nil {
			fmt.Println("Could not create Workspace variable.\n\n Error:", err.Error())
			return err
		}
	}

	return nil
}
