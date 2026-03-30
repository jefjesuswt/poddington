package node

import (
	"context"
	"fmt"

	"github.com/jefjesuswt/poddington/internal/container"
	"github.com/jefjesuswt/poddington/shared/ui"
	"github.com/spf13/cobra"
)

var listAll bool

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "List node's local containers",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		return WithContainerAction(ctx, args, func(ctx context.Context, svc *container.Service, target, highlighted string) error {
			instances, err := svc.List(ctx, listAll)
			if err != nil {
				return ui.PrintError("error listing containers: %w", err)
			}

			title := "Active containers"
			if listAll {
				title = "All containers (including stopped)"
			}

			ui.PrintTitle("%s: %d", title, len(instances))
			if len(instances) > 0 {
				ui.PrintListHeader()
				for _, instance := range instances {
					ui.PrintContainerRow(instance.ID, instance.Name, instance.State, instance.Image)
				}
			}

			fmt.Println()
			return nil
		})
	},
}
