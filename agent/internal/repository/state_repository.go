package repository

import "agent/internal/model"

type StateRepository interface {
	Load() (*model.State, error)
	Save(state *model.State) error
}
