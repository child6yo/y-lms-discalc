package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/child6yo/y-lms-discalc/orchestrator"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/service"
)

func CulculateExpression(input chan orchestrator.ExpAndId) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req orchestrator.ExpressionInput

		data, err := io.ReadAll(r.Body)
		if err != nil || len(data) == 0 {
			httpNewError(w, 500, "Internal server error")
			return
		}
		defer r.Body.Close()

		err = json.Unmarshal(data, &req)
		if err != nil {
			httpNewError(w, 422, "Expression is not valid")
			return
		}

		expression, err := service.PostfixExpression(req.Expression)
		if err != nil {
			httpNewError(w, 422, "Expression is not valid")
			return
		}
		mu.Lock()
		exps[currentId] = orchestrator.Expression{Id: currentId, Status: "Calculating...", Result: 0}
		mu.Unlock()

		res := orchestrator.ExpAndId{Id: currentId, Expression: expression}
		input <- res

		response := orchestrator.ExpressionId{Id: currentId}
		currentId++
		responseData, err := json.MarshalIndent(response, "", " ")
		if err != nil {
			httpNewError(w, 500, "Internal server error")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(responseData)
	}
}

func GetExpressions(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
    defer mu.Unlock()
	
	expressions := make([]orchestrator.Expression, 0, len(exps))

	for _, value := range exps {
		expressions = append(expressions, value)
	}

	response := orchestrator.ExpressionList{Expressions: expressions}
	responseData, err := json.MarshalIndent(response, "", " ")
	if err != nil {
		httpNewError(w, 500, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseData)
}

func GetExpressionById(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	re := regexp.MustCompile(`/api/v1/expressions/(\d+)`)
	matches := re.FindStringSubmatch(path)

	if len(matches) < 2 {
		httpNewError(w, 500, "Internal server error")
		return
	}
	id, err := strconv.Atoi(matches[1])
	if err != nil {
		httpNewError(w, 500, "Internal server error")
		return
	}

	if _, exists := exps[id]; !exists {
		httpNewError(w, 404, "Invalid expression id")
		return
	}

	responseData, err := json.MarshalIndent(exps[id], "", " ")
	if err != nil {
		httpNewError(w, 500, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseData)
}
