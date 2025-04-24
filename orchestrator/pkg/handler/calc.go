package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/child6yo/y-lms-discalc/orchestrator"
)

func (h *Handler) CulculateExpression() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

func (h *Handler) GetExpressions(w http.ResponseWriter, r *http.Request) {
	// mu.Lock()
	// defer mu.Unlock()

	// expressions := make([]orchestrator.Expression, 0, len(exps))

	// for _, value := range exps {
	// 	expressions = append(expressions, value)
	// }

	// response := orchestrator.ExpressionList{Expressions: expressions}
	// responseData, err := json.MarshalIndent(response, "", " ")
	// if err != nil {
	// 	httpNewError(w, 500, "Internal server error", err)
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte{})
}

func (h *Handler) GetExpressionById(w http.ResponseWriter, r *http.Request) {
	// path := r.URL.Path

	// re := regexp.MustCompile(`/api/v1/expressions/(\d+)`)
	// matches := re.FindStringSubmatch(path)

	// if len(matches) < 2 {
	// 	httpNewError(w, 500, "Internal server error", nil)
	// 	return
	// }
	// id, err := strconv.Atoi(matches[1])
	// if err != nil {
	// 	httpNewError(w, 500, "Internal server error", err)
	// 	return
	// }

	// if _, exists := exps[id]; !exists {
	// 	httpNewError(w, 404, "Invalid expression id", nil)
	// 	return
	// }

	// responseData, err := json.MarshalIndent(exps[id], "", " ")
	// if err != nil {
	// 	httpNewError(w, 500, "Internal server error", err)
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte{})
}
