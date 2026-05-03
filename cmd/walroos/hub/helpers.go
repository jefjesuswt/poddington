package hub

import (
	"context"
	"os"
	"path/filepath"

	"github.com/jefjesuswt/walroos/config"
	"github.com/jefjesuswt/walroos/internal/fleet"
	"github.com/jefjesuswt/walroos/shared/ui"
)

func InitFleetService(ctx context.Context) (*fleet.Service, func() error, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, nil, ui.WrapError("error getting home directory: %w", err)
	}

	// ~/.config/walroos/hub.db
	dbPath := filepath.Join(homeDir, ".config", "walroos", "hub.db")

	client, err := config.NewSQLite(dbPath)
	if err != nil {
		return nil, nil, ui.WrapError("error creating sqlite client: %w", err)
	}

	store := fleet.NewSQLiteStore(client)
	svc := fleet.NewService(store)

	if err := svc.Init(ctx); err != nil {

		errs := client.Close()
		return nil, nil, ui.WrapError("error initializing fleet service: %w", errs)
	}

	cleanup := func() error {
		return client.Close()
	}

	return svc, cleanup, nil
}
