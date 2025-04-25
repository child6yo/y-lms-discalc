package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/child6yo/y-lms-discalc/orchestrator"
)

func (h *Handler) CulculateExpression(w http.ResponseWriter, r *http.Request) {
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

	userId, err := getUserId(r)
	if err != nil {
		httpNewError(w, 401, "JWT is not valid", err)
		return
	}

	expId, err := h.service.CulculateExpression(userId, req.Expression)
	if err != nil {
		httpNewError(w, 500, "Internal server error", err)
		return
	}

	response := orchestrator.ExpressionId{Id: expId}
	responseData, err := json.MarshalIndent(response, "", " ")
	if err != nil {
		httpNewError(w, 500, "Internal server error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseData)
}

func (h *Handler) GetExpressions(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserId(r)
	if err != nil {
		httpNewError(w, 401, "JWT is not valid", err)
		return
	}

	result, err := h.service.GetExpressions(userId)
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

func (h *Handler) GetExpressionById(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	re := regexp.MustCompile(`/api/v1/expressions/(\d+)`)
	matches := re.FindStringSubmatch(path)

	if len(matches) < 2 {
		httpNewError(w, 500, "Internal server error", nil)
		return
	}
	expId, err := strconv.Atoi(matches[1])
	if err != nil {
		httpNewError(w, 500, "Internal server error", err)
		return
	}

	userId, err := getUserId(r)
	if err != nil {
		httpNewError(w, 401, "JWT is not valid", err)
		return
	}

	result, err := h.service.GetExpressioById(userId, expId)
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
