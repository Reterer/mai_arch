package repository

import (
	"delivery_system/item_service/config"
	"delivery_system/pkg/common_models"
)

type UserMock struct {
	cfg             *config.Repository
	usersByUsername map[string]common_models.User
}

func NewUserMock(cfg *config.Repository, users []common_models.User) *UserMock {
	userMock := &UserMock{
		cfg:             cfg,
		usersByUsername: make(map[string]common_models.User),
	}

	for _, user := range users {
		userMock.usersByUsername[user.Usernmame] = user
	}

	return userMock
}

func (u *UserMock) GetUserIDByUsername(username string) (common_models.UserID, error) {
	user, ok := u.usersByUsername[username]
	if !ok {
		return 0, NotExistsErr
	}
	return user.UserID, nil
}
