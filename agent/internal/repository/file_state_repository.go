package repository

import (
	"agent/internal/model"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type FileStateRepository struct {
	path string
}

func NewFileStateRepository(path string) *FileStateRepository {
	return &FileStateRepository{path: path}
}

func (r *FileStateRepository) Load() (*model.State, error) {
	raw, err := os.ReadFile(r.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &model.State{}, nil
		}
		return nil, err
	}

	var s model.State
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *FileStateRepository) Save(state *model.State) error {
	if err := os.MkdirAll(filepath.Dir(r.path), 0o755); err != nil {
		return err
	}

	raw, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.path, raw, 0o644)
}
