package mock

import "github.com/child6yo/y-lms-discalc/orchestrator"


type MockService struct {
	PostfixExpressionFunc    func(expression string) ([]string, error)
	CulculateExpressionFunc  func(userId int, expression string) (int, error)
	UpdateExpressionFunc     func(result *orchestrator.Expression) error
	GetExpressioByIdFunc     func(userId, expId int) (*orchestrator.Expression, error)
	GetExpressionsFunc       func(userId int) (*[]orchestrator.Expression, error)
	CreateUserFunc           func(user orchestrator.User) (int, error)
	GenerateTokenFunc        func(username, password string) (string, error)
	ParseTokenFunc           func(accessToken string) (int, error)
}

func (m *MockService) PostfixExpression(expression string) ([]string, error) {
	return m.PostfixExpressionFunc(expression)
}

func (m *MockService) CulculateExpression(userId int, expression string) (int, error) {
	return m.CulculateExpressionFunc(userId, expression)
}

func (m *MockService) UpdateExpression(result *orchestrator.Expression) error {
	return m.UpdateExpressionFunc(result)
}

func (m *MockService) GetExpressioById(userId, expId int) (*orchestrator.Expression, error) {
	return m.GetExpressioByIdFunc(userId, expId)
}

func (m *MockService) GetExpressions(userId int) (*[]orchestrator.Expression, error) {
	return m.GetExpressionsFunc(userId)
}

func (m *MockService) CreateUser(user orchestrator.User) (int, error) {
	return m.CreateUserFunc(user)
}

func (m *MockService) GenerateToken(username, password string) (string, error) {
	return m.GenerateTokenFunc(username, password)
}

func (m *MockService) ParseToken(accessToken string) (int, error) {
	return m.ParseTokenFunc(accessToken)
}