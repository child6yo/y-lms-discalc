package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/child6yo/y-lms-discalc/orchestrator"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/service/mock"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  string
		mockFunc     func(*mock.Service)
		wantStatus   int
		wantResponse string
	}{
		{
			name:        "successful creation",
			requestBody: `{"login":"test","password":"secret"}`,
			mockFunc: func(ms *mock.Service) {
				ms.CreateUserFunc = func(user orchestrator.User) (int, error) {
					return 123, nil
				}
			},
			wantStatus:   http.StatusCreated,
			wantResponse: `{"id":123}`,
		},
		{
			name:        "empty request body",
			requestBody: "",
			mockFunc: func(ms *mock.Service) {
				ms.CreateUserFunc = func(user orchestrator.User) (int, error) {
					return 0, nil
				}
			},
			wantStatus:   http.StatusInternalServerError,
			wantResponse: `{"error":"Internal server error"}`,
		},
		{
			name:        "invalid JSON",
			requestBody: `{"login":invalid}`,
			mockFunc: func(ms *mock.Service) {
				ms.CreateUserFunc = func(user orchestrator.User) (int, error) {
					return 0, nil
				}
			},
			wantStatus:   http.StatusUnprocessableEntity,
			wantResponse: `{"error":"Registration data is not valid"}`,
		},
		{
			name:        "service error",
			requestBody: `{"login":"test","password":"secret"}`,
			mockFunc: func(ms *mock.Service) {
				ms.CreateUserFunc = func(user orchestrator.User) (int, error) {
					return 123, errors.New("something went wrong")
				}
			},
			wantStatus:   http.StatusInternalServerError,
			wantResponse: `{"error":"Internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mock.Service{}
			tt.mockFunc(mockService)
			handler := &Handler{service: mockService}

			req := httptest.NewRequest("POST", "/register", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			handler.CreateUser(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("got status %d, want %d", resp.StatusCode, tt.wantStatus)
			}

			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			var got, want map[string]interface{}
			json.Unmarshal(body, &got)
			json.Unmarshal([]byte(tt.wantResponse), &want)

			if !reflect.DeepEqual(got, want) {
				t.Errorf("got response %v, want %v", got, want)
			}
		})
	}
}

func TestAuth(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  string
		mockFunc     func(*mock.Service)
		wantStatus   int
		wantResponse string
	}{
		{
			name:        "successfull case",
			requestBody: `{"login":"test","password":"secret"}`,
			mockFunc: func(ms *mock.Service) {
				ms.GenerateTokenFunc = func(username string, password string) (string, error) {
					return "token", nil
				}
			},
			wantStatus:   200,
			wantResponse: `{"jwt":"token"}`,
		},
		{
			name:        "invalid body case",
			requestBody: ``,
			mockFunc: func(ms *mock.Service) {
				ms.GenerateTokenFunc = func(username string, password string) (string, error) {
					return "", nil
				}
			},
			wantStatus:   500,
			wantResponse: `{"error":"Internal server error"}`,
		},
		{
			name:        "service error case",
			requestBody: `{"login":"test","password":"secret"}`,
			mockFunc: func(ms *mock.Service) {
				ms.GenerateTokenFunc = func(username string, password string) (string, error) {
					return "", errors.New("something went wrong")
				}
			},
			wantStatus:   422,
			wantResponse: `{"error":"Login data is not valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mock.Service{}
			tt.mockFunc(mockService)
			handler := &Handler{service: mockService}

			req := httptest.NewRequest("POST", "/login", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			handler.Auth(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("got status %d, want %d", resp.StatusCode, tt.wantStatus)
			}

			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			var got, want map[string]interface{}
			json.Unmarshal(body, &got)
			json.Unmarshal([]byte(tt.wantResponse), &want)

			if !reflect.DeepEqual(got, want) {
				t.Errorf("got response %v, want %v", got, want)
			}
		})
	}
}
