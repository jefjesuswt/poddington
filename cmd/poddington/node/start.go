package node

import (
	"context"
	"errors"

	"github.com/jefjesuswt/poddington/internal/core/container"
	"github.com/jefjesuswt/poddington/internal/ui"
	"github.com/spf13/cobra"
)

var startCommand = &cobra.Command{
	Use:   "start [container_name]",
	Short: "Start a container",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return WithContainerAction(args, func(ctx context.Context, svc *container.Service, target, highlighted string) error {
			if err := svc.Start(ctx, target); err != nil {
				if errors.Is(err, container.ErrAlreadyRunning) {
					ui.PrintWarning("Container %s is already running.", highlighted)
					return nil
				}
				return ui.PrintError("error starting container: %w", err)
			}
			ui.PrintSuccess("Container %s is now running.", highlighted)
			return nil
		})
	},
}
