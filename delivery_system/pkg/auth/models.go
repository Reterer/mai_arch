package auth

import "delivery_system/pkg/common_models"

type CheckAuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CheckAuthResponce struct {
	Status bool                 `json:"status"`
	UserID common_models.UserID `json:"user_id"`
}
