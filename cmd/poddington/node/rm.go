package node

import (
	"context"

	"github.com/jefjesuswt/poddington/internal/container"
	"github.com/jefjesuswt/poddington/shared/ui"
	"github.com/spf13/cobra"
)

var forceRemove bool

var rmCommand = &cobra.Command{
	Use: "rm [container_name]",
	Short: "Removes a container",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		return WithContainerAction(ctx, args, func(ctx context.Context, svc *container.Service, target, highlighted string) error {
			if err := svc.Remove(ctx, target, forceRemove); err != nil {
				return ui.PrintError("error removing container: %w", err)
			}
			ui.PrintSuccess("Container %s has been removed.", highlighted)
			return nil
		})
	},
}
