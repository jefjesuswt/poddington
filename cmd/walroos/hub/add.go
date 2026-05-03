package hub

import (
	"errors"
	"log/slog"

	"github.com/jefjesuswt/walroos/internal/fleet"
	"github.com/jefjesuswt/walroos/shared/ui"
	"github.com/spf13/cobra"
)

var addCommand = &cobra.Command{
	Use:   "add [name] [address]",
	Short: "Registers a new node to the hub",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		svc, cleanup, err := InitFleetService(ctx)
		if err != nil {
			return err
		}

		defer func() {
			if err := cleanup(); err != nil {
				slog.Error("error closing database", "err", err)
			} else {
				slog.Info("Database closed successfully")
			}
		}()

		name := args[0]
		address := args[1]

		node, err := svc.RegisterNode(ctx, name, address)
		if err != nil {
			if errors.Is(err, fleet.ErrNodeNameAlreadyExists) {
				return ui.WrapError("node name already exists: %s", name)
			}

			if errors.Is(err, fleet.ErrNodeAddressAlreadyExists) {
				return ui.WrapError("node address already exists: %s", address)
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
