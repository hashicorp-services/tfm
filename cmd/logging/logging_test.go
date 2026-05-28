// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logging_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp-services/tfm/cmd/logging"
)

// helper resets the package state around each test so env vars don't leak.
func resetEnv(t *testing.T) {
	t.Helper()
	orig := map[string]string{
		logging.EnvLog:     os.Getenv(logging.EnvLog),
		logging.EnvLogPath: os.Getenv(logging.EnvLogPath),
	}
	t.Cleanup(func() {
		for k, v := range orig {
			if v == "" {
				os.Unsetenv(k)
			} else {
				os.Setenv(k, v)
			}
		}
	})
	os.Unsetenv(logging.EnvLog)
	os.Unsetenv(logging.EnvLogPath)
}

func TestInitOffByDefault(t *testing.T) {
	resetEnv(t)
	logging.Init(false)
	if got := logging.CurrentLevel(); got != "OFF" {
		t.Errorf("expected OFF, got %s", got)
	}
}

func TestInitFromEnvDebug(t *testing.T) {
	resetEnv(t)
	os.Setenv(logging.EnvLog, "DEBUG")
	logging.Init(false)
	if got := logging.CurrentLevel(); got != "DEBUG" {
		t.Errorf("expected DEBUG, got %s", got)
	}
}

func TestInitFromEnvInfo(t *testing.T) {
	resetEnv(t)
	os.Setenv(logging.EnvLog, "info") // case-insensitive
	logging.Init(false)
	if got := logging.CurrentLevel(); got != "INFO" {
		t.Errorf("expected INFO, got %s", got)
	}
}

func TestInitVerboseFlag(t *testing.T) {
	resetEnv(t)
	logging.Init(true)
	if got := logging.CurrentLevel(); got != "INFO" {
		t.Errorf("expected INFO, got %s", got)
	}
}

func TestInitEnvWinsOverFlag(t *testing.T) {
	resetEnv(t)
	// TFM_LOG=WARN and --verbose both set: env var should win (WARN > INFO).
	os.Setenv(logging.EnvLog, "WARN")
	logging.Init(true)
	if got := logging.CurrentLevel(); got != "WARN" {
		t.Errorf("expected WARN (env wins over verbose flag), got %s", got)
	}
}

func TestInitInvalidLevelDefaultsOff(t *testing.T) {
	resetEnv(t)
	os.Setenv(logging.EnvLog, "BANANA")
	logging.Init(false)
	if got := logging.CurrentLevel(); got != "OFF" {
		t.Errorf("expected OFF for invalid level, got %s", got)
	}
}

func TestInitTraceLevel(t *testing.T) {
	resetEnv(t)
	os.Setenv(logging.EnvLog, "TRACE")
	logging.Init(false)
	if got := logging.CurrentLevel(); got != "TRACE" {
		t.Errorf("expected TRACE, got %s", got)
	}
}

func TestInitErrorLevel(t *testing.T) {
	resetEnv(t)
	os.Setenv(logging.EnvLog, "ERROR")
	logging.Init(false)
	if got := logging.CurrentLevel(); got != "ERROR" {
		t.Errorf("expected ERROR, got %s", got)
	}
}

func TestInitJSONMode(t *testing.T) {
	resetEnv(t)
	os.Setenv(logging.EnvLog, "JSON")
	logging.Init(false)
	// JSON mode implies TRACE level.
	if got := logging.CurrentLevel(); got != "TRACE" {
		t.Errorf("expected TRACE for JSON mode, got %s", got)
	}
}

func TestNewLogger(t *testing.T) {
	resetEnv(t)
	logging.Init(false)
	l := logging.NewLogger("test.subsystem")
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
	// Verify the name is embedded (hclog logger Name returns root+named).
	if !strings.Contains(l.Name(), "test.subsystem") {
		t.Errorf("expected logger name to contain 'test.subsystem', got %q", l.Name())
	}
}

func TestNewLoggerPanicsOnEmpty(t *testing.T) {
	resetEnv(t)
	logging.Init(false)
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for empty logger name")
		}
	}()
	logging.NewLogger("")
}

// TestLogFileOutput verifies that TFM_LOG_PATH causes log lines to be written
// to the specified file.
func TestLogFileOutput(t *testing.T) {
	resetEnv(t)

	dir := t.TempDir()
	logFile := filepath.Join(dir, "tfm-test.log")

	os.Setenv(logging.EnvLog, "DEBUG")
	os.Setenv(logging.EnvLogPath, logFile)

	logging.Init(false)

	l := logging.NewLogger("filetest")
	l.Debug("hello from test", "key", "value")

	// Give the writer time to flush (it's synchronous in hclog, but be safe).
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("could not read log file %s: %v", logFile, err)
	}
	if !strings.Contains(string(content), "hello from test") {
		t.Errorf("log file missing expected message; got:\n%s", string(content))
	}
}

func TestValidLevels(t *testing.T) {
	expected := []string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "OFF"}
	if len(logging.ValidLevels) != len(expected) {
		t.Fatalf("expected %d valid levels, got %d", len(expected), len(logging.ValidLevels))
	}
	for i, l := range expected {
		if logging.ValidLevels[i] != l {
			t.Errorf("ValidLevels[%d] = %q, want %q", i, logging.ValidLevels[i], l)
		}
	}
}
