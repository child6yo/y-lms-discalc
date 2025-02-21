package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/child6yo/y-lms-discalc/agent"
	"github.com/child6yo/y-lms-discalc/agent/pkg/service"
)

var url = "http://localhost:8000/internal/task"

func worker(g int) {
	for {
		resp, err := http.Get(url)
		if err != nil {
			log.Println("[AGENT] Get task error:", err)
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
			log.Println("[AGENT] Decode error:", err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		result := service.EvaluatePostfix(task)
		if result.Error != "" {
			log.Println("[AGENT] Evaluation error:", err)
			continue
		}
		
		resultJSON, err := json.Marshal(result)
		if err != nil {
			log.Println("[AGENT] Marshal error:", err)
			continue
		}

		postResp, err := http.Post(url, "application/json", bytes.NewBuffer(resultJSON))
		if err != nil {
			log.Println("[AGENT] Post result error:", err)
			continue
		}
		postResp.Body.Close()

		log.Printf("[AGENT] Worker %d Processed task %s", g, task.Id)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	for w := 1; w <= 5; w++ {
		go worker(w)
	}

	wg.Wait()
}