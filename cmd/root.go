package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

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
		// Check for flag-based operations
		listMode, _ := cmd.Flags().GetString("list")
		viewport, _ := cmd.Flags().GetString("port")
		liveview, _ := cmd.Flags().GetString("view")
		camera, _ := cmd.Flags().GetString("camera")
		preset, _ := cmd.Flags().GetInt("preset")
		showIDs, _ := cmd.Flags().GetBool("show-ids")
		launchTUI, _ := cmd.Flags().GetBool("tui")

		c, err := getClient()
		if err != nil {
			return err
		}

		// Handle TUI launch
		if launchTUI {
			return tui.Run(c)
		}

		// Handle list operations
		if listMode != "" {
			return handleListOperation(c, listMode, showIDs)
		}

		// Handle viewport switching
		if viewport != "" && liveview != "" {
			return handleViewportSwitch(c, viewport, liveview)
		}

		// Handle camera PTZ operations
		if camera != "" {
			return handleCameraOperation(c, camera, preset)
		}

		// If no flags specified, show help
		return cmd.Help()
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

	// Flag-based options
	rootCmd.Flags().BoolP("tui", "i", false, "Launch interactive TUI")
	rootCmd.Flags().StringP("port", "p", "", "Viewport name or ID (use with --view)")
	rootCmd.Flags().StringP("view", "v", "", "Liveview/camera name or ID (use with --port)")
	rootCmd.Flags().StringP("camera", "c", "", "Camera name or ID for PTZ operations (use with --preset)")
	rootCmd.Flags().IntP("preset", "P", -2, "PTZ preset position (-1 for home, 0-9 for presets)")
	rootCmd.Flags().StringP("list", "L", "", "List items: 'viewports', 'liveviews', or 'cameras'")
	rootCmd.Flags().Bool("show-ids", false, "Show IDs when listing")
}

func getClient() (*client.Client, error) {
	cfg := config.Get()
	return client.NewClient(cfg.ProtectURL, cfg.APIToken), nil
}

func handleListOperation(c *client.Client, listType string, showIDs bool) error {
	switch listType {
	case "viewports":
		return listViewports(c, showIDs)
	case "liveviews", "views":
		return listLiveviews(c, showIDs)
	case "cameras":
		return listCameras(c, showIDs)
	default:
		return fmt.Errorf("invalid list type: %s (use 'viewports', 'liveviews', or 'cameras')", listType)
	}
}

func handleViewportSwitch(c *client.Client, viewportIdentifier, liveviewIdentifier string) error {
	log := logger.Get()

	// Find viewport by name or ID
	viewports, err := c.ListViewports()
	if err != nil {
		return fmt.Errorf("failed to list viewports: %w", err)
	}

	var viewportID string
	for _, vp := range viewports {
		if vp.ID == viewportIdentifier || vp.Name == viewportIdentifier {
			viewportID = vp.ID
			break
		}
	}

	if viewportID == "" {
		return fmt.Errorf("viewport not found: %s", viewportIdentifier)
	}

	// Find liveview by name or ID
	liveviews, err := c.ListCameras()
	if err != nil {
		return fmt.Errorf("failed to list liveviews: %w", err)
	}

	var liveviewID string
	for _, lv := range liveviews {
		if lv.ID == liveviewIdentifier || lv.Name == liveviewIdentifier {
			liveviewID = lv.ID
			break
		}
	}

	if liveviewID == "" {
		return fmt.Errorf("liveview not found: %s", liveviewIdentifier)
	}

	if err := c.SwitchViewport(viewportID, liveviewID); err != nil {
		return fmt.Errorf("failed to switch viewport: %w", err)
	}

	fmt.Printf("Successfully switched viewport %s to liveview %s\n", viewportIdentifier, liveviewIdentifier)
	log.Infow("Switched viewport", "viewportID", viewportID, "liveviewID", liveviewID)

	return nil
}

func handleCameraOperation(c *client.Client, cameraNameOrID string, preset int) error {
	if preset == -2 {
		return fmt.Errorf("--preset flag is required when using --camera")
	}

	if preset < -1 || preset > 9 {
		return fmt.Errorf("invalid preset value: %d (must be between -1 and 9)", preset)
	}

	cameras, err := c.ListPTZCameras()
	if err != nil {
		return err
	}

	var cameraID string
	var cameraName string
	for _, cam := range cameras {
		if cam.ID == cameraNameOrID || cam.Name == cameraNameOrID {
			cameraID = cam.ID
			cameraName = cam.Name
			break
		}
	}

	if cameraID == "" {
		return fmt.Errorf("camera not found: %s", cameraNameOrID)
	}

	if err := c.MovePTZToPreset(cameraID, preset); err != nil {
		return err
	}

	presetLabel := fmt.Sprintf("preset %d", preset)
	if preset == -1 {
		presetLabel = "home position"
	}

	fmt.Printf("Successfully moved camera '%s' to %s\n", cameraName, presetLabel)
	return nil
}

func listViewports(c *client.Client, showIDs bool) error {
	log := logger.Get()

	viewports, err := c.ListViewports()
	if err != nil {
		return fmt.Errorf("failed to list viewports: %w", err)
	}

	if len(viewports) == 0 {
		fmt.Println("No viewports found")
		return nil
	}

	liveviews, err := c.ListCameras()
	if err != nil {
		return fmt.Errorf("failed to list liveviews: %w", err)
	}

	liveviewMap := make(map[string]string)
	for _, lv := range liveviews {
		liveviewMap[lv.ID] = lv.Name
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if showIDs {
		fmt.Fprintln(w, "NAME\tCURRENT LIVEVIEW\tID\tLIVEVIEW ID")
		fmt.Fprintln(w, "----\t----------------\t--\t-----------")
	} else {
		fmt.Fprintln(w, "NAME\tCURRENT LIVEVIEW")
		fmt.Fprintln(w, "----\t----------------")
	}

	for _, vp := range viewports {
		liveviewName := liveviewMap[vp.Liveview]
		if liveviewName == "" {
			liveviewName = vp.Liveview
		}
		if showIDs {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", vp.Name, liveviewName, vp.ID, vp.Liveview)
		} else {
			fmt.Fprintf(w, "%s\t%s\n", vp.Name, liveviewName)
		}
	}

	w.Flush()
	log.Infow("Listed viewports", "count", len(viewports))
	return nil
}

func listLiveviews(c *client.Client, showIDs bool) error {
	log := logger.Get()

	cameras, err := c.ListCameras()
	if err != nil {
		return fmt.Errorf("failed to list liveviews: %w", err)
	}

	if len(cameras) == 0 {
		fmt.Println("No liveviews found")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if showIDs {
		fmt.Fprintln(w, "NAME\tID")
		fmt.Fprintln(w, "----\t--")
		for _, cam := range cameras {
			fmt.Fprintf(w, "%s\t%s\n", cam.Name, cam.ID)
		}
	} else {
		fmt.Fprintln(w, "NAME")
		fmt.Fprintln(w, "----")
		for _, cam := range cameras {
			fmt.Fprintf(w, "%s\n", cam.Name)
		}
	}

	w.Flush()
	log.Infow("Listed liveviews", "count", len(cameras))
	return nil
}

func listCameras(c *client.Client, showIDs bool) error {
	cameras, err := c.ListPTZCameras()
	if err != nil {
		return err
	}

	if len(cameras) == 0 {
		fmt.Println("No PTZ cameras found")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if showIDs {
		fmt.Fprintln(w, "ID\tNAME")
		fmt.Fprintln(w, "--\t----")
		for _, camera := range cameras {
			fmt.Fprintf(w, "%s\t%s\n", camera.ID, camera.Name)
		}
	} else {
		fmt.Fprintln(w, "NAME")
		fmt.Fprintln(w, "----")
		for _, camera := range cameras {
			fmt.Fprintf(w, "%s\n", camera.Name)
		}
	}
	w.Flush()

	return nil
}
