package sqlite

import "database/sql"

type AgentRepository struct{ db *sql.DB }

func NewAgentRepository(db *sql.DB) *AgentRepository {
	return &AgentRepository{db}
}

func (r *AgentRepository) Save(id string) error {
	_, err := r.db.Exec(`
		INSERT OR IGNORE INTO agents (id)
		VALUES (?)
	`, id)

	return err
}
