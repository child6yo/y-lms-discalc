package worker

import (
	"bytes"
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

		result := service.EvaluatePostfix(task)
		if result.Error != "" {
			log.Println("Evaluation error:", err)
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
	}
}