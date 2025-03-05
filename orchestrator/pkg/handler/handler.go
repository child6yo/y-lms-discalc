package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/child6yo/y-lms-discalc/orchestrator"
)

var (
	currentId = 1
	exps      = make(map[int]orchestrator.Expression)
	mu        sync.Mutex
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

func addCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func StaticFileHandler(w http.ResponseWriter, r *http.Request) {
	addCORSHeaders(w)
	absPath, err := filepath.Abs("./client/index.html")
	if err != nil {
		http.Error(w, "Error resolving file path", http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, absPath)
}