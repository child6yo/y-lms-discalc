package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"

	"github.com/child6yo/y-lms-discalc/orchestrator"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func httpNewError(w http.ResponseWriter, statusCode int, message string, err error) {
	if err != nil {
		log.Println("Handling error: ", err)
	} else {
		log.Println("Unknown error: ", message)
	}

	response := orchestrator.ErrorModel{Error: message}
	responseData, _ := json.MarshalIndent(response, "", " ")

	http.Error(w, string(responseData), statusCode)
}

func addCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func (h *Handler) StaticFileHandler(w http.ResponseWriter, r *http.Request) {
	addCORSHeaders(w)
	absPath, err := filepath.Abs("./client/index.html")
	if err != nil {
		http.Error(w, "Error resolving file path", http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, absPath)
}