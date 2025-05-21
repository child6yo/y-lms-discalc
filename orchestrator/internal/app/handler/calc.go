package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/child6yo/y-lms-discalc/orchestrator"
)

// CalculateExpression - хендлер, отвечающий за прием арифметических выражений.
//
// В случае успеха направляет выражение на обработку;
// пользователю возвращает айди выражения, находящегося в обработке.
func (h *Handler) CalculateExpression(w http.ResponseWriter, r *http.Request) {
	var req orchestrator.ExpressionInput

	data, err := io.ReadAll(r.Body)
	if err != nil || len(data) == 0 {
		httpNewError(w, 500, "Internal server error", err)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(data, &req)
	if err != nil {
		httpNewError(w, 422, "Expression is not valid", err)
		return
	}

	userID, err := getUserID(r)
	if err != nil {
		httpNewError(w, 401, "JWT is not valid", err)
		return
	}

	expID, err := h.service.CulculateExpression(userID, req.Expression)
	if err != nil {
		httpNewError(w, 500, "Internal server error", err)
		return
	}

	response := orchestrator.ExpressionID{ID: expID}
	responseData, err := json.MarshalIndent(response, "", " ")
	if err != nil {
		httpNewError(w, 500, "Internal server error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseData)
}

// GetExpressions - хендлер, отвечающий за выдачу арифметических выражений.
//
// В случае успеха возвращает все выражения, созданные пользователем с их результатами.
func (h *Handler) GetExpressions(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		httpNewError(w, 401, "JWT is not valid", err)
		return
	}

	result, err := h.service.GetExpressions(userID)
	if err != nil {
		httpNewError(w, 500, "Internal server error", err)
		return
	}

	response := orchestrator.ExpressionListOutput{Expressions: *result}
	responseData, err := json.MarshalIndent(response, "", " ")
	if err != nil {
		httpNewError(w, 500, "Internal server error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseData)
}

// GetExpressionByID - хендлер, отвечающий за выдачу арифметического выражения по айди.
//
// Возвращает пользователю выражение по айди в том случае если оно существует и было создано этим пользователем.
func (h *Handler) GetExpressionByID(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	re := regexp.MustCompile(`/api/v1/expressions/(\d+)`)
	matches := re.FindStringSubmatch(path)

	if len(matches) < 2 {
		httpNewError(w, 500, "Internal server error", nil)
		return
	}
	expID, err := strconv.Atoi(matches[1])
	if err != nil {
		httpNewError(w, 500, "Internal server error", err)
		return
	}

	userID, err := getUserID(r)
	if err != nil {
		httpNewError(w, 401, "JWT is not valid", err)
		return
	}

	result, err := h.service.GetExpressioByID(userID, expID)
	if err != nil {
		httpNewError(w, 500, "Internal server error", err)
		return
	}

	responseData, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		httpNewError(w, 500, "Internal server error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseData)
}
