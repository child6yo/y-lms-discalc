package repository

import "database/sql"

type Repository struct {
	Db *sql.DB
}

func NewRepository() (*Repository, error) {
	db, err := NewSqliteDb()
	if err != nil {
		return nil, err
	}
	return &Repository{Db: db}, nil
}