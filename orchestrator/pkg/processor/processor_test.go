package processor

import (
	"math"
	"sync"
	"testing"
	"time"

	"github.com/child6yo/y-lms-discalc/orchestrator"
)

func TestExpressionProcessor(t *testing.T) {
	testCases := []struct {
		name       string
		expression orchestrator.ExpAndId
		expect     orchestrator.Expression
	}{
		{
			name:       "OK",
			expression: orchestrator.ExpAndId{Id: 1, Expression: []string{"2", "2", "2", "*", "+"}},
			expect:     orchestrator.Expression{Id: 1, Status: "Success", Result: 6},
		},
		{
			name:       "Invalid postfix",
			expression: orchestrator.ExpAndId{Id: 1, Expression: []string{"2", "2", "2"}},
			expect:     orchestrator.Expression{Id: 1, Status: "ERROR", Result: 0},
		},
		{
			name:       "Invalid operation",
			expression: orchestrator.ExpAndId{Id: 1, Expression: []string{"2", "+"}},
			expect:     orchestrator.Expression{Id: 1, Status: "ERROR", Result: 0},
		},
		{
			name:       "Computing error",
			expression: orchestrator.ExpAndId{Id: 1, Expression: []string{"2", "0", "/"}},
			expect:     orchestrator.Expression{Id: 1, Status: "ERROR", Result: 0},
		},
	}

	taskc := make(chan orchestrator.Task, 3)
	op := make(chan map[int]orchestrator.Expression, 3)
	var wg sync.WaitGroup

	t.Parallel()
	for _, test := range testCases {
		wg.Add(1)
		go func() {
			defer wg.Done()
			processExpression(test.expression, taskc, op)
		}()

		go func() {
			for task := range taskc {
				var r float64
				var err string
				switch task.Operation {
				case "+":
					r = task.Arg1 + task.Arg2
				case "-":
					r = task.Arg1 - task.Arg2
				case "*":
					r = task.Arg1 * task.Arg2
				case "/":
					if task.Arg2 == 0 {
						err = "Division by zero"
					} else {
						r = task.Arg1 / task.Arg2
					}
				}

				res := orchestrator.Result{Id: task.Id, Result: r, Error: err}
				chInterface, _ := TaskResultChannels.Load(task.Id)
				resultChan, _ := chInterface.(chan orchestrator.Result)
				resultChan <- res
			}
		}()

		select {
		case answer := <-op:
			res := answer[test.expression.Id]
			if res.Status != test.expect.Status {
				t.Fatalf("Test '%s' FAILED. Expected status: %s, Status: %s", test.name, test.expect.Status, res.Status)
			} else if math.Abs(res.Result-test.expect.Result) > 1e-9 { 
				t.Fatalf("Test '%s' FAILED. Expected result: %.2f, Result: %.2f", test.name, test.expect.Result, res.Result)
			}
		case <-time.After(5 * time.Second):
			t.Fatalf("Test '%s' FAILED. Timeout occurred", test.name)
		}
	}

	wg.Wait()
}
