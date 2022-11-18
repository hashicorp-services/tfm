// Copyright © 2022

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"

	"github.com/hashicorp/go-tfe"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var (

	// `tfe-migrate list organization` command
	orgListCmd = &cobra.Command{
		Use:     "organization",
		Aliases: []string{"org"},
		Short:   "List Organizations",
		Long:    "List of Organizations.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return orgList(
				GetClientContexts())
			//*viperString("search"),
			//*viperString("repository"),
			//*viperString("run-status"))
		},
		//		Run: func(cmd *cobra.Command, args []string) {
		//			orgList(getTfeDiscoverClientContextSource())
	}

	// `tfe-mig org show org-id` command
	orgShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show org attributes",
		Long:  "Show the attributes of a specific org.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return orgShow(
				*ViperString("name"))
		},
	}
)

func init() {
	// Flags().StringP, etc... - the "P" gives us the option for a short hand

	// `tfe-discover organization list` command
	orgListCmd.Flags().BoolP("all", "a", false, "List all? (optional)")

	// `tfe-discover organization show` command
	orgShowCmd.Flags().Int16P("id", "i", 0, "id of foo.")
	orgShowCmd.MarkFlagRequired("name")
	orgShowCmd.Flags().String("name", "n", "name of foo")

	// Add commands
	//RootCmd.AddCommand(listCmd)
	listCmd.AddCommand(orgListCmd)
	//orgCmd.AddCommand(orgShowCmd)

}

func orgList(c ClientContexts) error {
	o.AddMessageUserProvided("List of Organizations at: ", c.SourceHostname)
	allItems := []*tfe.Organization{}

	opts := tfe.OrganizationListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100},
	}

	for {
		items, err := c.SourceClient.Organizations.List(c.SourceContext, &opts)
		if err != nil {
			logError(err, "failed to list orgs")
		}

		allItems = append(allItems, items.Items...)

		o.AddFormattedMessageCalculated("Found %d Organizations", len(allItems))

		if items.CurrentPage >= items.TotalPages {
			break
		}
		opts.PageNumber = items.NextPage
	}

	o.AddTableHeaders("Name", "Created On", "Email")
	for _, i := range allItems {
		cr_created_at := FormatDateTime(i.CreatedAt)

		o.AddTableRows(i.Name, cr_created_at, i.Email)
	}

	return nil
}

func orgShow(name string) error {
	fmt.Println("Show org with name:", aurora.Bold(name))
	return nil
}
