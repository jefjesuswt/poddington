package node

import (
	"context"

	"charm.land/lipgloss/v2"
	"github.com/jefjesuswt/poddington/internal/core/container"
	"github.com/jefjesuswt/poddington/internal/infra/podman"
	"github.com/jefjesuswt/poddington/internal/ui"
)

var ctx context.Context

func InitContainerService() (*container.Service, error) {
	client, err := podman.NewClient()
	if err != nil {
		return nil, ui.PrintError("error creating podman client: %w", err)
	}
	return container.NewService(client), nil
}

func WithContainerAction(
	args []string,
	action func(ctx context.Context, svc *container.Service, target, highlighted string) error,
) error {
	svc, err := InitContainerService()
	if err != nil {
		return err
	}
	target := args[0]
	highlighted := lipgloss.NewStyle().Foreground(ui.PodFrosted).Bold(true).Render(target)

	return action(ctx, svc, target, highlighted)
}
