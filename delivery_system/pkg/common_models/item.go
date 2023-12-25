package common_models

type ItemID uint64

type CreateItemRequest struct {
	OwnerID UserID `json:"owner_id" validate:"required"`
	Data    string `json:"data" validate:"required"`
}

type Item struct {
	ItemID  ItemID `json:"item_id"`
	Data    string `json:"data"`
	OwnerID UserID `json:"owner_id"`
}

type PatchItemRequest struct {
	OwnerID UserID `json:"owner_id" validate:"required"`
	Data    string `json:"data" validate:"required"`
}
