package repository

import (
	"context"
	"database/sql"
)

const (
	userTable       = "user"
	expressionTable = "expression"
)

// NewSqliteDb создает новое подключение к локальному хранилищу sqlite.
func NewSqliteDb() (*sql.DB, error) {
	ctx := context.TODO()

	db, err := sql.Open("sqlite3", "database/discalc.db")
	if err != nil {
		return nil, err
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	err = initDatabase(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func initDatabase(db *sql.DB) error {
	user := `
	CREATE TABLE IF NOT EXISTS user(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		login TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`
	if _, err := db.Exec(user); err != nil {
		return err
	}

	result := `
	CREATE TABLE IF NOT EXISTS expression(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		user_id INTEGER NOT NULL,
		exp TEXT NOT NULL,
		result REAL,
		status TEXT NOT NULL,
	
		FOREIGN KEY (user_id) REFERENCES user (id)
	);`
	if _, err := db.Exec(result); err != nil {
		return err
	}

	return nil
}
