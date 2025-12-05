package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	// Test that Load initializes without error when config file doesn't exist
	cfg = nil

	config, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if config == nil {
		t.Fatal("Expected config to be initialized")
	}

	// Default log level should be "none"
	if config.LogLevel != "none" {
		t.Errorf("Expected default LogLevel to be 'none', got '%s'", config.LogLevel)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				ProtectURL: "https://protect.example.com",
				APIToken:   "test-token",
			},
			wantErr: false,
		},
		{
			name: "missing URL",
			config: Config{
				APIToken: "test-token",
			},
			wantErr: true,
		},
		{
			name: "missing token",
			config: Config{
				ProtectURL: "https://protect.example.com",
			},
			wantErr: true,
		},
		{
			name:    "empty config",
			config:  Config{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
