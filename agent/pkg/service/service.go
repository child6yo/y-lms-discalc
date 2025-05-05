package service

import (
	"github.com/child6yo/y-lms-discalc/agent"
)

// PostfixEvaluater определяет интерфейс вычислителя арифметичексих выражений.
type PostfixEvaluater interface {
	// PostfixEvaluate вычисляет простое арифметическое выражение (напр. 2 + 2).
	//
	// На вход принимает задачу с выражением в формате agent.Task, на выход результат в формате agent.Result.
	PostfixEvaluate(task agent.Task) agent.Result
}

// EvaluateService реализует вычислитель арифметических выражений.
type EvaluateService struct{}

// PostfixEvaluate - метод, реализующий вычисление арифметических выражений.
// 
// На вход принимает задачу с выражением в формате agent.Task, на выход результат в формате agent.Result.
func (p *EvaluateService) PostfixEvaluate(task agent.Task) agent.Result {
	var res float64

	switch task.Operation {
	case "+":
		res = task.Arg1 + task.Arg2
	case "-":
		res = task.Arg1 - task.Arg2
	case "*":
		res = task.Arg1 * task.Arg2
	case "/":
		if task.Arg2 == 0 {
			return agent.Result{ID: task.ID, Result: 0, Error: "division by zero"}
		}
		res = task.Arg1 / task.Arg2
	default:
		return agent.Result{ID: task.ID, Result: 0, Error: "unknown operator"}
	}

	return agent.Result{ID: task.ID, Result: res, Error: ""}
}
