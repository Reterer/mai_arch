package db

import (
	"database/sql"
	"delivery/common/models"
	"fmt"
	"strings"
)

func (db *DB) AddUser(req models.RegisterUserRequest) error {
	query := `INSERT INTO Users (username, first_name, last_name, pass_hash, email) VALUES (?, ?, ?, ?, ?)`
	_, err := db.db.Exec(query, req.Username, req.FirstName, req.LastName, req.Password, req.Email)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении пользователя: %v", err)
	}
	return nil
}

func (db *DB) GetUser(uid models.UserID) (models.UserWithPass, error) {
	query := `SELECT id, username, first_name, last_name, pass_hash, email FROM Users WHERE id = ?`
	row := db.db.QueryRow(query, uid)

	var user models.UserWithPass
	err := row.Scan(&user.UserID, &user.Username, &user.FirstName, &user.LastName, &user.Passhash, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("пользователь с UserID %d не найден", uid)
		}
		return user, fmt.Errorf("ошибка при получении пользователя: %v", err)
	}
	return user, nil
}

func (db *DB) UpdateUser(req models.UpdateUserRequest) error {
	query := "UPDATE Users SET"
	values := []interface{}{}
	if req.Username != nil {
		query += " username = ?,"
		values = append(values, *req.Username)
	}
	if req.FirstName != nil {
		query += " first_name = ?,"
		values = append(values, *req.FirstName)
	}
	if req.LastName != nil {
		query += " last_name = ?,"
		values = append(values, *req.LastName)
	}
	if req.Password != nil {
		query += " pass_hash = ?,"
		values = append(values, *req.Password)
	}
	if req.Email != nil {
		query += " email = ?,"
		values = append(values, *req.Email)
	}

	if len(values) != 0 {
		query = strings.TrimSuffix(query, ",")
	}

	query += " WHERE id = ?"
	values = append(values, req.UserID)

	_, err := db.db.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении пользователя: %v", err)
	}
	return nil
}

func (db *DB) SearchUser(firstName, lastName, username string) ([]models.User, error) {
	var users []models.User

	query := `SELECT id, username, first_name, last_name, email FROM Users WHERE`
	args := []interface{}{}

	if firstName != "" {
		query += ` first_name LIKE ? AND`
		args = append(args, firstName)
	}

	if lastName != "" {
		query += ` last_name LIKE ? AND`
		args = append(args, lastName)
	}

	if username != "" {
		query += ` username LIKE ? AND`
		args = append(args, username)
	}

	if len(args) == 0 {
		query = strings.TrimSuffix(query, " WHERE")
	} else {
		query = strings.TrimSuffix(query, " AND")
	}

	fmt.Println(query)

	rows, err := db.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.UserID, &user.Username, &user.FirstName, &user.LastName, &user.Email)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строк: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка после завершения перебора строк: %v", err)
	}

	return users, nil
}
