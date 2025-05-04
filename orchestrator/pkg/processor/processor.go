package processor

import (
	"log"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/child6yo/y-lms-discalc/orchestrator"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/service"
)

var (
	TaskResultChannels sync.Map
	globalTaskCounter  uint64
)

func processExpression(exp orchestrator.Expression, taskChan *chan *orchestrator.Task,
	config map[string]time.Duration, service service.Service, debug chan orchestrator.Expression) {
	var stack []float64
	taskCounter := 0

	expr, err := service.PostfixExpression(exp.Expression)
	if err != nil {
		log.Println("Something went wrong in processor: ", err, expr)
		debug <- orchestrator.Expression{}
		return
	}

	for _, token := range expr {
		if value, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, value)
		} else {
			if len(stack) < 2 {
				expession := orchestrator.Expression{Id: exp.Id, Result: 0, Status: "ERROR"}
				err := service.UpdateExpression(&expession)
				if err != nil {
					log.Println("Something went wrong in processor 1: ", err)
				}
				log.Printf("Expression %s: insufficient operands for operator %s\n", exp.Id, token)
				debug <- expession
				return
			}

			operandB := stack[len(stack)-1]
			operandA := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			resultChan := make(chan orchestrator.Expression, 1)

			localTaskCounter := atomic.AddUint64(&globalTaskCounter, 1)
			key := strconv.FormatUint(localTaskCounter, 10)
			TaskResultChannels.Store(key, resultChan)

			task := orchestrator.Task{
				Id:            key,
				Arg1:          operandA,
				Arg2:          operandB,
				Operation:     token,
				OperationTime: config[token],
			}

			*taskChan <- &task

			select {
			case res := <-resultChan:
				TaskResultChannels.Delete(taskCounter)
				if res.Status != "" {
					expession := orchestrator.Expression{Id: exp.Id, Result: 0, Status: "ERROR"}
					err := service.UpdateExpression(&expession)
					if err != nil {
						log.Println("Something went wrong in processor 2: ", err)
					}
					log.Printf("Expression %s, task %s error: %v\n", exp.Id, task.Id, res.Status)
					debug <- expession
					return
				}

				stack = append(stack, res.Result)
			case <-time.After(task.OperationTime + 1*time.Second):
				TaskResultChannels.Delete(taskCounter)
				expession := orchestrator.Expression{Id: exp.Id, Result: 0, Status: "ERROR"}
				err := service.UpdateExpression(&expession)
				if err != nil {
					log.Println("Something went wrong in processor 3: ", err)
				}
				log.Printf("Expression %s, task %s timeout\n", exp.Id, task.Id)
				debug <- expession
				return
			}
		}
	}

	if len(stack) != 1 {
		expession := orchestrator.Expression{Id: exp.Id, Result: 0, Status: "ERROR"}
		err := service.UpdateExpression(&expession)
		if err != nil {
			log.Println("Something went wrong in processor 4: ", err)
		}
		log.Printf("Expression %s: invalid RPN, stack: %v\n", exp.Id, stack)
		debug <- expession
		return
	}
	finalResult := stack[0]
	expession := orchestrator.Expression{Id: exp.Id, Expression: exp.Expression, Status: "Success", Result: finalResult}
	err = service.UpdateExpression(&expession)
	if err != nil {
		log.Println("Something went wrong in processor 5: ", err)
		debug <- expession
	}
	log.Printf("Expression %s computed successfully, result: %v\n", exp.Id, finalResult)
	debug <- expession
}

func StartExpressionProcessor(input *chan *orchestrator.Expression, taskChan *chan *orchestrator.Task,
	config map[string]time.Duration, service service.Service) {
	for exp := range *input {
		go processExpression(*exp, taskChan, config, service, nil)
	}
}
