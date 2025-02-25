package main

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/child6yo/y-lms-discalc/orchestrator"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/handler"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/processor"
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
	config := map[string]time.Duration{}
	config["+"] = time.Duration(getIntEnv("TIME_ADDITION_MS", 100) * int(time.Millisecond))
	config["-"] = time.Duration(getIntEnv("TIME_SUBTRACTION_MS", 100) * int(time.Millisecond))
	config["*"] = time.Duration(getIntEnv("TIME_MULTIPLICATIONS_MS", 100) * int(time.Millisecond))
	config["/"] = time.Duration(getIntEnv("TIME_DIVISIONS_MS", 100) * int(time.Millisecond))

	expressionInput := make(chan orchestrator.ExpAndId, 10)
	expressionsMap := make(chan map[int]orchestrator.Expression, 10)
	tasks := make(chan orchestrator.Task, 30)

	go processor.StartExpressionProcessor(expressionInput, tasks, expressionsMap, config)
	go handler.HandleExpressionsChanel(expressionsMap)

	http.HandleFunc("/api/v1/calculate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CulculateExpression(expressionInput)(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
	http.HandleFunc("/api/v1/expressions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.GetExpressions(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
	http.HandleFunc("/api/v1/expressions/", func(w http.ResponseWriter, r *http.Request) {
		pattern := `/api/v1/expressions/\d+`
		matched, err := regexp.MatchString(pattern, r.URL.Path)
		if err != nil || !matched {
			http.NotFound(w, r)
			return
		}
		if r.Method == http.MethodGet {
			handler.GetExpressionById(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	http.HandleFunc("/internal/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetTask(tasks)(w, r)
		case http.MethodPost:
			handler.Result()(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server successfully started")
	http.ListenAndServe(":8000", nil)
}
