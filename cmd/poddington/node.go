package main

import (
	"context"
	"errors"
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/log"
	"github.com/jefjesuswt/poddington/internal/core/container"
	"github.com/jefjesuswt/poddington/internal/infra/podman"
	"github.com/jefjesuswt/poddington/internal/ui"
	"github.com/spf13/cobra"
)

var hubAddress string
var listAll bool
var forceRemove bool
var ctx = context.Background()

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Run a node that connects to the hub",
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.PrintTitle("Poddington Node Initialization")
		log.Info("Poddington starting...", "mode", "NODE")
		log.Info("Connecting to hub...", "hub", hubAddress)
		return nil
	},
}

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "List node's local containers",
	RunE: func(cmd *cobra.Command, args []string) error {
		return withContainerAction(args, func(ctx context.Context, svc *container.Service, target, highlighted string) error {
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

var inspectCommand = &cobra.Command{
	Use:   "inspect [container_name]",
	Short: "Inspect a container, gets deep technical details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return withContainerAction(args, func(ctx context.Context, svc *container.Service, target, highlighted string) error {
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

var startCommand = &cobra.Command{
	Use:   "start [container_name]",
	Short: "Start a container",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return withContainerAction(args, func(ctx context.Context, svc *container.Service, target, highlighted string) error {
			if err := svc.Start(ctx, target); err != nil {
				if errors.Is(err, container.ErrAlreadyRunning) {
					ui.PrintWarning("Container %s is already running.", highlighted)
					return nil
				}
				return ui.PrintError("error starting container: %w", err)
			}
			ui.PrintSuccess("Container %s is now running.", highlighted)
			return nil
		})
	},
}

var stopCommand = &cobra.Command{
	Use:   "stop [container_name]",
	Short: "Stops a container",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return withContainerAction(args, func(ctx context.Context, svc *container.Service, target, highlighted string) error {
			if err := svc.Stop(ctx, target); err != nil {
				if errors.Is(err, container.ErrAlreadyStopped) {
					ui.PrintWarning("Container %s is already stopped.", highlighted)
					return nil
				}
				return ui.PrintError("error stopping container: %w", err)
			}

			ui.PrintSuccess("Container %s has been stopped.", highlighted)
			return nil
		})

	},
}

var restartCommand = &cobra.Command{
	Use: "restart [container_name]",
	Short: "Restarts a container",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return withContainerAction(args, func(ctx context.Context, svc *container.Service, target, highlighted string) error {
			if err := svc.Restart(ctx, target); err != nil {
				return ui.PrintError("error restarting container: %w", err)
			}
			ui.PrintSuccess("Container %s has been restarted.", highlighted)
			return nil
		})
	},
}

var rmCommand = &cobra.Command{
	Use: "rm [container_name]",
	Short: "Removes a container",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return withContainerAction(args, func(ctx context.Context, svc *container.Service, target, highlighted string) error {
			if err := svc.Remove(ctx, target, forceRemove); err != nil {
				return ui.PrintError("error removing container: %w", err)
			}
			ui.PrintSuccess("Container %s has been removed.", highlighted)
			return nil
		})
	},
}

var logsCommand = &cobra.Command{
	Use: "logs [container_name]",
	Short: "Fetch the logs of a container",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return withContainerAction(args, func(ctx context.Context, svc *container.Service, target, higlighted string) error {
			logsData, err := svc.GetLogs(ctx, target)
			if err != nil {
				return ui.PrintError("error fetching logs: %w", err)
			}
			ui.PrintLogs(target, logsData)
			return nil
		})
	},
}

func init() {
	nodeCmd.Flags().StringVar(&hubAddress, "hub-address", "127.0.0.1:8443", "IP address and port to connect to.")

	nodeCmd.AddCommand(listCommand)
	listCommand.Flags().BoolVarP(&listAll, "all", "a", false, "Show all containers (default shows just running)")

	nodeCmd.AddCommand(rmCommand)
	rmCommand.Flags().BoolVarP(&forceRemove, "force", "f", false, "Force the removal of an already container")

	nodeCmd.AddCommand(inspectCommand)
	nodeCmd.AddCommand(startCommand)
	nodeCmd.AddCommand(stopCommand)
	nodeCmd.AddCommand(restartCommand)
	nodeCmd.AddCommand(logsCommand)
}

func initContainerService() (*container.Service, error) {
	client, err := podman.NewClient()
	if err != nil {
		return nil, ui.PrintError("error creating podman client: %w", err)
	}
	return container.NewService(client), nil
}

func withContainerAction(
	args []string,
	action func(ctx context.Context, svc *container.Service, target, highlighted string) error,
) error {
	svc, err := initContainerService()
	if err != nil {
		return err
	}
	target := args[0]
	highlighted := lipgloss.NewStyle().Foreground(ui.PodFrosted).Bold(true).Render(target)

	return action(ctx, svc, target, highlighted)
}
