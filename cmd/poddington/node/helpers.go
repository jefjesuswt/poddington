package node

import (
	"context"

	"github.com/jefjesuswt/poddington/config"
	"github.com/jefjesuswt/poddington/internal/container"
	"github.com/jefjesuswt/poddington/shared/ui"
)

func InitContainerService(ctx context.Context) (*container.Service, error) {
	client, err := config.NewPodmanClient()
	if err != nil {
		return nil, ui.WrapError("error creating podman client: %w", err)
	}
	podmanRepo := container.NewRepository(client)
	return container.NewService(podmanRepo), nil
}
