package config

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Database struct {
	Read  *sql.DB
	Write *sql.DB
}

func NewSQLite(dbPath string) (*Database, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("could not create directory: %w", err)
	}

	dsn := fmt.Sprintf("file:%s?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)", dbPath)

	writePool, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	writePool.SetMaxOpenConns(1)

	readPool, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	readPool.SetMaxOpenConns(25)

	return &Database{Read: readPool, Write: writePool}, nil
}

func (db *Database) Close() error {
	var errs []error

	if err := db.Read.Close(); err != nil {
		errs = append(errs, fmt.Errorf("error closing read pool: %w", err))
	}

	if err := db.Write.Close(); err != nil {
		errs = append(errs, fmt.Errorf("error closing write pool: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to close database with errors: %v", errs)
	}

	return nil
}
