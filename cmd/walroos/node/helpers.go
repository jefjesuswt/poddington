package node

import (
	"context"

	"github.com/jefjesuswt/walroos/config"
	"github.com/jefjesuswt/walroos/internal/container"
	"github.com/jefjesuswt/walroos/shared/ui"
)

func InitContainerService(ctx context.Context) (*container.Service, error) {
	client, err := config.NewPodmanClient()
	if err != nil {
		return nil, ui.WrapError("error creating podman client: %w", err)
	}
	podmanRepo := container.NewRepository(client)
	return container.NewService(podmanRepo), nil
}
