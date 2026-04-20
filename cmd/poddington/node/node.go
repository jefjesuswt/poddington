package node

import (
	"log/slog"

	"github.com/spf13/cobra"
)

var HubAddress string

var Cmd = &cobra.Command{
	Use:   "node",
	Short: "Run a node that connects to the hub",
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.Info("Poddington starting...", "mode", "NODE")
		slog.Info("Connecting to hub...", "hub", HubAddress)
		return nil
	},
}

func init() {
	Cmd.Flags().StringVar(&HubAddress, "hub-address", "127.0.0.1:8443", "IP address and port to connect to.")

	Cmd.AddCommand(daemonCmd)
}
