package db

import (
	"database/sql"
	"delivery/common/models"
	"fmt"
	"strings"
)

func (db *DB) AddItem(req models.CreateItemRequest) error {
	query := `INSERT INTO Items (data, owner_id) VALUES (?, ?)`
	_, err := db.db.Exec(query, req.Data, req.OwnerID)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении элемента: %v", err)
	}

	return nil
}

func (db *DB) GetItem(itemID models.ItemID) (models.Item, error) {
	var item models.Item
	query := "SELECT id, data, owner_id FROM Items WHERE id = ?"

	err := db.db.QueryRow(query, itemID).Scan(&item.ItemID, &item.Data, &item.OwnerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Item{}, fmt.Errorf("элемент с ID %d не найден", itemID)
		}
		return models.Item{}, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}

	return item, nil
}

func (db *DB) UpdateItem(req models.UpdateItemRequest) error {
	query := "UPDATE Items SET"
	args := []interface{}{}

	if req.OwnerID != nil {
		query += " owner_id = ?,"
		args = append(args, *req.OwnerID)
	}

	if req.Data != nil {
		query += " data = ?,"
		args = append(args, *req.Data)
	}

	if len(args) > 0 {
		query = strings.TrimSuffix(query, ",")
	}

	query += " WHERE id = ?"
	args = append(args, req.ItemID)

	_, err := db.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении элемента: %v", err)
	}
	return nil
}

func (db *DB) GetItemsByUserID(userID models.UserID) ([]models.Item, error) {
	query := `SELECT id, data, owner_id FROM Items WHERE owner_id = ?`
	rows, err := db.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		err := rows.Scan(&item.ItemID, &item.Data, &item.OwnerID)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строк: %v", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка после завершения итерации строк: %v", err)
	}
	return items, nil
}
