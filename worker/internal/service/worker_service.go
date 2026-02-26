package service

import (
	"context"
	"worker/internal/client"
	"worker/internal/model"
	"worker/internal/repository"
)

type WorkerService interface {
	ApplyConfig(cfg *model.Config) error
	Hit(ctx context.Context) (status int, contentType string, body []byte, err error)
	GetCurrentConfig() (*model.Config, error)
}

type workerService struct {
	repo  repository.ConfigRepository
	fetch client.FetchClient
}

func NewWorkerService(repo repository.ConfigRepository, fetch client.FetchClient) WorkerService {
	return &workerService{repo: repo, fetch: fetch}
}

func (s *workerService) ApplyConfig(cfg *model.Config) error {
	return s.repo.Set(cfg)
}

func (s *workerService) Hit(ctx context.Context) (int, string, []byte, error) {
	cfg, err := s.repo.Get()
	if err != nil {
		return 0, "", nil, err
	}

	return s.fetch.Get(ctx, cfg.URL)
}

func (s *workerService) GetCurrentConfig() (*model.Config, error) {
	return s.repo.Get()
}
