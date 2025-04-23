package repository

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const (
	userTable = "user"
	resultTable = "result"
)

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
		password TEXT
	);`
	if _, err := db.Exec(user); err != nil {
		return err
	}

	result := `
	CREATE TABLE IF NOT EXISTS result(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		user_id INTEGER NOT NULL,
		expression TEXT NOT NULL,
		equal REAL NOT NULL,
	
		FOREIGN KEY (user_id) REFERENCES user (id)
	);`
	if _, err := db.Exec(result); err != nil {
		return err
	}

	return nil
}