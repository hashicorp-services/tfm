// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helper

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/hashicorp-services/tfm/cmd/logging"
)

// LogError logs err at ERROR level via the structured logger, prints a
// coloured message to stdout, then calls log.Fatalln to exit.
func LogError(err error, message string) {
	logging.NewLogger("helper").Error(message, "error", err)
	fmt.Println()
	fmt.Println()
	fmt.Println(color.RedString("Error: " + message))
	log.Fatalln(err)
}

// LogWarning logs err at WARN level via the structured logger and prints a
// coloured warning to stdout. Execution continues after this call.
func LogWarning(err error, message string) {
	logging.NewLogger("helper").Warn(message, "error", err)
	fmt.Println()
	fmt.Println()
	fmt.Println(color.YellowString("Warning: " + message))
}
