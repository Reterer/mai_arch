package repository

import "delivery_system/item_service/config"

type Repository struct {
	cfg *config.Repository
}

func New(cfg *config.Repository) *Repository {
	return &Repository{
		cfg: cfg,
	}
}
