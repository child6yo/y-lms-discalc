package service

import (
	"strconv"

	"github.com/child6yo/y-lms-discalc/orchestrator"
)

func (s *Service) CulculateExpression(userId int, expression string) (int, error) {
	cachedExp, exists := s.repo.GetCachedResult(expression)
	if !exists {
		expEntity := orchestrator.Result{Expression: expression, Status: "Calculating..."}
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

func (s *Service) UpdateExpression(result *orchestrator.Result) error {
	if result.Status != "ERROR" {
		s.repo.CacheResult(result)
	}

	return s.repo.UpdateExpression(result)
}

func (s *Service) GetExpressioById() {

}

func (s *Service) GetExpressions() {

}