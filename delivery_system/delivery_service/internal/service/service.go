package service

import (
	"delivery_system/delivery_service/config"
	"delivery_system/pkg/common_models"
)

type Repository interface {
	CreateDelivery(req common_models.CreateDeliveryRequest) error
	GetDeliveriesByFrom(userID common_models.UserID) ([]common_models.Delivery, error)
	GetDeliveriesByTo(userID common_models.UserID) ([]common_models.Delivery, error)
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

func (s *Service) CreateDelivery(req common_models.CreateDeliveryRequest) error {
	return s.repo.CreateDelivery(req)
}

func (s *Service) GetDeliveriesByFrom(userID common_models.UserID) ([]common_models.Delivery, error) {
	return s.repo.GetDeliveriesByFrom(userID)
}
func (s *Service) GetDeliveriesByTo(userID common_models.UserID) ([]common_models.Delivery, error) {
	return s.repo.GetDeliveriesByTo(userID)
}
