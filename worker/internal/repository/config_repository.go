package repository

import "worker/internal/model"

type ConfigRepository interface {
	Get() (*model.Config, error)
	Set(cfg *model.Config) error
}
