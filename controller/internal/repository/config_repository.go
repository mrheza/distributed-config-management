package repository

import (
	"controller/internal/model"
)

type ConfigRepository interface {
	GetLatest() (*model.Config, error)
	Create(url string, pollIntervalSeconds int) error
}
