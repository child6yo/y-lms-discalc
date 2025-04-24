package repository

import (
	"fmt"
	"strconv"

	"github.com/child6yo/y-lms-discalc/orchestrator"
)

func (r *Repository) AddExpression(userId int, expression *orchestrator.Result) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (user_id, exp, result, status) values ($1, $2, $3, $4)", expressionTable)

	res, err := r.Db.Exec(query, userId, expression.Expression, expression.Result, expression.Status)
	if err != nil {
		return 0, err
	}
	expId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(expId), nil
}

func (r *Repository) UpdateExpression(expression *orchestrator.Result) error {
	id, err := strconv.Atoi(expression.Id)
	if err != nil {
		return err
	}
	query := fmt.Sprintf("UPDATE %s SET result=$1, status=$2 WHERE id=%d", expressionTable, id)

	_, err = r.Db.Exec(query, expression.Result, expression.Status)
	if err != nil {
		return err
	}
	return nil
}
