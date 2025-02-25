package main

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/child6yo/y-lms-discalc/agent/pkg/worker"
)

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Print("Failed to load env")
		return defaultValue
	}
	return value
}

func getIntEnv(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Print("Failed to load env: ", err)
		return defaultValue
	}
	return value
}

func main() {
	var url = "http://localhost:8000/internal/task"

	computingPower := getIntEnv("COMPUTING_POWER", 10)

	var wg sync.WaitGroup

	for w := 1; w <= computingPower; w++ {
		go func() {
			defer wg.Done()
			worker.Worker(w, url)
		}()
		wg.Add(1)
	}

	log.Print("Agent successfully started")

	wg.Wait()
}