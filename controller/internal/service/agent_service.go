package service

import (
	"controller/internal/repository"

	"github.com/google/uuid"
)

type AgentService interface {
	Register() (string, error)
}

type agentService struct{ repo repository.AgentRepository }

func NewAgentService(r repository.AgentRepository) AgentService {
	return &agentService{r}
}

func (s *agentService) Register() (string, error) {
	id := uuid.New().String()
	return id, s.repo.Save(id)
}
