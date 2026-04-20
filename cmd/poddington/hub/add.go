package hub

import (
	"errors"

	"github.com/jefjesuswt/poddington/internal/fleet"
	"github.com/jefjesuswt/poddington/shared/ui"
	"github.com/spf13/cobra"
)

var addCommand = &cobra.Command{
	Use:   "add [name] [address]",
	Short: "Registers a new node to the hub",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		svc, err := InitFleetService(ctx)
		if err != nil {
			return err
		}

		name := args[0]
		address := args[1]

		node, err := svc.RegisterNode(ctx, name, address)
		if err != nil {
			if errors.Is(err, fleet.ErrNodeAlreadyExists) {
				return ui.WrapError("node already exists: %s", name)
			}
			return ui.WrapError("failed to register node: %w", err)
		}

		ui.PrintSuccess("Node '%s' successfully registered to the Hub.", node.Name)
		ui.PrintWarning("Save this token. It will not be shown again!")

		ui.PrintKeyValue("Node ID", node.ID)
		ui.PrintKeyValue("Address", node.Address)
		ui.PrintKeyValue("Token", node.Token)

		return nil
	},
}
