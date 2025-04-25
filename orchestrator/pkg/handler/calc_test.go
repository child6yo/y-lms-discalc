package handler

import (
	"context"
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

func TestCulculateExpression(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  string
		setupContext func(r *http.Request) *http.Request
		mockFunc     func(*mock.MockService)
		wantStatus   int
		wantResponse string
	}{
		{
			name:        "successfull case",
			requestBody: `{"expression":"2+2*2"}`,
			setupContext: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), userID, 123)
				return r.WithContext(ctx)
			},
			mockFunc: func(ms *mock.MockService) {
				ms.CulculateExpressionFunc = func(userId int, expr string) (int, error) {
					return 1, nil
				}
			},
			wantStatus:   201,
			wantResponse: `{"id":1}`,
		},
		{
			name:        "unauthorize error",
			requestBody: `{"expression":"2+2*2"}`,
			setupContext: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), userID, nil)
				return r.WithContext(ctx)
			},
			mockFunc: func(ms *mock.MockService) {
				ms.CulculateExpressionFunc = func(userId int, expr string) (int, error) {
					return 1, nil
				}
			},
			wantStatus:   401,
			wantResponse: `{"error":"JWT is not valid"}`,
		},
		{
			name:        "invalid data error",
			requestBody: `{"e":}`,
			setupContext: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), userID, 123)
				return r.WithContext(ctx)
			},
			mockFunc: func(ms *mock.MockService) {
				ms.CulculateExpressionFunc = func(userId int, expr string) (int, error) {
					return 1, nil
				}
			},
			wantStatus:   422,
			wantResponse: `{"error":"Expression is not valid"}`,
		},
		{
			name:        "invalid input error",
			requestBody: `{"expression":"2+"}`,
			setupContext: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), userID, 123)
				return r.WithContext(ctx)
			},
			mockFunc: func(ms *mock.MockService) {
				ms.CulculateExpressionFunc = func(userId int, expr string) (int, error) {
					return 0, errors.New("something went wrong")
				}
			},
			wantStatus:   500,
			wantResponse: `{"error":"Internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mock.MockService{}
			tt.mockFunc(mockService)
			handler := &Handler{service: mockService}

			req := httptest.NewRequest("POST", "/calculate", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			req = tt.setupContext(req)

			w := httptest.NewRecorder()
			handler.CulculateExpression(w, req)

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

func TestGetExpressions(t *testing.T) {
	tests := []struct {
		name         string
		setupContext func(r *http.Request) *http.Request
		mockFunc     func(*mock.MockService)
		wantStatus   int
		wantResponse string
	}{
		{
			name: "successfull case",
			setupContext: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), userID, 123)
				return r.WithContext(ctx)
			},
			mockFunc: func(ms *mock.MockService) {
				ms.GetExpressionsFunc = func(userId int) (*[]orchestrator.Expression, error) {
					return &[]orchestrator.Expression{
						{Id: "1", Result: 4, Expression: "2+2", Status: "Success"},
						{Id: "2", Result: 0, Expression: "2+", Status: "ERROR"}}, nil
				}
			},
			wantStatus: 200,
			wantResponse: `{"expressions": [
								{
									"id": "1",
									"result": 4,
									"expression": "2+2",
									"error": "Success"
								},
								{
									"id": "2",
									"result": 0,
									"expression": "2+",
									"error": "ERROR"
								}
							]}`,
		},
		{
			name: "unauthorize error",
			setupContext: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), userID, nil)
				return r.WithContext(ctx)
			},
			mockFunc: func(ms *mock.MockService) {
				ms.GetExpressionsFunc = func(userId int) (*[]orchestrator.Expression, error) {
					return &[]orchestrator.Expression{}, nil
				}
			},
			wantStatus:   401,
			wantResponse: `{"error": "JWT is not valid"}`,
		},
		{
			name: "service error",
			setupContext: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), userID, 123)
				return r.WithContext(ctx)
			},
			mockFunc: func(ms *mock.MockService) {
				ms.GetExpressionsFunc = func(userId int) (*[]orchestrator.Expression, error) {
					return nil, errors.New("something went wrong")
				}
			},
			wantStatus:   500,
			wantResponse: `{"error": "Internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mock.MockService{}
			tt.mockFunc(mockService)
			handler := &Handler{service: mockService}

			req := httptest.NewRequest("GET", "/calculate", nil)
			req.Header.Set("Content-Type", "application/json")
			req = tt.setupContext(req)

			w := httptest.NewRecorder()
			handler.GetExpressions(w, req)

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

func TestGetExpressionsById(t *testing.T) {
	tests := []struct {
		name         string
		setupContext func(r *http.Request) *http.Request
		mockFunc     func(*mock.MockService)
		path         string
		wantStatus   int
		wantResponse string
	}{
		{
			name: "successfull case",
			setupContext: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), userID, 123)
				return r.WithContext(ctx)
			},
			mockFunc: func(ms *mock.MockService) {
				ms.GetExpressioByIdFunc = func(userId int, expId int) (*orchestrator.Expression, error) {
					return &orchestrator.Expression{
						Id:         "1",
						Result:     4,
						Expression: "2+2",
						Status:     "Success"}, nil
				}
			},
			path:       "/api/v1/expressions/1",
			wantStatus: 200,
			wantResponse: `{
							"id": "1",
							"result": 4,
							"expression": "2+2",
							"error": "Success"
							}`,
		},
		{
			name: "unauthorize error",
			setupContext: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), userID, nil)
				return r.WithContext(ctx)
			},
			mockFunc: func(ms *mock.MockService) {
				ms.GetExpressioByIdFunc = func(userId int, expId int) (*orchestrator.Expression, error) {
					return &orchestrator.Expression{}, nil
				}
			},
			path:       "/api/v1/expressions/1",
			wantStatus: 401,
			wantResponse: `{"error":"JWT is not valid"}`,
		},
		{
			name: "unknown path error",
			setupContext: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), userID, 123)
				return r.WithContext(ctx)
			},
			mockFunc: func(ms *mock.MockService) {
				ms.GetExpressioByIdFunc = func(userId int, expId int) (*orchestrator.Expression, error) {
					return &orchestrator.Expression{}, nil
				}
			},
			path:       "/api/v1/expressions/one",
			wantStatus: 500,
			wantResponse: `{"error":"Internal server error"}`,
		},
		{
			name: "service error",
			setupContext: func(r *http.Request) *http.Request {
				ctx := context.WithValue(r.Context(), userID, 123)
				return r.WithContext(ctx)
			},
			mockFunc: func(ms *mock.MockService) {
				ms.GetExpressioByIdFunc = func(userId int, expId int) (*orchestrator.Expression, error) {
					return nil, errors.New("something went wrong")
				}
			},
			path:       "/api/v1/expressions/1",
			wantStatus: 500,
			wantResponse: `{"error":"Internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &mock.MockService{}
			tt.mockFunc(mockService)
			handler := &Handler{service: mockService}

			req := httptest.NewRequest("POST", tt.path, nil)
			req.Header.Set("Content-Type", "application/json")
			req = tt.setupContext(req)

			w := httptest.NewRecorder()
			handler.GetExpressionById(w, req)

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
