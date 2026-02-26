package service

import (
	"controller/internal/model"
	"controller/internal/repository"
	"sync"
)

type ConfigService interface {
	GetLatest() (*model.Config, error)
	Create(url string, pollIntervalSeconds int) error
}

type configService struct {
	repo repository.ConfigRepository

	mu       sync.RWMutex
	latest   *model.Config
	hasCache bool
}

func NewConfigService(r repository.ConfigRepository) ConfigService {
	return &configService{repo: r}
}

func (s *configService) GetLatest() (*model.Config, error) {
	s.mu.RLock()
	if s.hasCache && s.latest != nil {
		cached := cloneConfig(s.latest)
		s.mu.RUnlock()
		return cached, nil
	}
	s.mu.RUnlock()

	cfg, err := s.repo.GetLatest()
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.latest = cloneConfig(cfg)
	s.hasCache = true
	s.mu.Unlock()

	return cloneConfig(cfg), nil
}

func (s *configService) Create(url string, pollIntervalSeconds int) error {
	if err := s.repo.Create(url, pollIntervalSeconds); err != nil {
		return err
	}

	latest, err := s.repo.GetLatest()
	if err != nil {
		return err
	}

	// Config changed, refresh cache immediately with the latest DB value.
	s.mu.Lock()
	s.latest = cloneConfig(latest)
	s.hasCache = true
	s.mu.Unlock()

	return nil
}

func cloneConfig(c *model.Config) *model.Config {
	if c == nil {
		return nil
	}

	cp := *c
	return &cp
}
