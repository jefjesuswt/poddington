package config

import (
	"context"
	"fmt"
	"os"

	"github.com/containers/podman/v4/pkg/bindings"
)

type PodmanClient struct {
	// not a normal ctx, it's a custom podman ctx with parameters already injected
	// contains http and socket already inserted, podman api requires using it
	Ctx context.Context
}

func NewPodmanClient() (*PodmanClient, error) {
	uid := os.Getuid()
	socketUrl := fmt.Sprintf("unix://run/user/%d/podman/podman.sock", uid)
	ctx := context.Background()

	conn, err := bindings.NewConnection(ctx, socketUrl)
	if err != nil {
		return nil, fmt.Errorf("error creating connection: %w", err)
	}

	return &PodmanClient{
		Ctx: conn,
	}, nil
}
