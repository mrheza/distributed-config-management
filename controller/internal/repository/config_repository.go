package repository

import (
	"controller/internal/model"
	"errors"
)

var ErrConfigNotFound = errors.New("config not found")

type ConfigRepository interface {
	GetLatest() (*model.Config, error)
	Create(url string, pollIntervalSeconds int) error
}
