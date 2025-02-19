package handler

import (
	"encoding/json"
	"io"
	"net/http"

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

		exps[currentId] = orchestrator.Expression{Id: currentId, Status: "Calculating...", Result: 0}
		res := orchestrator.ExpAndId{Id: currentId, Expression: expression}
		input <- res
		currentId++

		responce := orchestrator.ExpressionId{Id: currentId}
		responseData, err := json.MarshalIndent(responce, "", " ")
		if err != nil { 
			httpNewError(w, 500, "Internal server error")
			return 
		}

		w.Header().Set("Content-Type", "application/json") 
		w.Write(responseData)
	}
}

func GetExpressions() {

}

func GetExpressionById() {
	
}