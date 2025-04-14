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

package compare

import (
	"fmt"

	"github.com/hashicorp-services/tfm/output"
	"github.com/hashicorp-services/tfm/tfclient"
	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

var (
	o       output.Output
	srcID   string
	srcType string
	dstID   string
	dstType string

	// `tfm compare workspace` command
	wsCmpCmd = &cobra.Command{
		Use:     "workspaces",
		Aliases: []string{"ws"},
		Short:   "Compare Workspaces",
		Long:    "Compare Workspaces",
		Run: func(cmd *cobra.Command, args []string) {
			wsCmp(tfclient.GetClientContexts())
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			o.Close()
		},
	}
)

func init() {
	CmpCmd.AddCommand(wsCmpCmd)
	wsCmpCmd.PersistentFlags().StringVar(&srcID, "src-id", "", "partial workspace name used to filter the results")
	wsCmpCmd.MarkFlagRequired("src-id")
	wsCmpCmd.PersistentFlags().StringVar(&srcType, "src-type", "", "comma-separated tag names used to filter the results")
	wsCmpCmd.MarkFlagRequired("src-type")
	wsCmpCmd.PersistentFlags().StringVar(&dstID, "dst-id", "", "comma-separated tag names to exclude")
	wsCmpCmd.MarkFlagRequired("dst-id")
	wsCmpCmd.PersistentFlags().StringVar(&dstType, "dst-type", "", "workspace name to match with a wildcard")
	wsCmpCmd.MarkFlagRequired("dst-type")
}

func wsCmp(c tfclient.ClientContexts) error {

	wsListOptions := tfe.WorkspaceListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100,
		},
		ProjectID: srcID,
	}

	srcWS, err := c.SourceClient.Workspaces.List(c.SourceContext, c.SourceOrganizationName, &wsListOptions)

	if err != nil {
		fmt.Printf("Error With retrieving Workspaces from %s : Error %s\n", c.SourceHostname, err)
	}

	dstWS, err := c.DestinationClient.Workspaces.List(c.DestinationContext, c.DestinationOrganizationName, &wsListOptions)
	if err != nil {
		fmt.Printf("Error With retrieving Workspaces from %s : Error %s\n", c.DestinationHostname, err)
	}

	// Create maps to track unique workspace names
	srcUnique := make(map[string]bool)
	dstUnique := make(map[string]bool)

	// Populate maps with workspace names
	for _, ws := range srcWS.Items {
		srcUnique[ws.Name] = true
	}
	for _, ws := range dstWS.Items {
		dstUnique[ws.Name] = true
	}

	// Remove identical workspace names from both maps
	for name := range srcUnique {
		if dstUnique[name] {
			delete(srcUnique, name)
			delete(dstUnique, name)
		}
	}

	// Convert maps back to slices for further processing
	srcWS.Items = nil
	for name := range srcUnique {
		srcWS.Items = append(srcWS.Items, &tfe.Workspace{Name: name})
	}

	dstWS.Items = nil
	for name := range dstUnique {
		dstWS.Items = append(dstWS.Items, &tfe.Workspace{Name: name})
	}

	if len(dstWS.Items) <= 0 && len(srcWS.Items) <= 0 {
		fmt.Println("\nMatched all workspaces names between source and destination.")
		return nil
	}

	// Create a single table with headers for source and destination workspaces
	o.AddTableHeaders("Source Workspaces", "Destination Workspaces")
	maxLen := len(srcWS.Items)
	if len(dstWS.Items) > maxLen {
		maxLen = len(dstWS.Items)
	}

	for i := 0; i < maxLen; i++ {
		var srcName, dstName string
		if i < len(srcWS.Items) {
			srcName = srcWS.Items[i].Name
		}
		if i < len(dstWS.Items) {
			dstName = dstWS.Items[i].Name
		}
		o.AddTableRows(srcName, dstName)
	}

	return nil
}
