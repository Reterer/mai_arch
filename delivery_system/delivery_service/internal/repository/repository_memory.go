package repository

import (
	"delivery_system/delivery_service/config"
	"delivery_system/pkg/common_models"
	"sync"
)

type RepositoryMemory struct {
	cfg        *config.Repository
	users      []common_models.UserWithPass
	nextUserID common_models.UserID
	mu         sync.RWMutex
}

func NewMemory(cfg *config.Repository, initUsers []common_models.UserWithPass) *RepositoryMemory {
	users := make([]common_models.UserWithPass, 0, len(initUsers))
	nextUserID := common_models.UserID(0)
	for _, user := range initUsers {
		users = append(users, user)
		if nextUserID < user.UserID {
			nextUserID = user.UserID
		}
	}
	nextUserID++

	return &RepositoryMemory{
		cfg:        cfg,
		users:      users,
		nextUserID: nextUserID,
	}
}
