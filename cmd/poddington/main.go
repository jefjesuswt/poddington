package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/containers/podman/v4/pkg/bindings"
	"github.com/containers/podman/v4/pkg/bindings/containers"
)

func main() {
	fmt.Println("poddington init...")

	uid := os.Getuid()

	socketUrl := fmt.Sprintf("unix://run/user/%d/podman/podman.sock", uid)

	ctx := context.Background()

	connText, err := bindings.NewConnection(ctx, socketUrl)
	if err != nil {
		fmt.Printf("Error creating connection: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Connection established.\n")

	containerList, err := containers.List(connText, nil)
	if err != nil {
		fmt.Printf("Error listing containers: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("containers found: %d\n", len(containerList))
	fmt.Println(strings.Repeat("-", 50))

	for _, container := range containerList {
		name := container.Names[0]
		if len(name) > 0 && name[0] == '/' {
			name = name[1:]
		}
		fmt.Printf("ID: %s | Nombre: %s | Estado: %s\n", container.ID[:12], name, container.State)
	}
}
