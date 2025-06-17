// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package copy

import (
	"fmt"
	"os"

	"github.com/hashicorp-services/tfm/cmd/helper"
	"github.com/hashicorp-services/tfm/tfclient"
	tfe "github.com/hashicorp/go-tfe"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (

	// `tfemigrate copy projects` command
	projectsCopyCmd = &cobra.Command{
		Use:     "projects",
		Short:   "Copy projects",
		Aliases: []string{"proj"},
		Long:    "Copy projects from source to destination org",
		//ValidArgs: []string{"state", "vars"},
		//Args:      cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			// Validate `projects-map` if it exists before any other functions can run.
			valid, projMapCfg, err := validateMap(tfclient.GetClientContexts(), "projects-map")
			if err != nil {
				return err
			}

			// Continue the application if `projects-map` is not provided. The valid and map output arent needed.
			_ = valid

			// switch {
			// case vars:
			// 	return copyVariables(tfclient.GetClientContexts())
			// case teamaccess:
			// 	return copyprojTeamAccess(tfclient.GetClientContexts())
			// }

			return copyProjects(
				tfclient.GetClientContexts(), projMapCfg)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			o.Close()
		},
	}
)

func init() {

	// `tfemigrate copy projects --project-id [projectID]`
	projectsCopyCmd.Flags().String("project-id", "", "Specify one single project ID to copy to destination")
	projectsCopyCmd.Flags().BoolVarP(&vars, "vars", "", false, "Copy project variables")
	projectsCopyCmd.Flags().SetInterspersed(false)

	// Add commands
	CopyCmd.AddCommand(projectsCopyCmd)
}

// List all projects in the source organization
func listSrcProjects(c tfclient.ClientContexts) ([]*tfe.Project, error) {
	o.AddMessageUserProvided("\nGetting list of Projects from: ", c.SourceHostname)

	options := tfe.ProjectListOptions{
		ListOptions: tfe.ListOptions{PageSize: 100},
	}

	srcProjects := []*tfe.Project{}

	for {
		projects, err := c.SourceClient.Projects.List(c.SourceContext, c.SourceOrganizationName, &options)
		if err != nil {
			return nil, err
		}

		srcProjects = append(srcProjects, projects.Items...)

		o.AddFormattedMessageCalculated("\nFound %d Projects", len(srcProjects))

		if projects.CurrentPage >= projects.TotalPages {
			break
		}
		options.PageNumber = projects.NextPage
	}
	return srcProjects, nil
}

// List all projects in the destination organization
func listDestProjects(c tfclient.ClientContexts, output bool) ([]*tfe.Project, error) {
	o.AddMessageUserProvided("\nGetting list of Projects from: ", c.DestinationHostname)

	options := tfe.ProjectListOptions{
		ListOptions: tfe.ListOptions{PageSize: 100},
	}

	destProjects := []*tfe.Project{}

	for {
		projects, err := c.DestinationClient.Projects.List(c.DestinationContext, c.DestinationOrganizationName, &options)
		if err != nil {
			return nil, err
		}

		destProjects = append(destProjects, projects.Items...)

		o.AddFormattedMessageCalculated("\nFound %d Projects", len(destProjects))

		if projects.CurrentPage >= projects.TotalPages {
			break
		}
		options.PageNumber = projects.NextPage
	}
	return destProjects, nil
}

// Gets all projects defined in the configuration file `projects` or `projects-map` lists from the source target
func getSrcProjectsCfg(c tfclient.ClientContexts) ([]*tfe.Project, error) {

	var srcProjects []*tfe.Project

	// Get source Project list from config list `projects` if it exists
	srcProjectsCfg := viper.GetStringSlice("projects")

	projMapCfg, err := helper.ViperStringSliceMap("projects-map")
	if err != nil {
		return srcProjects, errors.New("Invalid input for projects-map")
	}

	if len(srcProjectsCfg) > 0 {
		o.AddFormattedMessageCalculated("Found %d projects in `projects` list", len(srcProjectsCfg))
	}

	// If no Projects found in config (list or map), default to just assume all Projects from source will be chosen
	if len(srcProjectsCfg) > 0 && len(projMapCfg) > 0 {
		o.AddErrorUserProvided("'projects' list and 'projects-map' cannot be defined at the same time.")
		os.Exit(0)

	} else if len(projMapCfg) > 0 {

		// use config project from map
		var projList []string

		for key := range projMapCfg {
			projList = append(projList, key)
		}

		o.AddMessageUserProvided("Source Projects found in `projects-map`:", projList)

		// Set source projects
		srcProjects, err = getSrcProjectsFilter(tfclient.GetClientContexts(), projList)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to list Projects in map from source")
		}

	} else if len(srcProjectsCfg) > 0 {
		// use config projects from list

		fmt.Println("Using Projects config list:", srcProjectsCfg)

		//get source Projects
		srcProjects, err = getSrcProjectsFilter(tfclient.GetClientContexts(), srcProjectsCfg)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to list Projects from source")
		}

	} else {
		// Get ALL source Projects
		o.AddMessageUserProvided2("\nWarning:\n\n", "ALL Projects WILL BE MIGRATED from", viper.GetString("src_tfe_hostname"))

		srcProjects, err = listSrcProjects(tfclient.GetClientContexts())
		if !confirm() {
			fmt.Println("\n\n**** Canceling tfm run **** ")
			os.Exit(1)
		}

		if err != nil {
			return nil, errors.Wrap(err, "Failed to list Projects from source")
		}
	}

	// Check Projects exist in source from config
	for _, s := range srcProjectsCfg {
		//fmt.Println("\nFound Project ", s, "in config, check if it exists in", viper.GetString("src_tfe_hostname"))
		exists := doesProjectExist(s, srcProjects)
		if !exists {
			fmt.Printf("Defined Project in config %s DOES NOT exist in %s. \n Please validate your configuration.", s, viper.GetString("src_tfe_hostname"))
			break
		}
	}

	return srcProjects, nil

}

func getSrcProjectsFilter(c tfclient.ClientContexts, projList []string) ([]*tfe.Project, error) {
	o.AddMessageUserProvided("Getting list of Projects from: ", c.SourceHostname)
	srcProjects := []*tfe.Project{}

	for _, proj := range projList {

		var found bool

		for {
			opts := tfe.ProjectListOptions{
				ListOptions: tfe.ListOptions{
					PageNumber: 1,
					PageSize:   100},
				Name: proj,
			}

			items, err := c.SourceClient.Projects.List(c.SourceContext, c.SourceOrganizationName, &opts) // This should only return 1 result

			if err != nil {
				return nil, err
			}

			indexMatch := -1

			// If multiple Projects named similar, find exact match
			for i, result := range items.Items {
				if proj == result.Name {
					indexMatch = i
					found = true
					break
				}
			}

			// Append only if a matching project is found
			if found && indexMatch >= 0 {
				srcProjects = append(srcProjects, items.Items[indexMatch])
				break // Break the inner loop if a match is found
			}

			// If no project is found for a given name in projList, handle it accordingly
			if !found {
				// You might want to log this or handle it as per your logic
				o.AddMessageUserProvided("no project found with name: ", proj)
			}

			if items.CurrentPage >= items.TotalPages {
				break
			}
			opts.PageNumber = items.NextPage

		}
	}

	return srcProjects, nil
}

func getDstProjectsFilter(c tfclient.ClientContexts, projList []string) ([]*tfe.Project, error) {
	o.AddMessageUserProvided("Getting list of Projects from: ", c.DestinationHostname)
	dstProjects := []*tfe.Project{}

	fmt.Println("Project list from config:", projList)

	for _, proj := range projList {

		for {
			opts := tfe.ProjectListOptions{
				ListOptions: tfe.ListOptions{
					PageNumber: 1,
					PageSize:   100},
				Name: proj,
			}

			// This should only return 1 result
			items, err := c.DestinationClient.Projects.List(c.DestinationContext, c.DestinationOrganizationName, &opts)
			if err != nil {
				return nil, err
			}

			indexMatch := 0

			// If multiple Projects named similar, find exact match
			if len(items.Items) > 1 {
				for _, result := range items.Items {
					if proj == result.Name {
						// Finding matching Project name
						break
					}
					indexMatch++
				}
			}

			dstProjects = append(dstProjects, items.Items[indexMatch])

			if items.CurrentPage >= items.TotalPages {
				break
			}
			opts.PageNumber = items.NextPage

		}
	}

	return dstProjects, nil
}

func doesProjectExist(projectName string, proj []*tfe.Project) bool {
	for _, w := range proj {
		if projectName == w.Name {
			return true
		}
	}
	return false
}

// copyProject creates a new project in the target organization with settings from the source project.
func copyProjects(c tfclient.ClientContexts, projMapCfg map[string]string) error {

	// Get Projects from Config OR get ALL Projects from source
	srcProjects, err := getSrcProjectsCfg(c)
	if err != nil {
		return errors.Wrap(err, "Failed to list Projects from source target")
	}

	// Get Projects from Config OR get ALL Projects from source
	destProjects, err := listDestProjects(tfclient.GetClientContexts(), false)
	if err != nil {
		return errors.Wrap(err, "Failed to list Projects from destination target")
	}

	for _, srcproject := range srcProjects {
		destProjectName := srcproject.Name

		// Check if the destination Project name differs from the source name
		if len(projMapCfg) > 0 {
			o.AddMessageUserProvided3("Source Project:", srcproject.Name, "\nDestination Project:", projMapCfg[srcproject.Name])
			destProjectName = projMapCfg[srcproject.Name]
		}

		exists := doesProjectExist(destProjectName, destProjects)

		if exists {
			o.AddMessageUserProvided2(destProjectName, "exists in destination will not migrate", srcproject.Name)
		} else {

			srcproject, err := createProject(c, destProjectName)

			if err != nil {
				fmt.Println("Could not create Project.\n\n Error:", err.Error())
				return err
			}
			o.AddDeferredMessageRead("Migrated", srcproject.Name)
		}
	}
	return nil
}

func createProject(c tfclient.ClientContexts, projectName string) (*tfe.Project, error) {
	o.AddMessageUserProvided("Creating Project in destination: ", projectName)

	// Create the project in the destination organization
	project, err := c.DestinationClient.Projects.Create(c.DestinationContext, c.DestinationOrganizationName, tfe.ProjectCreateOptions{
		Name:        projectName,
		Description: helper.ViperString("created by tfm"),
	})

	if err != nil {
		return nil, errors.Wrap(err, "Failed to create Project in destination")
	}

	return project, nil
}
