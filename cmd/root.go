package cmd

import (
	"fmt"

	"github.com/methridge/protect/internal/client"
	"github.com/methridge/protect/internal/config"
	"github.com/methridge/protect/internal/logger"
	"github.com/methridge/protect/internal/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "protect",
	Short: "UniFi Protect View Switcher",
	Long:  `A TUI tool for switching between different camera views and viewports in UniFi Protect.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no subcommand is specified, launch the TUI
		c, err := getClient()
		if err != nil {
			return err
		}
		return tui.Run(c)
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		// Override config with flags if provided
		protectURL, _ := cmd.Flags().GetString("url")
		if protectURL != "" {
			cfg.ProtectURL = protectURL
		}

		apiToken, _ := cmd.Flags().GetString("token")
		if apiToken != "" {
			cfg.APIToken = apiToken
		}

		logLevel, _ := cmd.Flags().GetString("log-level")
		if logLevel != "" {
			cfg.LogLevel = logLevel
		}

		// Set log level
		if err := logger.SetLevel(cfg.LogLevel); err != nil {
			return fmt.Errorf("failed to set log level: %w", err)
		}

		// Validate configuration
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("invalid configuration: %w", err)
		}

		return nil
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP("url", "u", "", "UniFi Protect URL")
	rootCmd.PersistentFlags().StringP("token", "t", "", "API token for authentication")
	rootCmd.PersistentFlags().StringP("log-level", "l", "none", "Log level (none, debug, info, warn, error)")
}

func getClient() (*client.Client, error) {
	cfg := config.Get()
	return client.NewClient(cfg.ProtectURL, cfg.APIToken), nil
}
