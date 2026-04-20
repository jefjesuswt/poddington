package fleet

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jefjesuswt/poddington/config"
)

type Repository struct {
	db *config.Database
}

func NewRepository(db *config.Database) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Migrate(ctx context.Context) error {
	const query = `
	BEGIN IMMEDIATE;
	CREATE TABLE IF NOT EXISTS fleet_nodes (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		address TEXT NOT NULL,
		token TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		last_seen DATETIME NOT NULL
	);
	COMMIT;
	`

	_, err := r.db.Write.ExecContext(ctx, query)
	return err
}

func (r *Repository) Save(ctx context.Context, n Node) error {
	const query = `
		INSERT INTO fleet_nodes (id, name, address, token, created_at, last_seen)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Write.ExecContext(
		ctx,
		query,
		n.ID,
		n.Name,
		n.Address,
		n.Token,
		n.CreatedAt,
		n.LastSeen,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: fleet_nodes.name") {
			return ErrNodeAlreadyExists
		}
		return fmt.Errorf("failed to save node: %w", err)
	}

	return nil
}

func (r *Repository) GetAll(ctx context.Context) ([]Node, error) {
	const query = `SELECT id, name, address, token, created_at, last_seen FROM fleet_nodes`

	rows, err := r.db.Read.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes: %w", err)
	}
	defer rows.Close()

	var nodes []Node
	for rows.Next() {
		var n Node
		if err := rows.Scan(
			&n.ID,
			&n.Name,
			&n.Address,
			&n.Token,
			&n.CreatedAt,
			&n.LastSeen,
		); err != nil {
			return nil, fmt.Errorf("failed to scan node: %w", err)
		}
		nodes = append(nodes, n)
	}

	return nodes, rows.Err()
}

func (r *Repository) GetByID(ctx context.Context, id string) (Node, error) {
	const query = `SELECT id, name, address, token, created_at, last_seen FROM fleet_nodes WHERE id = ?`

	var n Node

	if err := r.db.Read.QueryRowContext(ctx, query, id).Scan(
		&n.ID,
		&n.Name,
		&n.Address,
		&n.Token,
		&n.CreatedAt,
		&n.LastSeen,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Node{}, ErrNodeNotFound
		}

		return Node{}, fmt.Errorf("failed to get node: %w", err)
	}

	return n, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	const query = `DELETE FROM fleet_nodes WHERE id = ?`

	result, err := r.db.Write.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNodeNotFound
	}
	return nil
}

func (r *Repository) UpdateLastSeen(ctx context.Context, id string) error {
	const query = `UPDATE fleet_nodes SET last_seen = ? WHERE id = ?`

	result, err := r.db.Write.ExecContext(ctx, query, time.Now().UTC(), id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNodeNotFound
	}

	return nil
}
