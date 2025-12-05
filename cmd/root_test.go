package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootCommand(t *testing.T) {
	if rootCmd.Use != "protect" {
		t.Errorf("Expected Use to be 'protect', got '%s'", rootCmd.Use)
	}

	if rootCmd.Short == "" {
		t.Error("Expected Short description to be set")
	}

	if rootCmd.Long == "" {
		t.Error("Expected Long description to be set")
	}
}

func TestRootCommandFlags(t *testing.T) {
	// Test persistent flags
	persistentFlags := rootCmd.PersistentFlags()

	urlFlag := persistentFlags.Lookup("url")
	if urlFlag == nil {
		t.Error("Expected 'url' flag to be registered")
	}

	tokenFlag := persistentFlags.Lookup("token")
	if tokenFlag == nil {
		t.Error("Expected 'token' flag to be registered")
	}

	logLevelFlag := persistentFlags.Lookup("log-level")
	if logLevelFlag == nil {
		t.Error("Expected 'log-level' flag to be registered")
	}

	if logLevelFlag != nil && logLevelFlag.DefValue != "none" {
		t.Errorf("Expected log-level default to be 'none', got '%s'", logLevelFlag.DefValue)
	}

	// Test command-specific flags
	flags := rootCmd.Flags()

	portFlag := flags.Lookup("port")
	if portFlag == nil {
		t.Error("Expected 'port' flag to be registered")
	}

	viewFlag := flags.Lookup("view")
	if viewFlag == nil {
		t.Error("Expected 'view' flag to be registered")
	}

	cameraFlag := flags.Lookup("camera")
	if cameraFlag == nil {
		t.Error("Expected 'camera' flag to be registered")
	}

	presetFlag := flags.Lookup("preset")
	if presetFlag == nil {
		t.Error("Expected 'preset' flag to be registered")
	}

	listFlag := flags.Lookup("list")
	if listFlag == nil {
		t.Error("Expected 'list' flag to be registered")
	}

	showIDsFlag := flags.Lookup("show-ids")
	if showIDsFlag == nil {
		t.Error("Expected 'show-ids' flag to be registered")
	}

	tuiFlag := flags.Lookup("tui")
	if tuiFlag == nil {
		t.Error("Expected 'tui' flag to be registered")
	}
}

func TestRootCommandHelp(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("Expected help to execute without error, got: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "protect") {
		t.Error("Expected help output to contain 'protect'")
	}

	// Check that flag-based options are documented
	if !strings.Contains(output, "--port") {
		t.Error("Expected help output to contain '--port' flag")
	}

	if !strings.Contains(output, "--view") {
		t.Error("Expected help output to contain '--view' flag")
	}

	if !strings.Contains(output, "--camera") {
		t.Error("Expected help output to contain '--camera' flag")
	}

	if !strings.Contains(output, "--list") {
		t.Error("Expected help output to contain '--list' flag")
	}

	rootCmd.SetArgs([]string{})
}

func TestRootCommandNoSubcommands(t *testing.T) {
	// Verify that legacy subcommands are not registered
	commands := rootCmd.Commands()

	for _, cmd := range commands {
		// Only help and completion commands should exist
		if cmd.Name() != "help" && cmd.Name() != "completion" {
			t.Errorf("Unexpected command '%s' found - only flag-based usage should be supported", cmd.Name())
		}
	}
}
