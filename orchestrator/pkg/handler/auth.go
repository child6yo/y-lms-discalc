package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/child6yo/y-lms-discalc/orchestrator"
)

type UserID string

var userID UserID

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

func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	var in input

	data, err := io.ReadAll(r.Body)
	if err != nil || len(data) == 0 {
		httpNewError(w, 500, "Internal server error", err)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(data, &in)
	if err != nil {
		httpNewError(w, 422, "Login data is not valid", err)
		return
	}

	token, err := h.service.GenerateToken(in.Login, in.Password)
	if err != nil {
		httpNewError(w, 422, "Login data is not valid", err)
		return
	}

	type output struct {
		JWT string `json:"jwt"`
	}

	out := output{JWT: token}

	responseData, err := json.MarshalIndent(out, "", " ")
	if err != nil {
		httpNewError(w, 500, "Internal server error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func (h *Handler) AuthorizeMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		if header == "" {
			httpNewError(w, 401, "Authorization failed", nil)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			httpNewError(w, 401, "Authorization failed", nil)
			return
		}

		if len(headerParts[1]) == 0 {
			httpNewError(w, 401, "Authorization failed", nil)
			return
		}

		userId, err := h.service.ParseToken(headerParts[1])
		if err != nil {
			httpNewError(w, 401, "Authorization failed", nil)
			return
		}

		ctx := context.WithValue(r.Context(), userID, userId)
		next(w, r.WithContext(ctx))
	}
}

func getUserId(r *http.Request) (int, error) {
	userId := r.Context().Value(userID)
	if userId == nil {
		return 0, errors.New("user ID not found")
	}
	return userId.(int), nil
}
