package worker

import (
	"context"
	"log"
	"time"

	"github.com/child6yo/y-lms-discalc/agent"
	"github.com/child6yo/y-lms-discalc/agent/pkg/service"
	pb "github.com/child6yo/y-lms-discalc/agent/proto"
)

func Worker(g int, grpcClient pb.OrchestratorServiceClient, evaluater service.PostfixEvaluater) {
	for {
		resp, err := grpcClient.GetTask(context.TODO(), nil)
		if err != nil {
			continue
		}

		task := agent.Task{Id: resp.Id,
			Arg1:          float64(resp.Arg1),
			Arg2:          float64(resp.Arg2),
			Operation:     resp.Operation,
			OperationTime: time.Duration(resp.OperationTime)}

		result := evaluater.PostfixEvaluate(task)

		if result.Error != "" {
			log.Println("Evaluation error:", result.Error)
			continue
		}

		_, err = grpcClient.TakeResult(context.TODO(), &pb.ResultResponse{Id: result.Id, Result: float32(result.Result), Error: result.Error})
		if err != nil {
			log.Println("Post result error:", err)
			continue
		}
		log.Printf("Task %s successfully done by worker %d", task.Id, g)
	}
}
