package repository

import (
	"database/sql"

	"github.com/child6yo/y-lms-discalc/orchestrator"
)

type Cache interface {
	Put(result *orchestrator.Expression)
	Get(expression string) (*orchestrator.Expression, bool)
}

type Database interface {
	CreateUser(user orchestrator.User) (int, error)
	GetUser(login, password string) (*orchestrator.User, error)

	AddExpression(userId int, expression *orchestrator.Expression) (int, error)
	UpdateExpression(expression *orchestrator.Expression) error
	GetExpressionById(expId, userId int) (*orchestrator.Expression, error)
	GetExpressions(userId int) (*[]orchestrator.Expression, error)
}

type Repository struct {
	Database
	Cache
}

type mainDatabase struct {
	db *sql.DB
}

func newMainDatabase(db *sql.DB) *mainDatabase {
	return &mainDatabase{db: db}
}

func NewRepository(cacheCap int) (*Repository, error) {
	db, err := newSqliteDb()
	if err != nil {
		return nil, err
	}
	cache := newExpressionCache(cacheCap)
	return &Repository{Database: newMainDatabase(db), Cache: cache}, nil
}