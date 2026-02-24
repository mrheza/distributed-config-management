package service

import (
	"controller/internal/model"
	"controller/internal/repository"
)

type ConfigService interface {
	GetLatest() (*model.Config, error)
	Create(url string, pollIntervalSeconds int) error
}

type configService struct{ repo repository.ConfigRepository }

func NewConfigService(r repository.ConfigRepository) ConfigService {
	return &configService{r}
}

func (s *configService) GetLatest() (*model.Config, error) {
	return s.repo.GetLatest()
}

func (s *configService) Create(url string, pollIntervalSeconds int) error {
	return s.repo.Create(url, pollIntervalSeconds)
}
