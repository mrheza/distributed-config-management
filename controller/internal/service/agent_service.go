package service

import (
	"controller/internal/repository"

	"github.com/google/uuid"
)

type AgentService interface {
	Register(existingID string) (string, error)
}

type agentService struct{ repo repository.AgentRepository }

func NewAgentService(r repository.AgentRepository) AgentService {
	return &agentService{repo: r}
}

func (s *agentService) Register(existingID string) (string, error) {
	id := existingID
	if _, err := uuid.Parse(existingID); existingID == "" || err != nil {
		id = uuid.New().String()
	}
	return id, s.repo.Save(id)
}
