package node

import (
	"context"

	"github.com/jefjesuswt/poddington/internal/core/container"
	"github.com/jefjesuswt/poddington/internal/ui"
	"github.com/spf13/cobra"
)
var restartCommand = &cobra.Command{
	Use: "restart [container_name]",
	Short: "Restarts a container",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return WithContainerAction(args, func(ctx context.Context, svc *container.Service, target, highlighted string) error {
			if err := svc.Restart(ctx, target); err != nil {
				return ui.PrintError("error restarting container: %w", err)
			}
			ui.PrintSuccess("Container %s has been restarted.", highlighted)
			return nil
		})
	},
}
