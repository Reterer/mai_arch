package repository

import (
	"delivery_system/item_service/config"
	"errors"
)

var (
	NotExistsErr = errors.New("no exists")
)

type Repository struct {
	cfg *config.Repository
}

func New(cfg *config.Repository) *Repository {
	return &Repository{
		cfg: cfg,
	}
}
