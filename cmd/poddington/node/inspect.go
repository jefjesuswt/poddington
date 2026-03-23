package node

import (
	"context"
	"fmt"

	"github.com/jefjesuswt/poddington/internal/core/container"
	"github.com/jefjesuswt/poddington/internal/ui"
	"github.com/spf13/cobra"
)

var inspectCommand = &cobra.Command{
	Use:   "inspect [container_name]",
	Short: "Inspect a container, gets deep technical details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return WithContainerAction(args, func(ctx context.Context, svc *container.Service, target, highlighted string) error {
			instance, err := svc.Inspect(ctx, target)
			if err != nil {
				return ui.PrintError("error inspecting container: %w", err)
			}

			ui.PrintTitle("Container Details: %s", instance.Name)
			ui.PrintKeyValue("ID", instance.ID)
			ui.PrintKeyValue("State", instance.State)
			ui.PrintKeyValue("Created", instance.Created)
			ui.PrintKeyValue("Image", instance.Image)
			ui.PrintKeyValue("IP Addr", instance.IPAddress)
			ui.PrintKeyValue("Command", instance.Cmd)

			ui.PrintList("Ports", instance.Ports)
			ui.PrintList("Mounts", instance.Mounts)

			fmt.Println()

			return nil
		})
	},
}
