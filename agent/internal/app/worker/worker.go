package worker

import (
	"context"
	"log"
	"time"

	"github.com/child6yo/y-lms-discalc/agent"
	"github.com/child6yo/y-lms-discalc/agent/internal/app/service"
	pb "github.com/child6yo/y-lms-discalc/agent/proto"
)

// Worker - воркер, выполняющий обработку арифметических выражений.
//
// Запрашивает задачи о оркестратора через grpcClient.GetTask
// и передает результат их вычисления через grpcClient.TakeResult.
func Worker(g int, grpcClient pb.OrchestratorServiceClient, evaluater service.PostfixEvaluater) {
	for {
		resp, err := grpcClient.GetTask(context.TODO(), nil)
		if err != nil {
			continue
		}

		task := agent.Task{ID: resp.Id,
			Arg1:          float64(resp.Arg1),
			Arg2:          float64(resp.Arg2),
			Operation:     resp.Operation,
			OperationTime: time.Duration(resp.OperationTime)}

		result := evaluater.PostfixEvaluate(task)

		if result.Error != "" {
			log.Println("Evaluation error:", result.Error)
			continue
		}

		_, err = grpcClient.TakeResult(context.TODO(), &pb.ResultResponse{Id: result.ID, Result: float32(result.Result), Error: result.Error})
		if err != nil {
			log.Println("Post result error:", err)
			continue
		}
		log.Printf("Task %s successfully done by worker %d", task.ID, g)
	}
}
