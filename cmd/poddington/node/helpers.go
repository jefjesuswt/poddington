package node

import (
	"context"

	"charm.land/lipgloss/v2"
	"github.com/jefjesuswt/poddington/cmd/config"
	"github.com/jefjesuswt/poddington/internal/container"
	"github.com/jefjesuswt/poddington/shared/ui"
)


func InitContainerService(ctx context.Context) (*container.Service, error) {
	client, err := config.NewPodmanClient()
	if err != nil {
		return nil, ui.PrintError("error creating podman client: %w", err)
	}
	podmanRepo := container.NewRepository(client)
	return container.NewService(podmanRepo), nil
}

func WithContainerAction(
	ctx context.Context,
	args []string,
	action func(ctx context.Context, svc *container.Service, target, highlighted string) error,
) error {
	svc, err := InitContainerService(ctx)
	if err != nil {
		return err
	}
	target := args[0]
	highlighted := lipgloss.NewStyle().Foreground(ui.PodFrosted).Bold(true).Render(target)

	return action(ctx, svc, target, highlighted)
}
