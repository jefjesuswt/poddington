package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jefjesuswt/walroos/internal/fleet"
	"github.com/jefjesuswt/walroos/shared/ui"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/spf13/cobra"
)

var (
	nodeID    string
	nodeToken string
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Starts the Walroos Node background worker",
	RunE: func(cmd *cobra.Command, args []string) error {
		if nodeID == "" || nodeToken == "" {
			return ui.WrapError("node-id and node-token are required to start the daemon")
		}

		ui.PrintTitle("Walroos Node Daemon Initialization")
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

func init() {
	daemonCmd.Flags().StringVar(&nodeID, "id", "", "The ID of this node")
	daemonCmd.Flags().StringVar(&nodeToken, "token", "", "The security token for this node")

	daemonCmd.MarkFlagRequired("id")
	daemonCmd.MarkFlagRequired("token")
}

func sendPing() {
	url := fmt.Sprintf("http://%s/api/nodes/%s/ping", HubAddress, nodeID)

	tel, err := getTelemetry()
	if err != nil {
		slog.Error("Failed to get telemetry", "hub", HubAddress, "error", err)
		return
	}

	jsonTel, err := json.Marshal(tel)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonTel))
	if err != nil {
		slog.Error("Failed to ping hub", "hub", HubAddress, "error", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
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

func getTelemetry() (*fleet.NodeTelemtry, error) {
	vMem, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	cpuPercents, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}

	return &fleet.NodeTelemtry{
		CpuPercentUsage: cpuPercents[0],
		MemPercentUsage: vMem.UsedPercent,
		MemTotalMB:      vMem.Total / 1024 / 1024,
		MemUsedMB:       vMem.Used / 1024 / 1024,
	}, nil
}
