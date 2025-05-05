package service

import (
	"strconv"

	"github.com/child6yo/y-lms-discalc/orchestrator"
)

// CulculateExpression запускает вычисление арифметического выражения
// и возвращает его айди.
//
// На вход принимает айди пользователя и выражение.
//
// В случае, если результат такого выражения есть в кеше, сразу же создает в базе данных
// выражение с результатом, хранящимся в кеше.
//
// В случае, если такого выражения в кеше нет - передает его в канал обработки,
// запуская процесс его вычисления.
func (s *MainService) CulculateExpression(userID int, expression string) (int, error) {
	cachedExp, exists := s.repo.Get(expression)
	if !exists {
		expEntity := orchestrator.Expression{Expression: expression, Status: "Calculating..."}
		expID, err := s.repo.AddExpression(userID, &expEntity)
		if err != nil {
			return 0, err
		}
		expEntity.ID = strconv.Itoa(expID)
		*s.expChannel <- &expEntity

		return expID, nil
	}

	expID, err := s.repo.AddExpression(userID, cachedExp)
	if err != nil {
		return 0, err
	}
	return expID, nil
}

// UpdateExpression используется обработчиком выражений для обновления результата вычисления
// в базе данных.
//
// В случае, если при вычислении не произошло ошибок, также записывает его в кеш.
func (s *MainService) UpdateExpression(result *orchestrator.Expression) error {
	if result.Status != "ERROR" {
		s.repo.Put(result)
	}

	return s.repo.UpdateExpression(result)
}

// GetExpressioByID возвращает выражение по его айди.
// На вход принимает айди пользователя и айди выражения.
func (s *MainService) GetExpressioByID(userID, expID int) (*orchestrator.Expression, error) {
	return s.repo.GetExpressionByID(expID, userID)
}

// GetExpressions возращает слайс всех выражений пользователя.
// На вход принимает айди пользователя.
func (s *MainService) GetExpressions(userID int) (*[]orchestrator.Expression, error) {
	return s.repo.GetExpressions(userID)
}
