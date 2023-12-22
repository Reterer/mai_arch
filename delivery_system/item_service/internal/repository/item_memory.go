package repository

import (
	"delivery_system/item_service/config"
	"delivery_system/pkg/common_models"
	"slices"
	"sync"
)

type ItemMemory struct {
	cfg     *config.Repository
	items   []common_models.Item
	nextIdx common_models.ItemID
	mu      sync.RWMutex
}

func NewItemMemory(cfg *config.Repository, initItems []common_models.Item) *ItemMemory {
	var nextIdx common_models.ItemID
	for i := range initItems {
		if nextIdx < initItems[i].ItemID {
			nextIdx = initItems[i].ItemID
		}
	}
	nextIdx++

	return &ItemMemory{
		cfg:     cfg,
		items:   slices.Clone(initItems),
		nextIdx: nextIdx,
		mu:      sync.RWMutex{},
	}
}

func (r *ItemMemory) findIdxByItemID(itemID common_models.ItemID) (int, error) {
	for i := range r.items {
		if r.items[i].ItemID == itemID {
			return i, nil
		}
	}
	return 0, NotExistsErr
}

func (r *ItemMemory) CreateItem(newItem common_models.CreateItemRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.items = append(r.items, common_models.Item{
		ItemID:  r.nextIdx,
		Data:    newItem.Data,
		OwnerID: newItem.OwnerID,
	})
	r.nextIdx++

	return nil
}

func (r *ItemMemory) GetItem(itemID common_models.ItemID) (common_models.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	idx, err := r.findIdxByItemID(itemID)
	if err != nil {
		return common_models.Item{}, err
	}

	return r.items[idx], nil
}

func (r *ItemMemory) UpdateItem(item common_models.Item) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	idx, err := r.findIdxByItemID(item.ItemID)
	if err != nil {
		return err
	}

	r.items[idx] = item
	return nil
}

func (r *ItemMemory) GetItems(userID common_models.UserID) ([]common_models.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var items []common_models.Item
	for i := range r.items {
		if r.items[i].OwnerID == userID {
			items = append(items, r.items[i])
		}
	}

	return items, nil
}
