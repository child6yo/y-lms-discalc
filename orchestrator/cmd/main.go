package main

import (
	"log"
	"net/http"
	"regexp"

	"github.com/child6yo/y-lms-discalc/orchestrator"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/handler"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/processor"
)


func main() {
	expressionInput := make(chan orchestrator.ExpAndId, 10)
	expressionsMap := make(chan map[int]orchestrator.Expression, 10)
	tasks := make(chan orchestrator.Task, 30)

	go processor.StartExpressionProcessor(expressionInput, tasks, expressionsMap)
	go handler.HandleExpressionsChanel(expressionsMap)

	http.HandleFunc("/api/v1/calculate", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            handler.GetExpressions(w, r)
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

	if err := http.ListenAndServe(":8000", nil); err != nil {
        log.Fatalf("Failed to start server")
    } else {
        log.Print("Server successfully started")
    }
}