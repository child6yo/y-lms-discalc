package main

import (
	"log/slog"
	"net/http"

	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/handler"
)

func main() {
	expressionInput := make(chan map[int][]string)

	// Запуск воркеров
	// worker.StartWorkers(input, output, 4)

	// Запуск HTTP-сервера
	http.HandleFunc("/calculate", handler.CulculateExpression(expressionInput))

	slog.Info("Server successfully started")
	http.ListenAndServe(":8000", nil)
}