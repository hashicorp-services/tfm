// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package list

import (
	"fmt"

	"github.com/hashicorp-services/tfm/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

var (
	searchString string
	tagsString   string
	excludedTags string
	wildcardName string
	includes     []string

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
	workspaceFilterCmd.PersistentFlags().StringSliceVar(&includes, "includes", nil, "Additional relations to include, comma separated, no space")
}

func workspaceFilter(c tfclient.ClientContexts) error {

	allItems := []*tfe.Workspace{}

	// converts the type of slice from []string to []tfe.WSIncludeOpt Not sure if there is a way to not need this?
	workspaceIncludes := make([]tfe.WSIncludeOpt, len(includes))
	for i, v := range includes {
		workspaceIncludes[i] = tfe.WSIncludeOpt(v)
	}

	workspaceFilterOpts := tfe.WorkspaceListOptions{
		ListOptions:  tfe.ListOptions{PageNumber: 1, PageSize: 100},
		Search:       searchString,
		Tags:         tagsString,
		ExcludeTags:  excludedTags,
		WildcardName: wildcardName,
		Include:      workspaceIncludes,
	}

	for {
		var items *tfe.WorkspaceList
		var err error

		if (ListCmd.Flags().Lookup("side").Value.String() == "source") || (!ListCmd.Flags().Lookup("side").Changed) {
			items, err = c.SourceClient.Workspaces.List(c.SourceContext, c.SourceOrganizationName, &workspaceFilterOpts)
		}

		if ListCmd.Flags().Lookup("side").Value.String() == "destination" {
			items, err = c.DestinationClient.Workspaces.List(c.DestinationContext, c.DestinationOrganizationName, &workspaceFilterOpts)
		}
		if err != nil {
			return nil
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
