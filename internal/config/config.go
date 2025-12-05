package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	ProtectURL string `mapstructure:"protect_url"`
	APIToken   string `mapstructure:"api_token"`
	LogLevel   string `mapstructure:"log_level"`
}

var cfg *Config

// Load loads the configuration from the XDG config directory
func Load() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Try XDG config directory first (Linux/BSD standard)
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		viper.AddConfigPath(filepath.Join(xdgConfig, "protect"))
	}

	// Try ~/.config (common convention)
	if home, err := os.UserHomeDir(); err == nil {
		viper.AddConfigPath(filepath.Join(home, ".config", "protect"))
	}

	// Try OS-specific config directory (macOS: ~/Library/Application Support)
	if configDir, err := os.UserConfigDir(); err == nil {
		appConfigDir := filepath.Join(configDir, "protect")
		viper.AddConfigPath(appConfigDir)
		// Create config directory if it doesn't exist
		os.MkdirAll(appConfigDir, 0755)
	}

	// Set defaults
	viper.SetDefault("log_level", "none")

	// Allow environment variables to override config
	// This must be set before reading the config file
	viper.SetEnvPrefix("PROTECT")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}

// Get returns the current configuration
func Get() *Config {
	if cfg == nil {
		c, err := Load()
		if err != nil {
			return &Config{}
		}
		return c
	}
	return cfg
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.ProtectURL == "" {
		return fmt.Errorf("protect_url is required")
	}
	if c.APIToken == "" {
		return fmt.Errorf("api_token is required")
	}
	return nil
}
