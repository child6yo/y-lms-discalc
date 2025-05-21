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

type userID string

var uID userID

// CreateUser - хендлер, отвечающий за регистрацию пользователей.
//
// В случае успеха возвращает полные регистрационные данные + айди пользователя.
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

	if user.Login == "" || user.Password == "" {
		httpNewError(w, 400, "Login and password are required", nil)
		return
	}

	user.ID, err = h.service.CreateUser(user)
	if err != nil {
		httpNewError(w, 500, "Internal server error", err)
		return
	}

	type output struct {
		ID int `json:"id"`
	}

	response := output{user.ID}
	responseData, err := json.MarshalIndent(response, "", " ")
	if err != nil {
		httpNewError(w, 500, "Internal server error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseData)
}

// Auth - хендлер, отвечающий за авторизацию пользователей
//
// В случае успеха возвращает JWT, валидный n кол-во часов,
// который может быть использован для авторизации.
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
		httpNewError(w, 400, "Login data is not valid", err)
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

// AuthorizeMiddleware - миддлвейр, отвечающий за авторизацию пользователей.
//
// Валидирует JWT пользователя, в случае успеха передает айди пользователя
// в контекст нижележащего уровня, используя тип userID.
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

		userID, err := h.service.ParseToken(headerParts[1])
		if err != nil {
			httpNewError(w, 401, "Authorization failed", nil)
			return
		}

		ctx := context.WithValue(r.Context(), uID, userID)
		next(w, r.WithContext(ctx))
	}
}

func getUserID(r *http.Request) (int, error) {
	userID := r.Context().Value(uID)
	if userID == nil {
		return 0, errors.New("user ID not found")
	}
	return userID.(int), nil
}
