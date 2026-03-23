package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type DB struct {
	ReadPool *sql.DB
	WritePool *sql.DB
}

func NewDB(dbPath string) (*DB, error) {
	if err := os.Mkdir(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("could not create directory: %w", err)
	}

	// dsn with pragmas for concurrence and tolerance
	dsn := fmt.Sprintf("file:%s?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)", dbPath)

	writePool, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("could not open write db: %w", err)
	}

	writePool.SetMaxOpenConns(1)

	readPool, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("could not open read db: %w", err)
	}

	readPool.SetMaxOpenConns(25)

	if err := writePool.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping write db: %w", err)
	}

	if err := readPool.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping read db: %w", err)
	}

	if err := runMigrations(writePool); err != nil {
		return nil, fmt.Errorf("could not run migrations: %w", err)
	}

	return &DB{
		ReadPool: readPool,
		WritePool: writePool,
	}, nil
}

func runMigrations(db *sql.DB) error {
	const createNodesTables = `
	BEGIN IMMEDIATE;
		CREATE TABLE IF NOT EXISTS nodes (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			address TEXT NOT NULL,
			token TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			last_seen DATETIME NOT NULL
		);
		COMMIT;
	`

	_, err := db.ExecContext(context.Background(), createNodesTables)
	return err
}
