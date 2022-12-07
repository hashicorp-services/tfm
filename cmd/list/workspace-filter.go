package list

import (

    //"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp-services/tfe-mig/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

var (
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
}

type workspaceFilterOptions struct {
	searchString string
	tagsString string
	excludedTags string
	workspaceIncludes []tfe.WSIncludeOpt
}

func workspaceFilter(c tfclient.ClientContexts, w workspaceFilterOptions) error {

	allItems := []*tfe.Workspace{}
	opts := tfe.WorkspaceListOptions{
		ListOptions: tfe.ListOptions{PageNumber: 1, PageSize: 100},
		Search:      w.searchString,
		Tags:        w.tagsString,
		ExcludeTags: w.excludedTags,
		Include:     w.workspaceIncludes,
	}

	for {
		items, err := c.SourceClient.Workspaces.List(c.SourceContext, c.SourceOrganizationName, &opts)

		if err != nil {
			return err
		}

		allItems = append(allItems, items.Items...)

		if items.CurrentPage >= items.TotalPages {
			break
		}
		opts.PageNumber = items.NextPage
	}
	//spew.Dump(len(allItems))
	return nil
}
