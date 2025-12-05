package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/methridge/protect/internal/logger"
	"github.com/spf13/cobra"
)

var viewportCmd = &cobra.Command{
	Use:   "viewport",
	Short: "Manage viewports",
	Long:  `List and switch between UniFi Protect viewports.`,
}

var viewportListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available viewports",
	RunE: func(cmd *cobra.Command, args []string) error {
		log := logger.Get()
		showIDs, _ := cmd.Flags().GetBool("show-ids")

		client, err := getClient()
		if err != nil {
			return err
		}

		viewports, err := client.ListViewports()
		if err != nil {
			return fmt.Errorf("failed to list viewports: %w", err)
		}

		if len(viewports) == 0 {
			fmt.Println("No viewports found")
			return nil
		}

		// Get all liveviews to map IDs to names
		liveviews, err := client.ListCameras()
		if err != nil {
			return fmt.Errorf("failed to list liveviews: %w", err)
		}

		// Create a map of liveview ID to name
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
				liveviewName = vp.Liveview // Fall back to ID if name not found
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
	},
}

var viewportSwitchCmd = &cobra.Command{
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

		if err := client.SwitchViewport(viewportID, liveviewID); err != nil {
			return fmt.Errorf("failed to switch viewport: %w", err)
		}

		fmt.Printf("Successfully switched viewport %s to liveview %s\n", viewportIdentifier, liveviewIdentifier)
		log.Infow("Switched viewport", "viewportID", viewportID, "liveviewID", liveviewID)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(viewportCmd)
	viewportCmd.AddCommand(viewportListCmd)
	viewportCmd.AddCommand(viewportSwitchCmd)

	viewportListCmd.Flags().Bool("show-ids", false, "Show IDs in addition to names")
}
