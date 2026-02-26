package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func New(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS agents (
			id TEXT PRIMARY KEY,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		return nil, fmt.Errorf("create agents table: %w", err)
	}

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS configurations (
			version BIGSERIAL PRIMARY KEY,
			url TEXT,
			poll_interval_seconds INTEGER NOT NULL DEFAULT 30,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		return nil, fmt.Errorf("create configurations table: %w", err)
	}

	return db, nil
}
