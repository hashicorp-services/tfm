// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package generate

import (
	"html/template"
	"os"

	"github.com/hashicorp-services/tfm/output"
	"github.com/spf13/cobra"
)

var (
	o output.Output

	// `tfm generate config` command
	generateConfigCmd = &cobra.Command{
		Use:     "config",
		Aliases: []string{"cfg"},
		Short:   "config command",
		Long:    "Generate a .tfm.hcl template",
		Run: func(cmd *cobra.Command, args []string) {
			generateConfigTemplate()
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			o.Close()
		},
	}
)

func init() {

	// Add commands
	GenerateCmd.AddCommand(generateConfigCmd)

}

// Template contents
const templateContent = `
{{.Description}}

src_tfe_hostname=""
src_tfe_org=""
src_tfe_token="Must have owner permissions"
dst_tfc_hostname=""
dst_tfc_org=""
dst_tfc_token="Must have owner permissions"
#dst_tfc_project_id=""

# A list of source=destination VCS oauth IDs. TFM will look at each workspace in the source for the source VCS oauth ID and assign the matching workspace in the destination with the destination VCS oauth ID.
#vcs-map=[
#  "ot-wF6KZMna4desiPRc=ot-JSQTcnWxqVL5zQ1w",
#]


# A List of Workspaces to create/check are migrated across to new TFC
#"workspaces" = [
#  "example-ws-1",
#  "example-ws-2"
#]

# A list of source=destination workspace names. TFM will look at each source workspace and recreate the workspace with the specified destination name.
#"workspaces-map" = [
#  "example-ws-1=new-ws-1",
#  "example-ws-2=new-ws-2"
#    ]

# A List of Projects to create/check are migrated across to new TFC
#"projects" = [
#  "example-proj-1",
#  "example-proj-2"
#]

# A list of source=destination project names. TFM will look at each source project and recreate the project with the specified destination name.
#"projects-map" = [
#  "example-proj-1=new-proj-1",
#  "example-proj-2=new-proj-2"
#    ]

# A list of source=destination agent pool IDs TFM will look at each workspace in the source for the source agent pool ID and assign the matching workspace in the destination the destination agent pool ID. Conflicts with 'agent-assignment'
#agents-map = [
#  "apool-DgzkahoomwHsBHcJ=apool-vbrJZKLnPy6aLVxE",
#  "apool-DgzkahoomwHsBHc3=apool-vbrJZKLnPy6aLVx4",
#]

# An agent Pool ID to assign to all workspaces in the destination. Conflicts with 'workspaces-map'
#agent-assignment-id="apool-h896pi2MeP4JJvsB"

# A list of source=destination variable set names. TFM will look at each source variable set and recreate the variable set with the specified destination name.
#varsets-map = [
#  "Azure-creds=New-Azure-Creds",
#  "aws-creds2=New-AWS-Creds"
# ]

# A list of source=destination SSH IDs. TFM will look at each workspace in the source for the source SSH  ID and assign the matching workspace in the destination with the destination SSH ID.
#ssh-map=[
#  "sshkey-sPLAKMcqnWtHPSgx=sshkey-CRLmPJpoHwsNFAoN"
#]

# THE FOLLOWING ARE ONLY USED FOR MIGRATING FROM TERRAFORM OPEN SOURCE / COMMUNITY EDITION TO TFE/TFC

#commit_message = "A commit message the tfm core remove-backend command uses when removing backend blocks from .tf files and commiting the changes back"
#commit_author_name = "the name that will appear as the commit author"
#commit_author_email = "the email that will appear for the commit author"
#vcs_type = "github"
#gitlab_token = "A gitlab token."
#gitlab_group = "gitlab group ID (This is usually found in the url bar of gitlab)."
#gitlab_username = "A gitlab username"
#github_token = "A github token used with the tfm core clone command. Must have read permissions to clone repos and write permissions to use the remove-backend command"
#github_organization = "The github Organization to clone repos from"
#github_username = "A github username"

#clone_repos_path = "/path/on/local/host/to/clone/repos/to"
#vcs_provider_id = "An Oauth ID of a VCS provider connection configured in TFC/TFE"

# A list of VCS repositories containing terraform code. TFM will clone each repo during the tfm core clone command for migrating opensource/commmunity edition terraform managed code to TFE/TFC.
# If one is not provided then tfm will attempt to clone every repo it has read access to in the Github Organization.
#repos_to_clone =  [
# "repo1",
# "repo2",
# "repo3"
#]
`

type TemplateData struct {
	Description string
}

func generateConfigTemplate() {
	tmpl, err := template.New(".tfm.hcl").Parse(templateContent)
	if err != nil {
		panic(err)
	}

	outputFile, err := os.Create(".tfm.hcl")
	if err != nil {
		panic(err)
	}

	defer outputFile.Close()

	data := TemplateData{
		Description: "# TFM Config File",
	}

	err = tmpl.Execute(outputFile, data)
	if err != nil {
		panic(err)
	}

	o.AddMessageUserProvided(".tfm.hcl template generated in current directroy ", "")
}
