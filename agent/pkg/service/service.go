package service

import (
	"github.com/child6yo/y-lms-discalc/agent"
)

type Evaluator interface {
	EvaluatePostfix(task agent.Task) agent.Result
}

type PostfixEvaluator struct{}

func NewPostfixEvaluator() *PostfixEvaluator {
	return &PostfixEvaluator{}
}

func (p *PostfixEvaluator) EvaluatePostfix(task agent.Task) agent.Result {
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
			return agent.Result{Id: task.Id, Result: 0, Error: "division by zero"}
		}
		res = task.Arg1 / task.Arg2
	default:
		return agent.Result{Id: task.Id, Result: 0, Error: "unknown operator"}
	}

	return agent.Result{Id: task.Id, Result: res, Error: ""}
}
