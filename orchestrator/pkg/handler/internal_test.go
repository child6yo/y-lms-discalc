package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/child6yo/y-lms-discalc/orchestrator"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/processor"
)

func TestGetTask(t *testing.T) {
	testCases := []struct {
		name         string
		prepare      func(chan orchestrator.Task)
		expectBody   orchestrator.Task
		expectStatus int
	}{
		{
			name: "200, Task received",
			prepare: func(output chan orchestrator.Task) {
				go func() {
					output <- orchestrator.Task{Id: "1", Arg1: 2, Arg2: 2, Operation: "+", OperationTime: 1 * time.Second}
				}()
			},
			expectBody:   orchestrator.Task{Id: "1", Arg1: 2, Arg2: 2, Operation: "+", OperationTime: 1 * time.Second},
			expectStatus: 200,
		},
		{
			name:         "404, No task received",
			prepare:      func(output chan orchestrator.Task) {},
			expectBody:   orchestrator.Task{},
			expectStatus: 404,
		},
	}

	for _, test := range testCases {
		output := make(chan orchestrator.Task, 1)
		test.prepare(output)
		handler := GetTask(output)

		req := httptest.NewRequest("GET", "http://localhost:8000/api/v1/task", nil)
		w := httptest.NewRecorder()

		handler(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Test %s: error reading response body: %v", test.name, err)
		}

		var result orchestrator.Task
		if test.expectStatus == 200 {
			err = json.Unmarshal(body, &result)
			if err != nil {
				t.Fatalf("Test %s: error unmarshalling response body: %v", test.name, err)
			}
		}

		if resp.StatusCode != test.expectStatus {
			t.Errorf("test %s failed: result: %d, expected: %d", test.name, resp.StatusCode, test.expectStatus)
		} else if test.expectStatus == 200 && !reflect.DeepEqual(result, test.expectBody) {
			t.Errorf("test %s failed: result: %+v, expected: %+v", test.name, result, test.expectBody)
		} else {
			t.Logf("Test %s success", test.name)
		}
	}
}

func TestResultHandler(t *testing.T) {
	testCases := []struct {
		name           string
		reqBody        interface{}
		expectedStatus int
		expectedBody   interface{}
		setupMocks     func()
	}{
		{
			name:           "422, Invalid JSON",
			reqBody:        `{invalid}`,
			expectedStatus: 422,
			expectedBody: orchestrator.Result{
				Error: "Invalid data",
			},
		},
		{
			name:           "500, Empty body",
			reqBody:        nil,
			expectedStatus: 500,
			expectedBody: orchestrator.Result{
				Error: "Internal server error",
			},
		},
		{
			name:           "404, Task not found",
			reqBody:        orchestrator.Result{Id: "123"},
			expectedStatus: 404,
			expectedBody: orchestrator.Result{
				Error: "Task not found or already processed",
			},
		},
		{
			name:           "200, Successful request",
			reqBody:        orchestrator.Result{Id: "123", Result: 42.0},
			expectedStatus: 200,
			expectedBody:   nil,
			setupMocks: func() {
				ch := make(chan orchestrator.Result, 1)
				processor.TaskResultChannels.Store("123", ch)
			},
		},
		{
			name:           "Channel type mismatch",
			reqBody:        orchestrator.Result{Id: "123"},
			expectedStatus: 500,
			expectedBody: orchestrator.Result{
				Error: "Internal server error",
			},
			setupMocks: func() {
				processor.TaskResultChannels.Store("123", "invalid channel")
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			if test.setupMocks != nil {
				test.setupMocks()
			}

			var reqBody []byte
			if test.reqBody != nil {
				var err error
				reqBody, err = json.Marshal(test.reqBody)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/result", bytes.NewReader(reqBody))
			rec := httptest.NewRecorder()

			handler := Result()
			handler.ServeHTTP(rec, req)

			if rec.Code != test.expectedStatus {
				t.Errorf("Expected status %d, got %d", test.expectedStatus, rec.Code)
			}

			if test.expectedBody != nil {
				var responseBody orchestrator.Result
				err := json.NewDecoder(rec.Body).Decode(&responseBody)
				if err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}

				if responseBody != test.expectedBody {
					t.Errorf("Expected body %+v, got %+v", test.expectedBody, responseBody)
				}
			} else {
				if rec.Body.Len() > 0 {
					t.Errorf("Expected empty body, got: %s", rec.Body.String())
				}
			}
		})
	}
}
