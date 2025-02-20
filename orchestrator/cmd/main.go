package main

import (
	"log/slog"
	"net/http"

	"github.com/child6yo/y-lms-discalc/orchestrator"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/handler"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/processor"
)

func main() {
	expressionInput := make(chan orchestrator.ExpAndId, 10)
	expressionsMap := make(chan map[int]orchestrator.Expression, 10)
	tasks := make(chan orchestrator.Task, 30)
	resultsChan := make(chan orchestrator.Result, 10)

	go processor.StartExpressionProcessor(expressionInput, tasks, expressionsMap)
	go handler.HandleExpressionsChanel(expressionsMap)

	http.HandleFunc("/api/v1/calculate", handler.CulculateExpression(expressionInput))

	http.HandleFunc("/internal/task", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            handler.GetTask(tasks)
        case http.MethodPost:
            handler.Result(resultsChan)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })

	slog.Info("Server successfully started")
	http.ListenAndServe(":8000", nil)
}