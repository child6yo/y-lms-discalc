package repository

import (
	"fmt"

	"github.com/child6yo/y-lms-discalc/orchestrator"
)

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

func (d *mainDatabase) GetUser(login, password string) (*orchestrator.User, error) {
	var user orchestrator.User

	query := fmt.Sprintf("SELECT * FROM %s WHERE login=$1 AND password=$2", userTable)
	row := d.db.QueryRow(query, login, password)

	err := row.Scan(&user.Id, &user.Login, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
