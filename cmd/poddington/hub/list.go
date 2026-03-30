package hub

import (
	"fmt"
	"time"

	"github.com/jefjesuswt/poddington/shared/ui"
	"github.com/spf13/cobra"
)

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "List registered nodes in the hub",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		svc, err := InitFleetService(ctx)
		if err != nil {
			return err
		}

		nodes, err := svc.ListNodes(ctx)
		if err != nil {
			return ui.PrintError("failed to list nodes: %w", err)
		}

		ui.PrintTitle("Fleet nodes: %d", len(nodes))

		if len(nodes) > 0 {
			for _, node := range nodes {
				lastSeen := time.Since(node.LastSeen).Round(time.Second).String()

				ui.PrintNodeRow(node.ID, node.Name, node.Address, lastSeen)
			}
		}

		fmt.Println()
		return nil
	},
}
