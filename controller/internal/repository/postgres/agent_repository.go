package postgres

import "database/sql"

type AgentRepository struct{ db *sql.DB }

func NewAgentRepository(db *sql.DB) *AgentRepository {
	return &AgentRepository{db}
}

func (r *AgentRepository) Save(id string) error {
	_, err := r.db.Exec(`
		INSERT INTO agents (id)
		VALUES ($1) ON CONFLICT (id) DO NOTHING
	`, id)

	return err
}
