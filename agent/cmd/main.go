package main

import (
	"sync"

	"github.com/child6yo/y-lms-discalc/agent/pkg/worker"
)


func main() {
	var url = "http://localhost:8000/internal/task"

	var wg sync.WaitGroup
	wg.Add(1)

	for w := 1; w <= 5; w++ {
		go worker.Worker(w, url)
	}

	wg.Wait()
}