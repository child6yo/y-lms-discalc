package processor

import (
	"log"
	"strconv"
	"sync"
	"time"
	"sync/atomic"

	"github.com/child6yo/y-lms-discalc/orchestrator"
)

var (
	TaskResultChannels sync.Map
 	globalTaskCounter uint64
)

func processExpression(exp orchestrator.ExpAndId, taskChan chan orchestrator.Task, output chan map[int]orchestrator.Expression) {
	var stack []float64
	taskCounter := 0
	m := make(map[int]orchestrator.Expression)

	for _, token := range exp.Expression {
		if value, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, value)
		} else {
			if len(stack) < 2 {
				m[exp.Id] = orchestrator.Expression{Id: exp.Id, Status: "ERROR", Result: 0}
				output <- m
				log.Printf("Expression %d: insufficient operands for operator %s\n", exp.Id, token)
				return
			}

			operandB := stack[len(stack)-1]
			operandA := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			resultChan := make(chan orchestrator.Result, 1)

			localTaskCounter := atomic.AddUint64(&globalTaskCounter, 1)
			key := strconv.FormatUint(localTaskCounter, 10)
			TaskResultChannels.Store(key, resultChan)
			
			task := orchestrator.Task{
				Id:            key, 
				Arg1:          operandA,
				Arg2:          operandB,
				Operation:     token,
				OperationTime: 5 * time.Second,
			}

			taskChan <- task

			select {
			case res := <-resultChan:
				TaskResultChannels.Delete(taskCounter)
				if res.Error != "" {
					m[exp.Id] = orchestrator.Expression{Id: exp.Id, Status: "ERROR", Result: 0}
					output <- m
					log.Printf("Expression %d, task %s error: %v\n", exp.Id, task.Id, res.Error)
					return
				}

				stack = append(stack, res.Result)
			case <-time.After(task.OperationTime + 5*time.Second):
				TaskResultChannels.Delete(taskCounter)
				m[exp.Id] = orchestrator.Expression{Id: exp.Id, Status: "ERROR", Result: 0}
				output <- m
				log.Printf("Expression %d, task %s timeout\n", exp.Id, task.Id)
				return
			}
		}
	}

	if len(stack) != 1 {
		m[exp.Id] = orchestrator.Expression{Id: exp.Id, Status: "ERROR", Result: 0}
		output <- m
		log.Printf("Expression %d: invalid RPN, stack: %v\n", exp.Id, stack)
		return
	}
	finalResult := stack[0]
	m[exp.Id] = orchestrator.Expression{Id: exp.Id, Status: "Success", Result: finalResult}
	output <- m
	log.Printf("Expression %d computed successfully, result: %v\n", exp.Id, finalResult)
}

func StartExpressionProcessor(input chan orchestrator.ExpAndId, taskChan chan orchestrator.Task, output chan map[int]orchestrator.Expression) {
	for exp := range input {
		go processExpression(exp, taskChan, output)
	}
}
