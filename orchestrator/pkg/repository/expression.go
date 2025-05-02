package repository

import (
	"fmt"
	"strconv"

	"github.com/child6yo/y-lms-discalc/orchestrator"
)

func (d *mainDatabase) AddExpression(userId int, expression *orchestrator.Expression) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (user_id, exp, result, status) values ($1, $2, $3, $4)", expressionTable)

	res, err := d.db.Exec(query, userId, expression.Expression, expression.Result, expression.Status)
	if err != nil {
		return 0, err
	}
	expId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(expId), nil
}

func (d *mainDatabase) UpdateExpression(expression *orchestrator.Expression) error {
	id, err := strconv.Atoi(expression.Id)
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

func (d *mainDatabase) GetExpressionById(expId, userId int) (*orchestrator.Expression, error) {
	var result orchestrator.Expression

	query := fmt.Sprintf("SELECT id, result, exp, status FROM %s WHERE user_id=$1 AND id=$2", expressionTable)

	row := d.db.QueryRow(query, userId, expId)
	err := row.Scan(&result.Id, &result.Result, &result.Expression, &result.Status)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *mainDatabase) GetExpressions(userId int) (*[]orchestrator.Expression, error) {
	var result []orchestrator.Expression

	query := fmt.Sprintf("SELECT id, result, exp, status FROM %s WHERE user_id=$1", expressionTable)
	rows, err := d.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		r := orchestrator.Expression{}
		err := rows.Scan(&r.Id, &r.Result, &r.Expression, &r.Status)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}

	return &result, nil
}