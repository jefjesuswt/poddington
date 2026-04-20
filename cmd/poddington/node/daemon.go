package node

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jefjesuswt/poddington/shared/ui"
	"github.com/spf13/cobra"
)

var (
	nodeID    string
	nodeToken string
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Starts the Poddington Node background worker",
	RunE: func(cmd *cobra.Command, args []string) error {
		if nodeID == "" || nodeToken == "" {
			return ui.WrapError("node-id and node-token are required to start the daemon")
		}

		ui.PrintTitle("Poddington Node Daemon Initialization")
		slog.Info("Starting background worker...", "hub", HubAddress, "id", nodeID)

		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		sendPing()

		for {
			select {
			case <-cmd.Context().Done():
				slog.Info("Daemon shutting down...")

			case <-ticker.C:
				sendPing()
			}
		}

	},
}

func sendPing() {
	url := fmt.Sprintf("http://%s/api/nodes/%s/ping", HubAddress, nodeID)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(nil))
	if err != nil {
		slog.Error("Failed to ping hub", "hub", HubAddress, "error", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+nodeToken)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Failed to ping hub", "hub", HubAddress, "error", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		slog.Info("Ping succesful", "hub", HubAddress)
	} else {
		slog.Error("Failed to ping hub", "hub", HubAddress, "status", resp.StatusCode)
	}

}

func init() {
	daemonCmd.Flags().StringVar(&nodeID, "id", "", "The ID of this node")
	daemonCmd.Flags().StringVar(&nodeToken, "token", "", "The security token for this node")

	daemonCmd.MarkFlagRequired("id")
	daemonCmd.MarkFlagRequired("token")
}
