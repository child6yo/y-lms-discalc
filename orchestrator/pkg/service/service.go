package service

import (
	"github.com/child6yo/y-lms-discalc/orchestrator"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/repository"
)

type Service interface {
	PostfixExpression(expression string) ([]string, error)

	CulculateExpression(userId int, expression string) (int, error)
	UpdateExpression(result *orchestrator.Expression) error
	GetExpressioById(userId, expId int) (*orchestrator.Expression, error)
	GetExpressions(userId int) (*[]orchestrator.Expression, error)

	CreateUser(user orchestrator.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(accessToken string) (int, error)
}

type MainService struct {
	repo       *repository.Repository
	expChannel *chan *orchestrator.Expression
}

func NewService(repository *repository.Repository, expChannel *chan *orchestrator.Expression) *MainService {
	return &MainService{repo: repository, expChannel: expChannel}
}
