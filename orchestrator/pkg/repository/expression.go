package repository

import (
	"fmt"
	"strconv"

	"github.com/child6yo/y-lms-discalc/orchestrator"
)

// AddExpression добавляет арифметическое выражение в БД.
// На вход принимает айди пользователя и модель выражения.
func (d *mainDatabase) AddExpression(userID int, expression *orchestrator.Expression) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (user_id, exp, result, status) values ($1, $2, $3, $4)", expressionTable)

	res, err := d.db.Exec(query, userID, expression.Expression, expression.Result, expression.Status)
	if err != nil {
		return 0, err
	}
	expID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(expID), nil
}

// UpdateExpression обновляет результат и статус арифметического выражения.
// На вход принимает модель выражения.
func (d *mainDatabase) UpdateExpression(expression *orchestrator.Expression) error {
	id, err := strconv.Atoi(expression.ID)
	if err != nil {
		return err
	}
	query := fmt.Sprintf("UPDATE %s SET result=$1, status=$2 WHERE id=%d", expressionTable, id)

	_, err = d.db.Exec(query, expression.Result, expression.Status)
	if err != nil {
		return err
	}
	return nil
}

// GetExpressionById возвращает арифметическое выражение по его айди.
// На вход принимает айди выражения и айди пользователя.
func (d *mainDatabase) GetExpressionByID(expID, userID int) (*orchestrator.Expression, error) {
	var result orchestrator.Expression

	query := fmt.Sprintf("SELECT id, result, exp, status FROM %s WHERE user_id=$1 AND id=$2", expressionTable)

	row := d.db.QueryRow(query, userID, expID)
	err := row.Scan(&result.ID, &result.Result, &result.Expression, &result.Status)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetExpressions возвращает слайс арифметических выражений, принадлежащих пользователю.
// На вход принимает айди пользователя.
func (d *mainDatabase) GetExpressions(userID int) (*[]orchestrator.Expression, error) {
	var result []orchestrator.Expression

	query := fmt.Sprintf("SELECT id, result, exp, status FROM %s WHERE user_id=$1", expressionTable)
	rows, err := d.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		r := orchestrator.Expression{}
		err := rows.Scan(&r.ID, &r.Result, &r.Expression, &r.Status)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}

	return &result, nil
}
