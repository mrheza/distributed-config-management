package repository

import (
	"database/sql"
	"errors"
	"testing"
	"worker/internal/model"

	"github.com/stretchr/testify/assert"
)

func TestMemoryConfigRepository_Get_NoConfig(t *testing.T) {
	repo := NewMemoryConfigRepository()

	cfg, err := repo.Get()
	assert.Nil(t, cfg)
	assert.True(t, errors.Is(err, sql.ErrNoRows))
}

func TestMemoryConfigRepository_SetAndGet(t *testing.T) {
	repo := NewMemoryConfigRepository()

	err := repo.Set(&model.Config{Version: 1, URL: "https://example.com"})
	assert.NoError(t, err)

	cfg, err := repo.Get()
	assert.NoError(t, err)
	assert.Equal(t, 1, cfg.Version)
	assert.Equal(t, "https://example.com", cfg.URL)
}

func TestMemoryConfigRepository_Get_ReturnsClone(t *testing.T) {
	repo := NewMemoryConfigRepository()
	_ = repo.Set(&model.Config{Version: 1, URL: "https://example.com"})

	cfg, err := repo.Get()
	assert.NoError(t, err)
	cfg.URL = "changed"

	cfg2, err := repo.Get()
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", cfg2.URL)
}
