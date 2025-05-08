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
	suffix  string
	prefix  string

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
	wsCmpCmd.PersistentFlags().StringVar(&srcID, "src-id", "Default Project", "partial workspace name used to filter the results")
	wsCmpCmd.MarkFlagRequired("src-id")
	wsCmpCmd.PersistentFlags().StringVar(&srcType, "src-type", "project", "specify the source type (organization or project)")
	wsCmpCmd.MarkFlagRequired("src-type")
	wsCmpCmd.PersistentFlags().StringVar(&dstID, "dst-id", "Default Project", "comma-separated tag names to exclude")
	wsCmpCmd.MarkFlagRequired("dst-id")
	wsCmpCmd.PersistentFlags().StringVar(&dstType, "dst-type", "project", "specify the destination type (organization or project)")
	wsCmpCmd.MarkFlagRequired("dst-type")
	wsCmpCmd.PersistentFlags().StringVar(&suffix, "suffix-filter", "", "(optional) only for destination workspaces, if they were copied over with a common suffix, it will remove them for the comparison")
	wsCmpCmd.PersistentFlags().StringVar(&prefix, "prefix-filter", "", "(optional) only for destination workspaces, if they were copied over with a common prefix, it will remove them for the comparison")
	wsCmpCmd.PersistentFlags().Args()
}

func wsCmp(c tfclient.ClientContexts) error {

	// Validate the source and destination types
	// They must be either "organization" or "project"
	if srcType != "organization" && srcType != "project" {
		fmt.Println("Invalid source type. Must be 'organization' or 'project'.")
		return nil
	}
	if dstType != "organization" && dstType != "project" {
		fmt.Println("Invalid destination type. Must be 'organization' or 'project'.")
		return nil
	}

	// Checks if source and destionation types are the project
	// and if the IDs are valid
	// Project IDs must start with "prj-"
	if srcType == "project" && len(srcID) >= 4 && srcID[:4] != "prj-" {
		fmt.Println("Invalid source ID for project. Project IDs must start with 'prj-'.")
		return nil
	}

	if dstType == "project" && len(dstID) >= 4 && dstID[:4] != "prj-" {
		fmt.Println("Invalid destination ID for project. Project IDs must start with 'prj-'.")
		return nil
	}

	var wsSRCListOptions tfe.WorkspaceListOptions
	var wsDSTListOptions tfe.WorkspaceListOptions

	// Using a helper function to create WorkspaceListOptions
	createWorkspaceListOptions := func(id, typ string) tfe.WorkspaceListOptions {
		options := tfe.WorkspaceListOptions{
			ListOptions: tfe.ListOptions{
				PageNumber: 1,
				PageSize:   100,
			},
		}
		if typ == "project" {
			options.ProjectID = id
		}
		return options
	}

	// Initialize options for source and destination
	wsSRCListOptions = createWorkspaceListOptions(srcID, srcType)
	wsDSTListOptions = createWorkspaceListOptions(dstID, dstType)

	// This was for debugging purposes, you can uncomment it if needed.
	// fmt.Printf("WorkspaceListOptions: %+v\n", wsSRCListOptions)
	// fmt.Printf("WorkspaceListOptions: %+v\n", wsDSTListOptions)

	srcWS, err := c.SourceClient.Workspaces.List(c.SourceContext, c.SourceOrganizationName, &wsSRCListOptions)

	if err != nil {
		fmt.Printf("Error With retrieving Workspaces from %s : Error %s\n", c.SourceHostname, err)
	}

	dstWS, err := c.DestinationClient.Workspaces.List(c.DestinationContext, c.DestinationOrganizationName, &wsDSTListOptions)

	if err != nil {
		fmt.Printf("Error With retrieving Workspaces from %s : Error %s\n", c.DestinationHostname, err)
	}

	// List of all workspaces

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

	// Remove suffix and prefix from destination workspace names if specified
	if suffix != "" || prefix != "" {
		for name := range dstUnique {
			originalName := name
			if suffix != "" && len(name) > len(suffix)+1 && name[len(name)-len(suffix)-1:] == "-"+suffix {
				name = name[:len(name)-len(suffix)-1]
			}
			if prefix != "" && len(name) > len(prefix)+1 && name[:len(prefix)+1] == prefix+"-" {
				name = name[len(prefix)+1:]
			}
			if name != originalName {
				delete(dstUnique, originalName)
				dstUnique[name] = true
			}
		}
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

	o.AddMessageUserProvided("Source Workspace Count:", len(srcWS.Items))
	o.AddMessageUserProvided("Destination Workspaces Count:", len(dstWS.Items))

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
