package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/child6yo/y-lms-discalc/agent"
	"github.com/child6yo/y-lms-discalc/agent/pkg/service"
)

func Worker(g int, url string) {
	for {
		resp, err := http.Get(url)
		if err != nil {
			log.Println("Get task error:", err)
			time.Sleep(5 * time.Second)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			time.Sleep(2 * time.Second)
			continue
		}

		var task agent.Task
		if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
			log.Println("Decode error:", err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

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

			resultJSON, err := json.Marshal(result)
			if err != nil {
				log.Println("Marshal error:", err)
				continue
			}

			postResp, err := http.Post(url, "application/json", bytes.NewBuffer(resultJSON))
			if err != nil {
				log.Println("Post result error:", err)
				continue
			}
			postResp.Body.Close()

			log.Printf("Worker %d Processed task %s", g, task.Id)
		case <-ctx.Done():
			log.Printf("Worker %d: Task %s exceeded time limit of %d seconds", g, task.Id, task.OperationTime)
			continue
		}
	}
}
