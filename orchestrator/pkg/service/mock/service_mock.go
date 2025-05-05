package mock

import "github.com/child6yo/y-lms-discalc/orchestrator"

// Service реализует мок сервиса.
type Service struct {
	PostfixExpressionFunc   func(expression string) ([]string, error)
	CulculateExpressionFunc func(userID int, expression string) (int, error)
	UpdateExpressionFunc    func(result *orchestrator.Expression) error
	GetExpressioByIDFunc    func(userID, expId int) (*orchestrator.Expression, error)
	GetExpressionsFunc      func(userID int) (*[]orchestrator.Expression, error)
	CreateUserFunc          func(user orchestrator.User) (int, error)
	GenerateTokenFunc       func(username, password string) (string, error)
	ParseTokenFunc          func(accessToken string) (int, error)
}

// PostfixExpression mock
func (m *Service) PostfixExpression(expression string) ([]string, error) {
	return m.PostfixExpressionFunc(expression)
}

// CulculateExpression mock
func (m *Service) CulculateExpression(userID int, expression string) (int, error) {
	return m.CulculateExpressionFunc(userID, expression)
}

// UpdateExpression mock
func (m *Service) UpdateExpression(result *orchestrator.Expression) error {
	return m.UpdateExpressionFunc(result)
}

// GetExpressioByID mock
func (m *Service) GetExpressioByID(userID, expID int) (*orchestrator.Expression, error) {
	return m.GetExpressioByIDFunc(userID, expID)
}

// GetExpressions mock
func (m *Service) GetExpressions(userID int) (*[]orchestrator.Expression, error) {
	return m.GetExpressionsFunc(userID)
}

// CreateUser mock
func (m *Service) CreateUser(user orchestrator.User) (int, error) {
	return m.CreateUserFunc(user)
}

// GenerateToken mock
func (m *Service) GenerateToken(username, password string) (string, error) {
	return m.GenerateTokenFunc(username, password)
}

// ParseToken mock
func (m *Service) ParseToken(accessToken string) (int, error) {
	return m.ParseTokenFunc(accessToken)
}
