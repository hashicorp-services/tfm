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

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var (
	// `tfe-migrate foo` commands
	fooCmd = &cobra.Command{
		Use:   "foo",
		Short: "Foo commands",
		Long:  "The Foo commands for tfe-migrate foo create/list/show/delete",
	}
	// `tfe-migrate foo create` command
	fooCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create Foo",
		Long:  "Create Foo things in the system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fooCreate(
				*viperString("name"))
		},
	}

	// `tfe-migrate foo list` command
	fooListCmd = &cobra.Command{
		Use:   "list",
		Short: "List Foo",
		Long:  "List Foo things in the system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fooList(
				*viperBool("all"))
		},
	}

	// `tfe-migrate foo show` command
	fooShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show Foo",
		Long:  "Show Foo things in the system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fooShow(
				*viperInt("id"))
		},
	}

	// `tfe-migrate foo delete` command
	fooDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete Foo",
		Long:  "Delete Foo things in the system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fooDelete(
				*viperInt("id"))
		},
	}
)

func init() {
	// Flags().StringP, etc... - the "P" gives us the option for a short hand

	// `tfe-migrate foo create` command
	fooCreateCmd.Flags().StringP("name", "n", "", "Name of foo.")
	fooCreateCmd.MarkFlagRequired("name")

	// `tfe-migrate foo list` command
	fooListCmd.Flags().BoolP("all", "a", false, "List all? (optional)")

	// `tfe-migrate foo show` command
	fooShowCmd.Flags().Int16P("id", "i", 0, "id of foo.")
	fooShowCmd.MarkFlagRequired("id")

	// `tfe-migrate foo delete` command
	fooDeleteCmd.Flags().Int16P("id", "i", 0, "id of foo.")
	fooDeleteCmd.MarkFlagRequired("id")

	// Add commands
	rootCmd.AddCommand(fooCmd)
	fooCmd.AddCommand(fooCreateCmd)
	fooCmd.AddCommand(fooListCmd)
	fooCmd.AddCommand(fooShowCmd)
	fooCmd.AddCommand(fooDeleteCmd)
}

func fooCreate(name string) error {
	fmt.Println("Create foo with name:", aurora.Bold(name))
	return nil
}

func fooList(all bool) error {
	fmt.Println("List foo with all:", aurora.Bold(all))
	return nil
}

func fooShow(id int) error {
	fmt.Println("Delete foo with id:", aurora.Bold(id))
	return nil
}

func fooDelete(id int) error {
	fmt.Println("Delete foo with id:", aurora.Bold(id))
	return nil
}
