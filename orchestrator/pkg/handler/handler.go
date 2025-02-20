package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/child6yo/y-lms-discalc/orchestrator"
)

var (
	currentId = 1
	exps = make(map[int]orchestrator.Expression)
	mu sync.Mutex
)

func HandleExpressionsChanel(c chan map[int]orchestrator.Expression) {
	for expmap := range c {
		for id, exp := range expmap {
			mu.Lock()
			exps[id] = exp
			mu.Unlock()
		}
	}
}

func httpNewError(w http.ResponseWriter, statusCode int, message string) {
	log.Println("Handling error: ", message)

	response := orchestrator.ErrorModel{Error: message}
	responseData, _ := json.MarshalIndent(response, "", " ")

	http.Error(w, string(responseData), statusCode)
}