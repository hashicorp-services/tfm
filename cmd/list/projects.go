// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package list

import (
	"context"
	"fmt"

	"encoding/json"

	"github.com/hashicorp-services/tfm/cmd/logging"
	"github.com/hashicorp-services/tfm/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

var (

	// `tfm list project` command
	projectListCmd = &cobra.Command{
		Use:     "projects",
		Aliases: []string{"prj"},
		Short:   "Projects command",
		Long:    "List Projects in an org",
		Run: func(cmd *cobra.Command, args []string) {
			listProjects(tfclient.GetClientContexts(), jsonOut)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			o.Close()
		},
	}
)

func init() {

	// Add commands
	ListCmd.AddCommand(projectListCmd)

}

func listProjects(c tfclient.ClientContexts, jsonOut bool) error {
	log := logging.NewLogger("list.projects")

	srcProjects := []*tfe.Project{}
	projectJSON := make(map[string]interface{}) // Parent JSON object "project-names"
	projectNamesAndIDs := []map[string]string{} // project names slice to go inside parent object "project-names"

	opts := tfe.ProjectListOptions{
		ListOptions: tfe.ListOptions{
			PageNumber: 1,
			PageSize:   100},
	}

	if (ListCmd.Flags().Lookup("side").Value.String() == "source") || (!ListCmd.Flags().Lookup("side").Changed) {

		log.Info("listing projects", "org", c.SourceOrganizationName, "host", c.SourceHostname)

		if jsonOut == false {
			o.AddMessageUserProvided("Getting list of projects from: ", c.SourceHostname)
		}

		for {
			items, err := c.SourceClient.Projects.List(c.SourceContext, c.SourceOrganizationName, &opts)
			if err != nil {
				log.Error("failed to list projects", "org", c.SourceOrganizationName, "error", err)
				return err
			}

			srcProjects = append(srcProjects, items.Items...)
			log.Debug("fetched project page", "page", items.CurrentPage, "total_pages", items.TotalPages, "count", len(srcProjects))

			if jsonOut == false {
				o.AddFormattedMessageCalculated("Found %d Projects", len(srcProjects))
			}

			if items.CurrentPage >= items.TotalPages {
				break
			}
			opts.PageNumber = items.NextPage

		}
		if jsonOut == false {
			o.AddTableHeaders("Name", "ID")
		}
		log.Info("found projects", "org", c.SourceOrganizationName, "count", len(srcProjects))
		for _, i := range srcProjects {
			projectInfo := map[string]string{
				"name": i.Name,
				"id":   i.ID,
			}

			if jsonOut {
				projectNamesAndIDs = append(projectNamesAndIDs, projectInfo) // Store project name in slice
			}
			if jsonOut == false {
				o.AddTableRows(i.Name, i.ID)
			}
		}
		if jsonOut {
			projectJSON["projects"] = projectNamesAndIDs // Assign projects to the "projects" key

			jsonData, err := json.Marshal(projectJSON)
			if err != nil {
				fmt.Println("Error marshaling projects to JSON:", err)
				return err
			}

			fmt.Println(string(jsonData))
		}
	}

	if ListCmd.Flags().Lookup("side").Value.String() == "destination" {
		log.Info("listing projects", "org", c.DestinationOrganizationName, "host", c.DestinationHostname)

		if jsonOut == false {
			o.AddMessageUserProvided("Getting list of projects from: ", c.DestinationHostname)
		}

		for {
			items, err := c.DestinationClient.Projects.List(c.DestinationContext, c.DestinationOrganizationName, &opts)
			if err != nil {
				log.Error("failed to list destination projects", "org", c.DestinationOrganizationName, "error", err)
				return err
			}

			srcProjects = append(srcProjects, items.Items...)
			log.Debug("fetched destination project page", "page", items.CurrentPage, "total_pages", items.TotalPages, "count", len(srcProjects))

			if jsonOut == false {
				o.AddFormattedMessageCalculated("Found %d Projects", len(srcProjects))
			}

			if items.CurrentPage >= items.TotalPages {
				break
			}
			opts.PageNumber = items.NextPage

		}
		if jsonOut == false {
			o.AddTableHeaders("Name", "ID")
		}
		log.Info("found projects", "org", c.DestinationOrganizationName, "count", len(srcProjects))

		for _, i := range srcProjects {
			projectInfo := map[string]string{
				"name": i.Name,
				"id":   i.ID,
			}

			if jsonOut {
				projectNamesAndIDs = append(projectNamesAndIDs, projectInfo) // Store project name in slice
			}

			if jsonOut == false {
				o.AddTableRows(i.Name, i.ID)
			}
		}
		if jsonOut {
			projectJSON["projects"] = projectNamesAndIDs // Assign projects to the "project" key

			jsonData, err := json.Marshal(projectJSON)
			if err != nil {
				fmt.Println("Error marshaling projects to JSON:", err)
				return err
			}

			fmt.Println(string(jsonData))
		}
	}

	return nil
}

func getProjectName(client *tfe.Client, ctx context.Context, projectId string) (string, error) {

	prj, err := client.Projects.Read(ctx, projectId)

	if err != nil {
		return "error reading project", err
	}

	return prj.Name, nil
}
