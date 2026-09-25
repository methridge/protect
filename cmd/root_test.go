package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/methridge/protect/internal/client"
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

// newTestServer serves fixed viewports, liveviews and cameras, and records
// PATCH/POST request paths and bodies in requests.
func newTestServer(t *testing.T, requests *[]string) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET" && r.URL.Path == "/proxy/protect/integration/v1/viewers":
			json.NewEncoder(w).Encode([]client.Viewport{
				{ID: "vp1", Name: "Office"},
				{ID: "vp2", Name: "Office:TV"},
				{ID: "vp3", Name: "A"},
				{ID: "vp4", Name: "A:B"},
			})
		case r.Method == "GET" && r.URL.Path == "/proxy/protect/integration/v1/liveviews":
			json.NewEncoder(w).Encode([]client.Liveview{
				{ID: "lv1", Name: "Driveway"},
				{ID: "lv2", Name: "Cam 1:Cam 2"},
				{ID: "lv3", Name: "B:C"},
				{ID: "lv4", Name: "C"},
			})
		case r.Method == "GET" && r.URL.Path == "/proxy/protect/integration/v1/cameras":
			json.NewEncoder(w).Encode([]client.PTZCamera{
				{ID: "cam1", Name: "Front Door"},
				{ID: "cam2", Name: "Gate:North"},
			})
		default:
			body, _ := io.ReadAll(r.Body)
			*requests = append(*requests, r.Method+" "+r.URL.Path+" "+string(body))
		}
	}))
	t.Cleanup(server.Close)
	return server
}

func TestHandleSwitchCommand(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    string
		wantErr string
	}{
		{name: "simple names", arg: "Office:Driveway", want: `PATCH /proxy/protect/integration/v1/viewers/vp1 {"liveview":"lv1"}`},
		{name: "colon in viewport name", arg: "Office:TV:Driveway", want: `PATCH /proxy/protect/integration/v1/viewers/vp2 {"liveview":"lv1"}`},
		{name: "colon in liveview name", arg: "Office:Cam 1:Cam 2", want: `PATCH /proxy/protect/integration/v1/viewers/vp1 {"liveview":"lv2"}`},
		{name: "ambiguous split", arg: "A:B:C", wantErr: "ambiguous"},
		{name: "no match with colons", arg: "X:Y:Z", wantErr: "no viewport and liveview match"},
		{name: "unknown viewport", arg: "Nope:Driveway", wantErr: "viewport not found: Nope"},
		{name: "missing separator", arg: "Office", wantErr: "invalid switch format"},
		{name: "empty side", arg: "Office:", wantErr: "cannot be empty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requests []string
			server := newTestServer(t, &requests)
			c := client.NewClient(server.URL, "test-token")

			err := handleSwitchCommand(c, tt.arg)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("Expected error containing %q, got %v", tt.wantErr, err)
				}
				if len(requests) != 0 {
					t.Errorf("Expected no switch request, got %v", requests)
				}
				return
			}
			if err != nil {
				t.Fatalf("handleSwitchCommand(%q) error = %v", tt.arg, err)
			}
			if len(requests) != 1 || strings.TrimSpace(requests[0]) != tt.want {
				t.Errorf("Expected request %q, got %v", tt.want, requests)
			}
		})
	}
}

func TestHandlePTZCommand(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    string
		wantErr string
	}{
		{name: "home position", arg: "Front Door:-1", want: "POST /proxy/protect/integration/v1/cameras/cam1/ptz/goto/-1"},
		{name: "colon in camera name", arg: "Gate:North:5", want: "POST /proxy/protect/integration/v1/cameras/cam2/ptz/goto/5"},
		{name: "missing separator", arg: "Front Door", wantErr: "invalid ptz format"},
		{name: "non-numeric preset", arg: "Front Door:x", wantErr: "invalid preset value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requests []string
			server := newTestServer(t, &requests)
			c := client.NewClient(server.URL, "test-token")

			err := handlePTZCommand(c, tt.arg)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("Expected error containing %q, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("handlePTZCommand(%q) error = %v", tt.arg, err)
			}
			if len(requests) != 1 || strings.TrimSpace(requests[0]) != tt.want {
				t.Errorf("Expected request %q, got %v", tt.want, requests)
			}
		})
	}
}

func TestVersionSkipsConfigValidation(t *testing.T) {
	flag := rootCmd.Flags().Lookup("version")
	if err := flag.Value.Set("true"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		flag.Value.Set("false")
		flag.Changed = false
	})

	if err := rootCmd.PersistentPreRunE(rootCmd, nil); err != nil {
		t.Errorf("Expected --version to skip configuration, got: %v", err)
	}
}
