package hub

import (
	"errors"

	"charm.land/lipgloss/v2"
	"github.com/jefjesuswt/poddington/internal/fleet"
	"github.com/jefjesuswt/poddington/shared/ui"
	"github.com/spf13/cobra"
)

var removeCommand = &cobra.Command{
	Use: "rm [node_id]",
	Short: "Remove a node from the hub",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		svc, err := InitFleetService(ctx)
		if err != nil {
			return err
		}

		id := args[0]

		if err := svc.RemoveNode(ctx, id); err != nil {
			if errors.Is(err, fleet.ErrNodeNotFound) {
				return ui.PrintError("node %s not found", id)
			}
			return ui.PrintError("failed to remove node %s", id)
		}

		highlighted := lipgloss.NewStyle().Foreground(ui.PodGrape).Bold(true).Render(id)
		ui.PrintSuccess("Node %s has been permanently removed from the fleet.", highlighted)

		return nil
	},
}
