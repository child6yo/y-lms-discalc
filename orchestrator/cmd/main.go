package main

import (
	"log/slog"
	"net/http"

	"github.com/child6yo/y-lms-discalc/orchestrator"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/handler"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/processor"
)

func main() {
	expressionInput := make(chan orchestrator.ExpAndId)
	expressionsMap := make(chan map[int]orchestrator.Expression)
	tasks := make(chan orchestrator.Task)

	go processor.StartExpressionProcessor(expressionInput, tasks, expressionsMap)

	http.HandleFunc("/calculate", handler.CulculateExpression(expressionInput))

	slog.Info("Server successfully started")
	http.ListenAndServe(":8000", nil)
}