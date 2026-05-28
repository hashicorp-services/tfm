// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package logging provides the global structured logger for tfm.
//
// The design mirrors hashicorp/terraform/internal/logging: a single
// go-hclog InterceptLogger is initialised once (via Init) and all packages
// obtain named child loggers via NewLogger.
//
// Log level is controlled by:
//
//  1. TFM_LOG environment variable (TRACE|DEBUG|INFO|WARN|ERROR|OFF|JSON)
//  2. --verbose / -V CLI flag  →  effective level INFO
//
// The environment variable takes precedence over the CLI flag.
// TFM_LOG_PATH, when set, redirects all log output to the named file.
package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/hashicorp/go-hclog"
)

// Environment variable names, following the TF_LOG convention.
const (
	EnvLog     = "TFM_LOG"
	EnvLogPath = "TFM_LOG_PATH"
)

// ValidLevels lists the accepted values for TFM_LOG.
var ValidLevels = []string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "OFF"}

var (
	mu        sync.Mutex
	logger    hclog.InterceptLogger
	logWriter io.Writer
)

func init() {
	// Bootstrap a silent logger so callers never receive a nil logger
	// even if Init has not been called yet.
	logger = hclog.NewInterceptLogger(&hclog.LoggerOptions{
		Name:   "tfm",
		Level:  hclog.Off,
		Output: io.Discard,
	})
	logWriter = io.Discard
}

// Init initialises the global logger from the environment and the --verbose
// flag value. It must be called once, early in cobra's initialisation chain
// (inside cobra.OnInitialize or initConfig), before any commands run.
//
// Level precedence: TFM_LOG env var > verboseFlag > default (OFF).
func Init(verboseFlag bool) {
	mu.Lock()
	defer mu.Unlock()

	level, jsonFmt := resolveLevel(verboseFlag)
	output := resolveOutput()

	logger = hclog.NewInterceptLogger(&hclog.LoggerOptions{
		Name:              "tfm",
		Level:             level,
		Output:            output,
		IndependentLevels: true,
		JSONFormat:        jsonFmt,
	})

	logWriter = logger.StandardWriter(&hclog.StandardLoggerOptions{InferLevels: true})

	// Redirect the stdlib log package to hclog so any legacy log.Printf calls
	// are captured.
	log.SetFlags(0)
	log.SetPrefix("")
	log.SetOutput(logWriter)
}

// NewLogger returns a named child of the global logger. The name is appended
// to "tfm" (e.g. "tfm.copy.workspaces").
func NewLogger(name string) hclog.Logger {
	mu.Lock()
	defer mu.Unlock()
	if name == "" {
		panic("logging.NewLogger: name must not be empty")
	}
	return logger.Named(name)
}

// CurrentLevel returns the string representation of the active log level.
func CurrentLevel() string {
	mu.Lock()
	defer mu.Unlock()
	return strings.ToUpper(logger.GetLevel().String())
}

// resolveLevel determines the hclog.Level and JSON flag from the environment
// and the --verbose flag.
func resolveLevel(verboseFlag bool) (hclog.Level, bool) {
	envVal := strings.ToUpper(strings.TrimSpace(os.Getenv(EnvLog)))

	// JSON is a special pseudo-level: emit JSON at TRACE verbosity.
	if envVal == "JSON" {
		return hclog.Trace, true
	}

	if envVal != "" {
		if isValidLevel(envVal) {
			return hclog.LevelFromString(envVal), false
		}
		fmt.Fprintf(os.Stderr,
			"[WARN] Invalid TFM_LOG value: %q. Defaulting to OFF. Valid levels: %s\n",
			envVal, strings.Join(ValidLevels, " "))
		return hclog.Off, false
	}

	// No env var — check the --verbose flag.
	if verboseFlag {
		return hclog.Info, false
	}

	return hclog.Off, false
}

// resolveOutput opens TFM_LOG_PATH for append-write if set, otherwise returns
// stderr.
func resolveOutput() io.Writer {
	if path := os.Getenv(EnvLogPath); path != "" {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[WARN] TFM_LOG_PATH: could not open %q: %v\n", path, err)
			return os.Stderr
		}
		return f
	}
	return os.Stderr
}

func isValidLevel(level string) bool {
	for _, l := range ValidLevels {
		if level == l {
			return true
		}
	}
	return false
}
