package common_models

import "time"

type DeliveryID uint64

type CreateDeliveryRequest struct {
	FromUserID UserID   `json:"from_user_id" validate:"required"`
	ToUserID   UserID   `json:"to_user_id" validate:"required"`
	ToAddr     string   `json:"to_addr" validate:"required"`
	Items      []ItemID `json:"items" validate:"required"`
}

type Delivery struct {
	DeliveryID   DeliveryID `json:"delivery_id"`
	FromUserID   UserID     `json:"from_user_id"`
	FromAddr     string     `json:"from_addr"`
	ToUserID     UserID     `json:"to_user_id"`
	ToAddr       string     `json:"to_addr"`
	Status       int        `json:"status"`
	CreationDate time.Time  `json:"creation_date"`
}
