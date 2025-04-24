package repository

import (
	"database/sql"
)

type Repository struct {
	Db *sql.DB
	Cache *Cache
}

func NewRepository(cacheCap int) (*Repository, error) {
	db, err := newSqliteDb()
	if err != nil {
		return nil, err
	}
	cache := newCache(cacheCap)
	return &Repository{Db: db, Cache: cache}, nil
}