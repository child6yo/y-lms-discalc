package service

import (
	"strconv"

	"github.com/child6yo/y-lms-discalc/orchestrator"
)

func (s *MainService) CulculateExpression(userId int, expression string) (int, error) {
	cachedExp, exists := s.repo.GetCachedResult(expression)
	if !exists {
		expEntity := orchestrator.Expression{Expression: expression, Status: "Calculating..."}
		expId, err := s.repo.AddExpression(userId, &expEntity)
		if err != nil {
			return 0, err
		}
		expEntity.Id = strconv.Itoa(expId)
		*s.expChannel <- &expEntity

		return expId, nil
	}

	expId, err := s.repo.AddExpression(userId, cachedExp)
	if err != nil {
		return 0, err
	}
	return expId, nil
}

func (s *MainService) UpdateExpression(result *orchestrator.Expression) error {
	if result.Status != "ERROR" {
		s.repo.CacheResult(result)
	}

	return s.repo.UpdateExpression(result)
}

func (s *MainService) GetExpressioById(userId, expId int) (*orchestrator.Expression, error) {
	return s.repo.GetExpressionById(expId, userId)
}

func (s *MainService) GetExpressions(userId int) (*[]orchestrator.Expression, error) {
	return s.repo.GetExpressions(userId)
}