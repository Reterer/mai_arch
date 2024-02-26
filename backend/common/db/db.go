package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	db *sql.DB
}

func InitDB(addr string) (*DB, error) {
	db, err := sql.Open("mysql", addr)
	if err != nil {
		return nil, err
	}

	return &DB{
		db: db,
	}, nil
}

func (db *DB) Close() error {
	return db.db.Close()
}
