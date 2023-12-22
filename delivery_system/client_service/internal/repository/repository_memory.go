package repository

import (
	"delivery_system/client_service/config"
	"delivery_system/pkg/common_models"
	"errors"
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

func (r *RepositoryMemory) findByUsername(username string) (int, bool) {
	for i := range r.users {
		if r.users[i].Usernmame == username {
			return i, true
		}
	}
	return 0, false
}
func (r *RepositoryMemory) findByUserID(userID common_models.UserID) (int, bool) {
	for i := range r.users {
		if r.users[i].UserID == userID {
			return i, true
		}
	}
	return 0, false
}

func (r *RepositoryMemory) genUserID() common_models.UserID {
	next := r.nextUserID
	r.nextUserID++
	return next
}

func (r *RepositoryMemory) AddUser(req common_models.RegisterUserRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.findByUsername(req.Usernmame)
	if ok {
		return errors.New("exists")
	}

	r.users = append(r.users, common_models.UserWithPass{
		User: common_models.User{
			UserID:    r.genUserID(),
			Usernmame: req.Usernmame,
			FirstName: req.FirstName,
			LastName:  req.LastName,
		},
		Passhash: req.Password,
	})
	return nil
}

func (r *RepositoryMemory) GetUserByID(userID common_models.UserID) (common_models.UserWithPass, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	idx, ok := r.findByUserID(userID)
	if !ok {
		return common_models.UserWithPass{}, errors.New("not found")
	}
	return r.users[idx], nil
}
func (r *RepositoryMemory) UpdateUser(req common_models.UpdateUserRequest) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	idx, ok := r.findByUserID(req.UserID)
	if !ok {
		return false, errors.New("not found")
	}

	if req.Usernmame != nil {
		if _, ok := r.findByUsername(*req.Usernmame); ok {
			return false, nil
		}
		r.users[idx].Usernmame = *req.Usernmame
	}
	if req.FirstName != nil {
		r.users[idx].FirstName = *req.FirstName
	}
	if req.LastName != nil {
		r.users[idx].LastName = *req.LastName
	}
	if req.Password != nil {
		r.users[idx].Passhash = *req.Password
	}
	return true, nil
}

func (r *RepositoryMemory) SearchMask(mask string) ([]common_models.UserWithPass, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	panic("не реализовано")
}

func (r *RepositoryMemory) GetUserByUsername(username string) (common_models.UserWithPass, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	idx, ok := r.findByUsername(username)
	if !ok {
		return common_models.UserWithPass{}, false, nil
	}
	return r.users[idx], true, nil
}
