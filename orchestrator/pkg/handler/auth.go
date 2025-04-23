package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/child6yo/y-lms-discalc/orchestrator"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user orchestrator.User

	data, err := io.ReadAll(r.Body)
	if err != nil || len(data) == 0 {
		httpNewError(w, 500, "Internal server error", err)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(data, &user)
	if err != nil {
		httpNewError(w, 422, "Registration data is not valid", err)
		return
	}

	user.Id, err = h.service.CreateUser(user)
	if err != nil {
		httpNewError(w, 500, "Internal server error", err)
		return
	}

	response := orchestrator.User{Id: user.Id, Login: user.Login, Password: user.Password}
	responseData, err := json.MarshalIndent(response, "", " ")
	if err != nil {
		httpNewError(w, 500, "Internal server error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseData)
}