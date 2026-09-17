// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cmd_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/hashicorp-services/tfm/cmd"
)

// TestCompletionCommandRegistered verifies that the built-in Cobra completion
// command is accessible on RootCmd (i.e. DisableDefaultCmd is not set).
// The completion command is registered lazily on the first Execute call, so
// we trigger that call here before inspecting the command tree.
func TestCompletionCommandRegistered(t *testing.T) {
	// Trigger lazy registration.
	var discard bytes.Buffer
	cmd.RootCmd.SetOut(&discard)
	cmd.RootCmd.SetErr(&discard)
	cmd.RootCmd.SetArgs([]string{"completion", "bash"})
	_ = cmd.RootCmd.Execute()
	cmd.RootCmd.SetArgs(nil)

	completionCmd, _, err := cmd.RootCmd.Find([]string{"completion"})
	if err != nil {
		t.Fatalf("Find(completion): unexpected error: %v", err)
	}
	if completionCmd == nil {
		t.Fatal("completion command is nil; expected it to be registered on RootCmd")
	}
	if completionCmd.Name() != "completion" {
		t.Errorf("expected command name %q, got %q", "completion", completionCmd.Name())
	}
}

// TestCompletionSubcommands checks that all four shell sub-commands are
// present under the completion command.
func TestCompletionSubcommands(t *testing.T) {
	// Ensure the completion parent is registered.
	var discard bytes.Buffer
	cmd.RootCmd.SetOut(&discard)
	cmd.RootCmd.SetErr(&discard)
	cmd.RootCmd.SetArgs([]string{"completion", "bash"})
	_ = cmd.RootCmd.Execute()
	cmd.RootCmd.SetArgs(nil)

	completionCmd, _, err := cmd.RootCmd.Find([]string{"completion"})
	if err != nil || completionCmd == nil {
		t.Fatalf("could not find completion command: %v", err)
	}

	shells := []string{"bash", "zsh", "fish", "powershell"}
	for _, shell := range shells {
		t.Run(shell, func(t *testing.T) {
			shellCmd, _, findErr := completionCmd.Find([]string{shell})
			if findErr != nil {
				t.Fatalf("Find(%s): unexpected error: %v", shell, findErr)
			}
			if shellCmd == nil || shellCmd.Name() != shell {
				t.Errorf("expected sub-command %q to be registered under completion", shell)
			}
		})
	}
}

// TestCompletionBashOutput verifies that GenBashCompletionV2 produces a
// non-empty script that mentions the root command name.
func TestCompletionBashOutput(t *testing.T) {
	var buf bytes.Buffer
	if err := cmd.RootCmd.GenBashCompletionV2(&buf, true); err != nil {
		t.Fatalf("GenBashCompletionV2: %v", err)
	}
	out := buf.String()
	if len(out) == 0 {
		t.Fatal("bash completion script is empty")
	}
	if !strings.Contains(out, "bash") {
		t.Errorf("expected bash completion script content; first 200 chars: %.200s", out)
	}
}

// TestCompletionZshOutput verifies that GenZshCompletion produces a
// non-empty zsh completion script.
func TestCompletionZshOutput(t *testing.T) {
	var buf bytes.Buffer
	if err := cmd.RootCmd.GenZshCompletion(&buf); err != nil {
		t.Fatalf("GenZshCompletion: %v", err)
	}
	out := buf.String()
	if len(out) == 0 {
		t.Fatal("zsh completion script is empty")
	}
	if !strings.Contains(out, "zsh") && !strings.Contains(out, "compdef") {
		t.Errorf("expected zsh completion script content; first 200 chars: %.200s", out)
	}
}

// TestCompletionFishOutput verifies that GenFishCompletion produces a
// non-empty fish completion script.
func TestCompletionFishOutput(t *testing.T) {
	var buf bytes.Buffer
	if err := cmd.RootCmd.GenFishCompletion(&buf, true); err != nil {
		t.Fatalf("GenFishCompletion: %v", err)
	}
	out := buf.String()
	if len(out) == 0 {
		t.Fatal("fish completion script is empty")
	}
	if !strings.Contains(out, "fish") && !strings.Contains(out, "complete") {
		t.Errorf("expected fish completion script content; first 200 chars: %.200s", out)
	}
}

// TestCompletionPowershellOutput verifies that GenPowerShellCompletion produces
// a non-empty PowerShell completion script.
func TestCompletionPowershellOutput(t *testing.T) {
	var buf bytes.Buffer
	if err := cmd.RootCmd.GenPowerShellCompletion(&buf); err != nil {
		t.Fatalf("GenPowerShellCompletion: %v", err)
	}
	out := buf.String()
	if len(out) == 0 {
		t.Fatal("PowerShell completion script is empty")
	}
	if !strings.Contains(out, "Register-ArgumentCompleter") {
		t.Errorf("expected PowerShell completion script content; first 200 chars: %.200s", out)
	}
}
