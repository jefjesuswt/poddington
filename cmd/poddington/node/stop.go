package node

import (
	"context"
	"errors"

	"github.com/jefjesuswt/poddington/internal/container"
	"github.com/jefjesuswt/poddington/shared/ui"
	"github.com/spf13/cobra"
)

var stopCommand = &cobra.Command{
	Use:   "stop [container_name]",
	Short: "Stops a container",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		return WithContainerAction(ctx, args, func(ctx context.Context, svc *container.Service, target, highlighted string) error {
			if err := svc.Stop(ctx, target); err != nil {
				if errors.Is(err, container.ErrAlreadyStopped) {
					ui.PrintWarning("Container %s is already stopped.", highlighted)
					return nil
				}
				return ui.PrintError("error stopping container: %w", err)
			}

			ui.PrintSuccess("Container %s has been stopped.", highlighted)
			return nil
		})

	},
}
