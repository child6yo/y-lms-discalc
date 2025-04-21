package worker

import (
	"context"
	"log"
	"time"

	"github.com/child6yo/y-lms-discalc/agent"
	"github.com/child6yo/y-lms-discalc/agent/pkg/service"
	pb "github.com/child6yo/y-lms-discalc/agent/proto"
)

func Worker(g int, url string, grpcClient pb.OrchestratorServiceClient) {
	for {
		resp, err := grpcClient.GetTask(context.TODO(), nil)
		if err != nil {
			log.Println("error: ", err)
			continue
		}

		task := agent.Task{Id: resp.Id,
			Arg1:          float64(resp.Arg1),
			Arg2:          float64(resp.Arg2),
			Operation:     resp.Opetation,
			OperationTime: time.Duration(resp.OperationTime)}

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(task.OperationTime)*time.Second)

		resultCh := make(chan agent.Result)
		go func() {
			defer cancel()
			result := service.EvaluatePostfix(task)
			resultCh <- result
		}()

		select {
		case result := <-resultCh:
			if result.Error != "" {
				log.Println("Evaluation error:", result.Error)
				continue
			}

			_, err := grpcClient.TakeResult(context.TODO(), &pb.ResultResponse{Id: result.Id, Result: float32(result.Result), Error: result.Error})
			if err != nil {
				log.Println("Post result error:", err)
				continue
			}
		case <-ctx.Done():
			log.Printf("Worker %d: Task %s exceeded time limit of %d seconds", g, task.Id, task.OperationTime)
			continue
		}
	}
}
