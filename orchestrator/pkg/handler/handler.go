package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

var (
	expMap = make(map[int][]string)
	currentId = 1
) 

type ErrorModel struct {
	Error string `json:"error"`
}

func httpNewError(w http.ResponseWriter, statusCode int, message string) {
	slog.Error(message)

	response := ErrorModel{Error: message}
	responseData, _ := json.MarshalIndent(response, "", " ")

	http.Error(w, string(responseData), statusCode)
}