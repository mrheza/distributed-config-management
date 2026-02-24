package repository

type AgentRepository interface {
	Save(id string) error
}
