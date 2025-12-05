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
	flags := rootCmd.PersistentFlags()

	urlFlag := flags.Lookup("url")
	if urlFlag == nil {
		t.Error("Expected 'url' flag to be registered")
	}

	tokenFlag := flags.Lookup("token")
	if tokenFlag == nil {
		t.Error("Expected 'token' flag to be registered")
	}

	logLevelFlag := flags.Lookup("log-level")
	if logLevelFlag == nil {
		t.Error("Expected 'log-level' flag to be registered")
	}

	if logLevelFlag != nil && logLevelFlag.DefValue != "none" {
		t.Errorf("Expected log-level default to be 'none', got '%s'", logLevelFlag.DefValue)
	}
}

func TestRootCommandHasSubcommands(t *testing.T) {
	expectedCommands := []string{"viewport", "liveview", "camera"}

	commands := rootCmd.Commands()
	commandMap := make(map[string]bool)
	for _, cmd := range commands {
		commandMap[cmd.Name()] = true
	}

	for _, expected := range expectedCommands {
		if !commandMap[expected] {
			t.Errorf("Expected subcommand '%s' to be registered", expected)
		}
	}

	// Test that we have at least our 3 main commands
	if len(commands) < 3 {
		t.Errorf("Expected at least 3 commands, got %d", len(commands))
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

	if !strings.Contains(output, "Available Commands") {
		t.Error("Expected help output to contain 'Available Commands'")
	}

	rootCmd.SetArgs([]string{})
}

func TestViewportCommand(t *testing.T) {
	viewportCmd := rootCmd.Commands()
	found := false
	for _, cmd := range viewportCmd {
		if cmd.Name() == "viewport" {
			found = true
			if cmd.Short == "" {
				t.Error("Expected viewport command to have Short description")
			}
			break
		}
	}

	if !found {
		t.Error("Expected viewport command to be registered")
	}
}

func TestLiveviewCommand(t *testing.T) {
	liveviewCmd := rootCmd.Commands()
	found := false
	for _, cmd := range liveviewCmd {
		if cmd.Name() == "liveview" {
			found = true
			if cmd.Short == "" {
				t.Error("Expected liveview command to have Short description")
			}
			break
		}
	}

	if !found {
		t.Error("Expected liveview command to be registered")
	}
}

func TestCameraCommand(t *testing.T) {
	cameraCmd := rootCmd.Commands()
	found := false
	for _, cmd := range cameraCmd {
		if cmd.Name() == "camera" {
			found = true
			if cmd.Short == "" {
				t.Error("Expected camera command to have Short description")
			}
			break
		}
	}

	if !found {
		t.Error("Expected camera command to be registered")
	}
}
