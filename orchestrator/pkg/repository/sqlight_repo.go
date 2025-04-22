package repository

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const (
	userTable = "user"
)

type Repository struct {
	Db *sql.DB
}

func NewRepository() (*Repository, error) {
	db, err := NewSqlightDb()
	if err != nil {
		return nil, err
	}
	return &Repository{Db: db}, nil
}

func NewSqlightDb() (*sql.DB, error) {
	ctx := context.TODO()

	db, err := sql.Open("sqlite3", "database/discalc.db")
	if err != nil {
		return nil, err
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}