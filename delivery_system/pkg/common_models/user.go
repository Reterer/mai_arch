package common_models

type RegisterUserRequest struct {
	Usernmame string `json:"username" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

type UserID uint64
type User struct {
	UserID    UserID `json:"user_id"`
	Usernmame string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	PassHash  string `json:"-"`
}

type UpdateUserRequest struct {
	UserID    UserID  `json:"user_id,omitempty"`
	Usernmame *string `json:"username,omitempty"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Password  *string `json:"password,omitempty"`
}

type SerarchUsersResponce []User
