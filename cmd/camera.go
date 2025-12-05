package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/methridge/protect/internal/client"
	"github.com/methridge/protect/internal/config"
	"github.com/spf13/cobra"
)

var showCameraIDs bool

var cameraCmd = &cobra.Command{
	Use:   "camera",
	Short: "Manage PTZ cameras",
	Long:  `List PTZ cameras and move them to preset positions.`,
}

var listCamerasCmd = &cobra.Command{
	Use:   "list",
	Short: "List all PTZ cameras",
	Long:  `List all available PTZ-capable cameras in UniFi Protect.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		c := client.NewClient(cfg.ProtectURL, cfg.APIToken)
		cameras, err := c.ListPTZCameras()
		if err != nil {
			return err
		}

		if len(cameras) == 0 {
			fmt.Println("No PTZ cameras found")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		if showCameraIDs {
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
	},
}

var movePTZCmd = &cobra.Command{
	Use:   "goto <camera-name-or-id> <preset>",
	Short: "Move PTZ camera to preset position",
	Long: `Move a PTZ camera to a specific preset position.

Preset values:
  -1: Home position (use -- before -1, e.g., "protect camera goto Tower -- -1")
  0-9: Preset slots 0 through 9

Examples:
  protect camera goto "Front Door" -- -1
  protect camera goto 68faa8300029d203e4115927 0
  protect camera goto Driveway 3`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cameraNameOrID := args[0]
		presetStr := args[1]

		preset, err := strconv.Atoi(presetStr)
		if err != nil {
			return fmt.Errorf("invalid preset value: %s (must be a number between -1 and 9)", presetStr)
		}

		if preset < -1 || preset > 9 {
			return fmt.Errorf("invalid preset value: %d (must be between -1 and 9)", preset)
		}

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		c := client.NewClient(cfg.ProtectURL, cfg.APIToken)

		// Get all PTZ cameras to find the camera by name or ID
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

		// Move the camera to the preset
		if err := c.MovePTZToPreset(cameraID, preset); err != nil {
			return err
		}

		presetLabel := fmt.Sprintf("preset %d", preset)
		if preset == -1 {
			presetLabel = "home position"
		}

		fmt.Printf("Successfully moved camera '%s' to %s\n", cameraName, presetLabel)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cameraCmd)
	cameraCmd.AddCommand(listCamerasCmd)
	cameraCmd.AddCommand(movePTZCmd)

	listCamerasCmd.Flags().BoolVar(&showCameraIDs, "show-ids", false, "Show camera IDs alongside names")
}
