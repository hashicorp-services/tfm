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
	"log"

	"github.com/hashicorp-services/tfm/cmd/copy"
	"github.com/hashicorp-services/tfm/cmd/core"
	"github.com/hashicorp-services/tfm/cmd/delete"
	"github.com/hashicorp-services/tfm/cmd/generate"
	"github.com/hashicorp-services/tfm/cmd/helper"
	"github.com/hashicorp-services/tfm/cmd/list"
	"github.com/hashicorp-services/tfm/cmd/lock"
	// "github.com/hashicorp-services/tfm/cmd/nuke"
	"github.com/hashicorp-services/tfm/cmd/unlock"
	"github.com/hashicorp-services/tfm/output"
	"github.com/hashicorp-services/tfm/version"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	o       *output.Output
	jsonOut bool

	// Required to leverage viper defaults for optional Flags
	bindPFlags = func(cmd *cobra.Command, args []string) {
		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			log.Fatal(aurora.Red(err))
		}
	}
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:              "tfm",
	Short:            "A CLI to assist with Terraform community edition, Terraform Cloud, and Terraform Enterprise migrations.",
	SilenceUsage:     true,
	SilenceErrors:    true,
	Version:          version.String(),
	PersistentPreRun: bindPFlags, // Bind here to avoid having to call this in every subcommand
}

// `tfemig copy` commands
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy command",
	Long:  "Copy objects from source TFC/TFE org to destination TFC/TFE org",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// // Close output stream always before exiting
	if err := RootCmd.Execute(); err != nil {
		o.Close()
		log.Fatal(aurora.Red(err))
	} else {
		o.Close()
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file, can be used to store common flags, (default is ~/.tfm.hcl).")
	RootCmd.PersistentFlags().BoolP("autoapprove", "", false, "Auto approve the tfm run. --autoapprove=true . false by default")
	RootCmd.PersistentFlags().BoolVar(&jsonOut, "json", false, "Print the output in JSON format")

	// Available commands required after "tfm"
	RootCmd.AddCommand(copy.CopyCmd)
	RootCmd.AddCommand(list.ListCmd)
	// RootCmd.AddCommand(nuke.NukeCmd)
	RootCmd.AddCommand(delete.DeleteCmd)
	RootCmd.AddCommand(generate.GenerateCmd)
	RootCmd.AddCommand(lock.LockCmd)
	RootCmd.AddCommand(unlock.UnlockCmd)
	RootCmd.AddCommand(core.CoreCmd)
	// Turn off completion option
	RootCmd.CompletionOptions.DisableDefaultCmd = true

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in current & home directory with name ".tfm" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("hcl")
		viper.AddConfigPath(".")
		viper.SetConfigName(".tfm.hcl")

	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	isConfigFile := false
	if err := viper.ReadInConfig(); err == nil {
		isConfigFile = true // Capture information here to bring after all flags are loaded (namely which output type)
	}

	// Some hacking here to let viper use the cobra required flags, simplifies this checking
	// in one place rather than each command
	// More info: https://github.com/spf13/viper/issues/397
	postInitCommands(RootCmd.Commands())

	// // Initialize output
	o = output.New(*helper.ViperBool("json"))

	// check to see if the --json flag was provided and return bool value assigned to "json"
	json := viper.IsSet("json")

	// Print if config file was found and json output is desired
	if isConfigFile && !json {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// copy.pasta function
func postInitCommands(commands []*cobra.Command) {
	for _, cmd := range commands {
		presetRequiredFlags(cmd)
		if cmd.HasSubCommands() {
			postInitCommands(cmd.Commands())
		}
	}
}

// copy.pasta function
func presetRequiredFlags(cmd *cobra.Command) {
	viper.BindPFlags(cmd.Flags())
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			cmd.Flags().Set(f.Name, viper.GetString(f.Name))
		}
	})
}
