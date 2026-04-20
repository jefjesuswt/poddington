package hub

import (
	"context"
	"os"
	"path/filepath"

	"github.com/jefjesuswt/poddington/config"
	"github.com/jefjesuswt/poddington/internal/fleet"
	"github.com/jefjesuswt/poddington/shared/ui"
)

func InitFleetService(ctx context.Context) (*fleet.Service, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, ui.WrapError("error getting home directory: %w", err)
	}

	// ~/.config/poddington/hub.db
	dbPath := filepath.Join(homeDir, ".config", "poddington", "hub.db")

	client, err := config.NewSQLite(dbPath)
	if err != nil {
		return nil, ui.WrapError("error creating sqlite client: %w", err)
	}

	repo := fleet.NewRepository(client)
	svc := fleet.NewService(repo)

	// auto-migrations
	if err := svc.Init(ctx); err != nil {
		return nil, ui.WrapError("error initializing fleet service: %w", err)
	}

	return svc, nil
}
