package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func New(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	db.Exec(`CREATE TABLE IF NOT EXISTS 
			agents (
    			id TEXT PRIMARY KEY,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP
			)`)

	db.Exec(`CREATE TABLE IF NOT EXISTS 
			configurations(
				version INTEGER PRIMARY KEY AUTOINCREMENT,
				url TEXT,
				poll_interval_seconds INTEGER NOT NULL DEFAULT 30,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP
			)`)

	return db, nil
}
