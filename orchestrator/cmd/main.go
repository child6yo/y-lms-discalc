package main

import (
	"log/slog"
	"net/http"
)

func main() {

	slog.Info("Server successfully started")
	http.ListenAndServe(":8000", nil)
}