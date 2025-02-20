package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/child6yo/y-lms-discalc/orchestrator"
)

var (
	currentId = 1
	exps = make(map[int]orchestrator.Expression)
)

func HandleExpressionsChanel(c chan map[int]orchestrator.Expression) {
	for exp := range c {
		exps[currentId] = exp[currentId]
	}
}

func httpNewError(w http.ResponseWriter, statusCode int, message string) {
	slog.Error(message)

	response := orchestrator.ErrorModel{Error: message}
	responseData, _ := json.MarshalIndent(response, "", " ")

	http.Error(w, string(responseData), statusCode)
}