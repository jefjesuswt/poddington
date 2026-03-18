package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jefjesuswt/poddington/internal/core/container"
	"github.com/jefjesuswt/poddington/internal/infra/podman"
)

func main() {

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]


	podmanClient, err := podman.NewClient()
	if err != nil {
		fmt.Printf("Error creating podman client: %v\n", err)
		os.Exit(1)
	}
	containerService := container.NewService(podmanClient)
	ctx := context.Background()

	switch command {
	case "list":
		instances, err := containerService.ListContainers(ctx)
		if err != nil {
			fmt.Printf("Error listing containers: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("containers found: %d\n", len(instances))
		for _, instance := range instances {
			fmt.Printf("ID: %s | Name: %s | State: %s\n", instance.ID[:12], instance.Name, instance.State)
		}

	case "inspect":
		target := getTargetArgs()
		instance, err := containerService.InspectContainer(ctx, target)
		if err != nil {
			fmt.Printf("Error inspecting container: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("container id: %s\n", instance.ID)
		fmt.Printf("container name: %s\n", instance.Name)
		fmt.Printf("container image: %s\n", instance.Image)
		fmt.Printf("container state: %s\n", instance.State)

	case "start":
		target := getTargetArgs()
		err := containerService.StartContainer(ctx, target)
		if err != nil {
			fmt.Printf("Error starting container: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("started container %s\n", target)

	case "stop":
		target := getTargetArgs()
		err := containerService.StopContainer(ctx, target)
		if err != nil {
			fmt.Printf("Error stopping container: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("stopped container %s\n", target)

	default:
		fmt.Printf("unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Uso de Poddington Mini-CLI:")
	fmt.Println("  go run cmd/poddington/main.go list")
	fmt.Println("  go run cmd/poddington/main.go inspect <container_name>")
	fmt.Println("  go run cmd/poddington/main.go start <container_name>")
	fmt.Println("  go run cmd/poddington/main.go stop <container_name>")
}

func getTargetArgs() string {
	if len(os.Args) < 3 {
		fmt.Println("Error: missing target argument")
		printUsage()
		os.Exit(1)
	}
	return os.Args[2]
}
