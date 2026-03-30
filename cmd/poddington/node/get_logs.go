package node

import (
	"context"

	"github.com/jefjesuswt/poddington/internal/container"
	"github.com/jefjesuswt/poddington/shared/ui"
	"github.com/spf13/cobra"
)

var logsCommand = &cobra.Command{
	Use: "logs [container_name]",
	Short: "Fetch the logs of a container",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		return WithContainerAction(ctx, args, func(ctx context.Context, svc *container.Service, target, higlighted string) error {
			logsData, err := svc.GetLogs(ctx, target)
			if err != nil {
				return ui.PrintError("error fetching logs: %w", err)
			}
			ui.PrintLogs(target, logsData)
			return nil
		})
	},
}
