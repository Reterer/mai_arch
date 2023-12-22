package service

import (
	"delivery_system/item_service/config"
	"delivery_system/pkg/common_models"
	"errors"
)

var (
	AuthErr   = errors.New("auth err")
	NotExists = errors.New("not exists")
)

type ItemRepository interface {
	CreateItem(newItem common_models.CreateItemRequest) error
	GetItem(itemID common_models.ItemID) (common_models.Item, error)
	UpdateItem(item common_models.Item) error
	GetItems(userID common_models.UserID) ([]common_models.Item, error)
}

type UserRepository interface {
	GetUserIDByUsername(username string) (common_models.UserID, error)
}

type Service struct {
	cfg      *config.Service
	itemRepo ItemRepository
	userRepo UserRepository
}

func New(cfg *config.Service, itemRepo ItemRepository, userRepo UserRepository) *Service {
	return &Service{
		cfg:      cfg,
		itemRepo: itemRepo,
		userRepo: userRepo,
	}
}

func (s *Service) CreateItem(newItem common_models.CreateItemRequest) error {
	return s.itemRepo.CreateItem(newItem)
}

func (s *Service) GetItem(itemID common_models.ItemID) (common_models.Item, error) {
	return s.itemRepo.GetItem(itemID)
}

// Обновляет существующий item с itemID
func (s *Service) UpdateItem(updateItem common_models.Item) error {
	item, err := s.itemRepo.GetItem(updateItem.ItemID)
	if err != nil {
		// TODO no rows -> no exists
		return err
	}

	if item.ItemID != updateItem.ItemID || item.OwnerID != updateItem.OwnerID {
		return AuthErr
	}

	return s.itemRepo.UpdateItem(updateItem)
}

func (s *Service) GetItemsByUsername(username string) ([]common_models.Item, error) {
	// Нужно узнать userID
	userID, err := s.userRepo.GetUserIDByUsername(username)
	if err != nil {
		// TODO no exists -> invalid user
		return nil, err
	}

	// Затем находим поссылки по userID
	items, err := s.itemRepo.GetItems(userID)
	if err != nil {
		return nil, err
	}

	// Нужно
	return items, nil
}
