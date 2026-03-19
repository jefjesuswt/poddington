package main

import (
	"context"
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/log"
	"github.com/jefjesuswt/poddington/internal/core/container"
	"github.com/jefjesuswt/poddington/internal/infra/podman"
	"github.com/jefjesuswt/poddington/internal/ui"
	"github.com/spf13/cobra"
)

var hubAddress string

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Run a node that connects to the hub",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("Poddington starting...", "mode", "NODE")
		log.Info("Connecting to hub...", "hub", hubAddress)
		return nil
	},
}

var listAll bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List node's local containers",
	RunE: func(cmd *cobra.Command, args []string) error {
		svc, err := initContainerService()
		if err != nil {
			return err
		}

		instances, err := svc.ListContainers(context.Background(), listAll)
		if err != nil {
			return err
		}

		title := "Active containers"
		if listAll {
			title = "All containers (Including stopped)"
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
	},
}

var inspectCommand = &cobra.Command{
	Use:   "inspect [container_name]",
	Short: "Inspect a container, gets deep technical details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		svc, err := initContainerService()
		if err != nil {
			return err
		}

		instance, err := svc.InspectContainer(context.Background(), args[0])
		if err != nil {
			return err
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
	},
}

var startCommand = &cobra.Command{
	Use:   "start [container_name]",
	Short: "Start a container",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		svc, err := initContainerService()
		if err != nil {
			return err
		}

		target := args[0]
		err = svc.StartContainer(context.Background(), target)
		if err != nil {
			return ui.PrintError("error starting container: %w", err)
		}

		highlightedName := lipgloss.NewStyle().Foreground(ui.ColorText).Bold(true).Render(target)
		ui.PrintSuccess("Container %s is now running.", highlightedName)

		return nil
	},
}

var stopCommand = &cobra.Command{
	Use:   "stop [container_name]",
	Short: "Stops a container",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		svc, err := initContainerService()
		if err != nil {
			return err
		}

		target := args[0]
		if err := svc.StopContainer(context.Background(), target); err != nil {
			return ui.PrintError("error stopping container: %w", err)
		}

		highlightedName := lipgloss.NewStyle().Foreground(ui.ColorText).Bold(true).Render(target)
		ui.PrintInfo("Container %s has been stopped.", highlightedName)

		return nil
	},
}

func init() {
	nodeCmd.Flags().StringVar(&hubAddress, "hub-address", "127.0.0.1:8443", "IP address and port to connect to.")

	nodeCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "Show all containers (default shows just running)")
	nodeCmd.AddCommand(inspectCommand)
	nodeCmd.AddCommand(startCommand)
	nodeCmd.AddCommand(stopCommand)
}

func initContainerService() (*container.Service, error) {
	client, err := podman.NewClient()
	if err != nil {
		return nil, ui.PrintError("error creating podman client: %w", err)
	}
	return container.NewService(client), nil
}
