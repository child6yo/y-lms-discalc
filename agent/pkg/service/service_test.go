package service

import (
	"testing"

	"github.com/child6yo/y-lms-discalc/agent"
)

func TestEvaluatePostfix(t *testing.T) {
	p := EvaluateService{}

	tests := []struct {
		name           string
		task           agent.Task
		expectedResult agent.Result
	}{
		{
			name: "Addition",
			task: agent.Task{ID: "1", Operation: "+", Arg1: 5, Arg2: 3},
			expectedResult: agent.Result{
				ID:     "1",
				Result: 8,
				Error:  "",
			},
		},
		{
			name: "Subtraction",
			task: agent.Task{ID: "2", Operation: "-", Arg1: 5, Arg2: 3},
			expectedResult: agent.Result{
				ID:     "2",
				Result: 2,
				Error:  "",
			},
		},
		{
			name: "Multiplication",
			task: agent.Task{ID: "3", Operation: "*", Arg1: 5, Arg2: 3},
			expectedResult: agent.Result{
				ID:     "3",
				Result: 15,
				Error:  "",
			},
		},
		{
			name: "Division",
			task: agent.Task{ID: "4", Operation: "/", Arg1: 6, Arg2: 3},
			expectedResult: agent.Result{
				ID:     "4",
				Result: 2,
				Error:  "",
			},
		},
		{
			name: "Division by Zero",
			task: agent.Task{ID: "5", Operation: "/", Arg1: 5, Arg2: 0},
			expectedResult: agent.Result{
				ID:     "5",
				Result: 0,
				Error:  "division by zero",
			},
		},
		{
			name: "Unknown Operator",
			task: agent.Task{ID: "6", Operation: "^", Arg1: 5, Arg2: 3},
			expectedResult: agent.Result{
				ID:     "6",
				Result: 0,
				Error:  "unknown operator",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := p.PostfixEvaluate(tt.task)

			if result != tt.expectedResult {
				t.Errorf("For task %+v, expected result %+v, got %+v", tt.task, tt.expectedResult, result)
			}
		})
	}
}
