package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	clientMocks "worker/internal/mocks/client"
	repositoryMocks "worker/internal/mocks/repository"
	"worker/internal/model"

	"github.com/stretchr/testify/assert"
)

func TestWorkerService_ApplyConfig(t *testing.T) {
	repo := new(repositoryMocks.ConfigRepository)
	fetch := new(clientMocks.FetchClient)
	svc := NewWorkerService(repo, fetch)

	cfg := &model.Config{Version: 1, URL: "https://example.com"}
	repo.On("Set", cfg).Return(nil).Once()

	err := svc.ApplyConfig(cfg)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestWorkerService_Hit_NoConfig(t *testing.T) {
	repo := new(repositoryMocks.ConfigRepository)
	fetch := new(clientMocks.FetchClient)
	svc := NewWorkerService(repo, fetch)

	repo.On("Get").Return((*model.Config)(nil), sql.ErrNoRows).Once()

	_, _, _, err := svc.Hit(context.Background())
	assert.Error(t, err)
	assert.True(t, errors.Is(err, sql.ErrNoRows))
}

func TestWorkerService_Hit_Success(t *testing.T) {
	repo := new(repositoryMocks.ConfigRepository)
	fetch := new(clientMocks.FetchClient)
	svc := NewWorkerService(repo, fetch)

	cfg := &model.Config{Version: 1, URL: "https://example.com"}
	repo.On("Get").Return(cfg, nil).Once()
	fetch.On("Get", context.Background(), cfg.URL).Return(200, "text/plain", []byte("ok"), nil).Once()

	status, ct, body, err := svc.Hit(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 200, status)
	assert.Equal(t, "text/plain", ct)
	assert.Equal(t, []byte("ok"), body)
}

func TestWorkerService_GetCurrentConfig(t *testing.T) {
	repo := new(repositoryMocks.ConfigRepository)
	fetch := new(clientMocks.FetchClient)
	svc := NewWorkerService(repo, fetch)

	cfg := &model.Config{Version: 2, URL: "https://example.com/v2"}
	repo.On("Get").Return(cfg, nil).Once()

	result, err := svc.GetCurrentConfig()
	assert.NoError(t, err)
	assert.Equal(t, cfg, result)
}
