package service

import (
	"delivery_system/client_service/config"
	"delivery_system/pkg/common_models"

	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	AddUser(req common_models.RegisterUserRequest) error
	GetUserByID(user_id common_models.UserID) (common_models.UserWithPass, error)
	UpdateUser(req common_models.UpdateUserRequest) (bool, error)
	SearchMask(mask string) ([]common_models.UserWithPass, error)
	GetUserByUsername(username string) (common_models.UserWithPass, bool, error)
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

// Создает нового пользователя
func (s *Service) Register(req common_models.RegisterUserRequest) error {
	var err error
	req.Password, err = s.hashPassword(req.Password)
	if err != nil {
		return err
	}
	return s.repo.AddUser(req)
}

// Возвращает пользователя по user_id
func (s *Service) GetUser(user_id common_models.UserID) (common_models.User, error) {
	user, err := s.repo.GetUserByID(user_id)
	if err != nil {
		return common_models.User{}, err
	}
	return common_models.User{
		UserID:    user.UserID,
		Usernmame: user.Usernmame,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}

// Обновляет не нулевые поля
// возвращает true, если успешно обновлен
// false, если нет (username занят)
// error - если произошла какая-то ошибка
func (s *Service) UpdateUser(req common_models.UpdateUserRequest) (bool, error) {
	if req.Password != nil {
		var err error
		*req.Password, err = s.hashPassword(*req.Password)
		if err != nil {
			return false, err
		}
	}
	return s.repo.UpdateUser(req)
}

// Поиск пользователя по маске имени/фамилии
func (s *Service) SearchMask(mask string) ([]common_models.User, error) {
	users, err := s.repo.SearchMask(mask)
	if err != nil {
		return nil, err
	}

	res := make([]common_models.User, 0, len(users))
	for _, user := range users {
		res = append(res, user.User)
	}
	return res, nil
}

// Поиск пользователя по username
func (s *Service) SearchUsername(username string) (common_models.User, bool, error) {
	users, ok, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return common_models.User{}, false, err
	}
	if !ok {
		return common_models.User{}, false, nil
	}

	return users.User, true, nil
}

func (s *Service) CheckUser(username, password string) (common_models.UserID, bool) {
	user, ok, err := s.repo.GetUserByUsername(username)
	if err != nil || !ok {
		return 0, false
	}
	if !s.checkPasswordHash(password, user.Passhash) {
		return 0, false
	}
	return user.UserID, true
}

func (s *Service) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *Service) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
