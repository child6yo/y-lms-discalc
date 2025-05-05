package repository

import (
	"fmt"

	"github.com/child6yo/y-lms-discalc/orchestrator"
)

// CreateUser создает пользователя в БД.
// Принимает на вход модель пользователя,
// в случае успеха создает пользователя в базе данных.
func (d *mainDatabase) CreateUser(user orchestrator.User) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (login, password) values ($1, $2)", userTable)

	result, err := d.db.Exec(query, user.Login, user.Password)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// GetUser возвращает пользователя из БД.
// Принимает на вход логин и пароль,
// в случае успеха возвращает модель пользователя.
func (d *mainDatabase) GetUser(login, password string) (*orchestrator.User, error) {
	var user orchestrator.User

	query := fmt.Sprintf("SELECT * FROM %s WHERE login=$1 AND password=$2", userTable)
	row := d.db.QueryRow(query, login, password)

	err := row.Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
