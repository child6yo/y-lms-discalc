package processor

import (
	"math"
	"sync"
	"testing"
	"time"

	"github.com/child6yo/y-lms-discalc/orchestrator"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/service/mock"
)

func TestExpressionProcessorSuccess(t *testing.T) {
	expression := orchestrator.Expression{ID: "1", Expression: "2+2*2"}
	expect := orchestrator.Expression{Result: 6, Status: "Success"}

	taskc := make(chan *orchestrator.Task, 3)
	debugChan := make(chan orchestrator.Expression)
	var wg sync.WaitGroup

	config := map[string]time.Duration{
		"+": 100 * time.Millisecond,
		"-": 100 * time.Millisecond,
		"*": 100 * time.Millisecond,
		"/": 100 * time.Millisecond,
	}

	mockService := &mock.Service{}
	mockService.PostfixExpressionFunc = func(expression string) ([]string, error) {
		return []string{"2", "2", "2", "*", "+"}, nil
	}
	mockService.UpdateExpressionFunc = func(result *orchestrator.Expression) error {
		return nil
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		processExpression(expression, &taskc, config, mockService, debugChan)
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

			res := orchestrator.Expression{ID: task.ID, Result: r, Status: err}
			chInterface, _ := TaskResultChannels.Load(task.ID)
			resultChan, _ := chInterface.(chan orchestrator.Expression)
			resultChan <- res
		}
	}()

	select {
	case answer := <-debugChan:
		res := answer
		if res.Status != expect.Status {
			t.Fatalf("Test FAILED. Expected status: %s, Status: %s", expect.Status, res.Status)
		} else if math.Abs(res.Result-expect.Result) > 1e-9 {
			t.Fatalf("Test FAILED. Expected result: %.2f, Result: %.2f", expect.Result, res.Result)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Test FAILED. Timeout occurred")
	}

	wg.Wait()
}

func TestExpressionProcessorFail(t *testing.T) {
	expression := orchestrator.Expression{ID: "1", Expression: "2+2*"}
	expect := orchestrator.Expression{Result: 0, Status: "ERROR"}

	taskc := make(chan *orchestrator.Task, 3)
	debugChan := make(chan orchestrator.Expression)
	var wg sync.WaitGroup

	config := map[string]time.Duration{
		"+": 100 * time.Millisecond,
		"-": 100 * time.Millisecond,
		"*": 100 * time.Millisecond,
		"/": 100 * time.Millisecond,
	}

	mockService := &mock.Service{}
	mockService.PostfixExpressionFunc = func(expression string) ([]string, error) {
		return []string{"2", "2", "2", "+"}, nil
	}
	mockService.UpdateExpressionFunc = func(result *orchestrator.Expression) error {
		return nil
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		processExpression(expression, &taskc, config, mockService, debugChan)
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

			res := orchestrator.Expression{ID: task.ID, Result: r, Status: err}
			chInterface, _ := TaskResultChannels.Load(task.ID)
			resultChan, _ := chInterface.(chan orchestrator.Expression)
			resultChan <- res
		}
	}()

	select {
	case answer := <-debugChan:
		res := answer
		if res.Status != expect.Status {
			t.Fatalf("Test FAILED. Expected status: %s, Status: %s", expect.Status, res.Status)
		} else if math.Abs(res.Result-expect.Result) > 1e-9 {
			t.Fatalf("Test FAILED. Expected result: %.2f, Result: %.2f", expect.Result, res.Result)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Test FAILED. Timeout occurred")
	}

	wg.Wait()
}
