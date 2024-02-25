package models

type CreateItemRequest struct {
	OwnerID UserID `json:"owner_id"`
	Data    string `json:"data"`
}

type ItemID uint64
type Item struct {
	ItemID  ItemID `json:"item_id"`
	Data    string `json:"data"`
	OwnerID UserID `json:"owner_id"`
}

type UpdateItemRequest struct {
	ItemID  ItemID  `json:"-"`
	OwnerID *UserID `json:"owner_id"`
	Data    *string `json:"data"`
}
