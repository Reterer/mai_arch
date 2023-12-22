package common_models

type ItemID uint64

type CreateItemRequest struct {
	Data    string `json:"data"`
	OwnerID UserID `json:"-"`
}

type Item struct {
	ItemID  ItemID `json:"item_id"`
	Data    string `json:"data"`
	OwnerID UserID `json:"owner_id"`
}

type PatchItemRequest struct {
	Data string `json:"data"`
}
