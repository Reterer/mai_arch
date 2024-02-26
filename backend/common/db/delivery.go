package db

import (
	"database/sql"
	"delivery/common/models"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (db *DB) AddDelivery(req models.CreateDeliveryRequest) error {
	tx, err := db.db.Begin()
	if err != nil {
		return fmt.Errorf("ошибка при начале транзакции: %v", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			fmt.Println(err)
		}
	}()

	result, err := tx.Exec("INSERT INTO Deliveries (from_user_id, from_addr, to_user_id, to_addr, status, creation_date) VALUES (?, ?, ?, ?, ?, ?)",
		req.FromUserID, "", req.ToUserID, req.ToAddr, 0, time.Now())
	if err != nil {
		return fmt.Errorf("ошибка при добавлении доставки: %v", err)
	}

	deliveryID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("ошибка при получении ID доставки: %v", err)
	}

	for _, itemID := range req.Items {
		_, err := tx.Exec("INSERT INTO ItemsDeliveries (item_id, delivery_id) VALUES (?, ?)",
			itemID, deliveryID)
		if err != nil {
			return fmt.Errorf("ошибка при добавлении предмета доставки: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("ошибка при коммите транзакции: %v", err)
	}

	return nil
}

func (db *DB) GetDeliveriesByFromUserID(uid models.UserID) ([]models.Delivery, error) {
	query := `
		SELECT d.id, d.from_user_id, d.from_addr, d.to_user_id, d.to_addr, d.status, d.creation_date, GROUP_CONCAT(id.item_id) AS item_ids
		FROM Deliveries d
		LEFT JOIN ItemsDeliveries id ON d.id = id.delivery_id
		WHERE d.from_user_id = ?
		GROUP BY d.id;
	`
	return db.getDeliveriesHelper(query, uid)
}

func (db *DB) GetDeliveriesByToUserID(uid models.UserID) ([]models.Delivery, error) {
	query := `
		SELECT d.id, d.from_user_id, d.from_addr, d.to_user_id, d.to_addr, d.status, d.creation_date, GROUP_CONCAT(id.item_id) AS item_ids
		FROM Deliveries d
		LEFT JOIN ItemsDeliveries id ON d.id = id.delivery_id
		WHERE d.to_user_id = ?
		GROUP BY d.id;
	`
	return db.getDeliveriesHelper(query, uid)
}

func (db *DB) getDeliveriesHelper(query string, uid models.UserID) ([]models.Delivery, error) {
	rows, err := db.db.Query(query, uid)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer rows.Close()

	var deliveries []models.Delivery
	for rows.Next() {
		var delivery models.Delivery
		var itemIDs sql.NullString // Используем sql.NullString для обработки NULL значений
		var creationDateStr string
		err := rows.Scan(&delivery.DeliveryID, &delivery.FromUserID, &delivery.FromAddr, &delivery.ToUserID, &delivery.ToAddr, &delivery.Status, &creationDateStr, &itemIDs)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки: %v", err)
		}

		creationDate, err := time.Parse("2006-01-02 15:04:05", creationDateStr)
		if err != nil {
			return nil, fmt.Errorf("ошибка при преобразовании времени: %v", err)
		}
		delivery.CreationDate = creationDate

		if itemIDs.Valid {
			var items []models.ItemID
			for _, idStr := range strings.Split(itemIDs.String, ",") {
				id, err := strconv.ParseUint(idStr, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("ошибка при преобразовании ID предмета: %v", err)
				}
				items = append(items, models.ItemID(id))
			}
			delivery.Items = items
		}

		deliveries = append(deliveries, delivery)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка после прохода по строкам: %v", err)
	}

	return deliveries, nil
}
