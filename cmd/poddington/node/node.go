package node

import (
	"github.com/charmbracelet/log"
	"github.com/jefjesuswt/poddington/internal/ui"
	"github.com/spf13/cobra"
)

var HubAddress string

var Cmd = &cobra.Command{
	Use:   "node",
	Short: "Run a node that connects to the hub",
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.PrintTitle("Poddington Node Initialization")
		log.Info("Poddington starting...", "mode", "NODE")
		log.Info("Connecting to hub...", "hub", HubAddress)
		return nil
	},
}

func init() {
	Cmd.Flags().StringVar(&HubAddress, "hub-address", "127.0.0.1:8443", "IP address and port to connect to.")

	Cmd.AddCommand(listCommand)
	listCommand.Flags().BoolVarP(&listAll, "all", "a", false, "Show all containers (default shows just running)")

	Cmd.AddCommand(rmCommand)
	rmCommand.Flags().BoolVarP(&forceRemove, "force", "f", false, "Force the removal of an already container")

	Cmd.AddCommand(inspectCommand)
	Cmd.AddCommand(startCommand)
	Cmd.AddCommand(stopCommand)
	Cmd.AddCommand(restartCommand)
	Cmd.AddCommand(logsCommand)
}
