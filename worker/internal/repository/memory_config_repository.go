package repository

import (
	"database/sql"
	"sync"
	"worker/internal/model"
)

type MemoryConfigRepository struct {
	mu     sync.RWMutex
	config *model.Config
}

func NewMemoryConfigRepository() *MemoryConfigRepository {
	return &MemoryConfigRepository{}
}

func (r *MemoryConfigRepository) Get() (*model.Config, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.config == nil {
		return nil, sql.ErrNoRows
	}

	c := *r.config
	return &c, nil
}

func (r *MemoryConfigRepository) Set(cfg *model.Config) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if cfg == nil {
		r.config = nil
		return nil
	}

	c := *cfg
	r.config = &c
	return nil
}
