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
src_tfe_token=""
dst_tfc_hostname=""
dst_tfc_org=""
dst_tfc_token=""
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
