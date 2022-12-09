package list

import (
	"fmt"
	"github.com/hashicorp-services/tfe-mig/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

var (
	searchString      string
	tagsString        string
	excludedTags      string
	wildcardName      string
	//workspaceIncludes tfe.WSIncludeOpt

	workspaceFilterCmd = &cobra.Command{
		Use:     "workspace-filter",
		Aliases: []string{"workspace-filter"},
		Short:   "Filter workspaces",
		Long:    "Filter Workspaces. Trying different ways to return workspaces",
		Run: func(cmd *cobra.Command, args []string) {
			workspaceFilter(tfclient.GetClientContexts())
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			o.Close()
		},
	}
)

func init() {
	ListCmd.AddCommand(workspaceFilterCmd)
	workspaceFilterCmd.PersistentFlags().StringVar(&searchString, "name", "", "partial workspace name used to filter the results")
	workspaceFilterCmd.PersistentFlags().StringVar(&tagsString, "tags", "", "comma-separated tag names used to filter the results")
	workspaceFilterCmd.PersistentFlags().StringVar(&excludedTags, "excluded-tags", "", "comma-separated tag names to exclude")
	workspaceFilterCmd.PersistentFlags().StringVar(&wildcardName, "wildcard-name", "", "workspace name to match with a wildcard")
	//workspaceFilterCmd.PersistentFlags().StringSlice("workspace-includes", []string{}, "Additional relations to include")
}

func workspaceFilter(c tfclient.ClientContexts) error {

	allItems := []*tfe.Workspace{}

	workspaceFilterOpts := tfe.WorkspaceListOptions{
		ListOptions:  tfe.ListOptions{PageNumber: 1, PageSize: 100},
		Search:       searchString,
		Tags:         tagsString,
		ExcludeTags:  excludedTags,
		WildcardName: wildcardName,
		//Include:      workspaceIncludes,
		//Include:      workspaceIncludes,
	}

	for {
		items, err := c.SourceClient.Workspaces.List(c.SourceContext, c.SourceOrganizationName, &workspaceFilterOpts)

		if err != nil {
			return err
		}

		allItems = append(allItems, items.Items...)

		if items.CurrentPage >= items.TotalPages {
			break
		}
		workspaceFilterOpts.PageNumber = items.NextPage
	}

	fmt.Print(len(allItems))

	return nil

}
