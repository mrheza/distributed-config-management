package sqlite

import (
	"controller/internal/model"
	"database/sql"
	"errors"
)

type ConfigRepository struct {
	db *sql.DB
}

func NewConfigRepository(db *sql.DB) *ConfigRepository {
	return &ConfigRepository{db}
}

func (r *ConfigRepository) GetLatest() (*model.Config, error) {

	row := r.db.QueryRow(`
		SELECT version, url, poll_interval_seconds
		FROM configurations
		ORDER BY version DESC
		LIMIT 1
	`)

	var c model.Config

	err := row.Scan(
		&c.Version,
		&c.URL,
		&c.PollIntervalSeconds,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &c, nil
}

func (r *ConfigRepository) Create(url string, pollIntervalSeconds int) error {

	_, err := r.db.Exec(`
		INSERT INTO configurations (url, poll_interval_seconds)
		VALUES (?, ?)
	`,
		url,
		pollIntervalSeconds,
	)

	return err
}
