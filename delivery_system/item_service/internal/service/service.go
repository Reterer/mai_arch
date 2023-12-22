package service

import (
	"delivery_system/item_service/config"
)

type Repository interface {
}

type Service struct {
	cfg  *config.Service
	repo Repository
}

func New(cfg *config.Service, repo Repository) *Service {
	return &Service{
		cfg:  cfg,
		repo: repo,
	}
}
