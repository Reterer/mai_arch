package repository

import "delivery_system/client_service/config"

type Repository struct {
	cfg *config.Repository
}

func New(cfg *config.Repository) *Repository {
	return &Repository{
		cfg: cfg,
	}
}
