package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/methridge/protect/internal/logger"
	"github.com/spf13/cobra"
)

var liveviewCmd = &cobra.Command{
	Use:   "liveview",
	Short: "Manage liveviews",
	Long:  `List and switch between UniFi Protect liveviews.`,
}

var liveviewListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available liveviews",
	RunE: func(cmd *cobra.Command, args []string) error {
		log := logger.Get()
		showIDs, _ := cmd.Flags().GetBool("show-ids")

		client, err := getClient()
		if err != nil {
			return err
		}

		cameras, err := client.ListCameras()
		if err != nil {
			return fmt.Errorf("failed to list cameras: %w", err)
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
	},
}

var liveviewSwitchCmd = &cobra.Command{
	Use:   "switch [viewport-id-or-name] [liveview-id-or-name]",
	Short: "Switch a viewport to a specific liveview",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		log := logger.Get()
		viewportIdentifier := args[0]
		liveviewIdentifier := args[1]

		client, err := getClient()
		if err != nil {
			return err
		}

		// Try to find viewport by name or ID
		viewports, err := client.ListViewports()
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

		// Try to find liveview by name or ID
		liveviews, err := client.ListCameras()
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

		if err := client.SwitchCamera(viewportID, liveviewID); err != nil {
			return fmt.Errorf("failed to switch liveview: %w", err)
		}

		fmt.Printf("Successfully switched viewport %s to liveview %s\n", viewportIdentifier, liveviewIdentifier)
		log.Infow("Switched liveview", "viewportID", viewportID, "liveviewID", liveviewID)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(liveviewCmd)
	liveviewCmd.AddCommand(liveviewListCmd)
	liveviewCmd.AddCommand(liveviewSwitchCmd)

	liveviewListCmd.Flags().Bool("show-ids", false, "Show IDs in addition to names")
}
