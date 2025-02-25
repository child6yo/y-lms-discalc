package handler

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/child6yo/y-lms-discalc/orchestrator"
)

type TestResponse1 struct {
	Id    int    `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func TestCalculateExpression(t *testing.T) {
	testCases := []struct {
		name         string
		requestBody  *strings.Reader
		expectBody   TestResponse1
		expectStatus int
	}{
		{
			name:         "201, Created",
			requestBody:  strings.NewReader(`{"expression":"2+2"}`),
			expectBody:   TestResponse1{Id: 1},
			expectStatus: 201,
		},
		{
			name:         "422, Expression is not valid 1",
			requestBody:  strings.NewReader(`{"expression":1}`),
			expectBody:   TestResponse1{Error: "Expression is not valid"},
			expectStatus: 422,
		},
		{
			name:         "422, Expression is not valid 2",
			requestBody:  strings.NewReader(`{"expression":"2("}`),
			expectBody:   TestResponse1{Error: "Expression is not valid"},
			expectStatus: 422,
		},
		{
			name:         "500, Internal server error",
			requestBody:  strings.NewReader(``),
			expectBody:   TestResponse1{Error: "Internal server error"},
			expectStatus: 500,
		},
	}

	c := make(chan orchestrator.ExpAndId, 5)

	for _, test := range testCases {
		req := httptest.NewRequest("POST", "http://localhost:8000/api/v1/calculate", test.requestBody)
		w := httptest.NewRecorder()
		CulculateExpression(c)(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Test %s: error reading response body: %v", test.name, err)
		}

		var r TestResponse1
		json.Unmarshal(body, &r)

		if r != test.expectBody {
			t.Errorf("test %s failed: result: %+v, expected: %+v", test.name, r, test.expectBody)
		} else if resp.StatusCode != test.expectStatus {
			t.Errorf("test %s failed: result: %d, expected: %d", test.name, resp.StatusCode, test.expectStatus)
		} else {
			t.Logf("Test %s success", test.name)
		}
	}
}

func TestGetExpressions(t *testing.T) {
	testCases := []struct {
		name         string
		exps         map[int]orchestrator.Expression
		expectBody   orchestrator.ExpressionList
		expectStatus int
	}{
		{
			name: "200, Non-empty expressions",
			exps: map[int]orchestrator.Expression{
				1: {Id: 1, Status: "Success", Result: 2},
				2: {Id: 2, Status: "ERROR", Result: 0},
			},
			expectBody: orchestrator.ExpressionList{
				Expressions: []orchestrator.Expression{
					{Id: 1, Status: "Success", Result: 2},
					{Id: 2, Status: "ERROR", Result: 0},
				},
			},
			expectStatus: 200,
		},
		{
			name:         "200, Empty expressions",
			exps:         map[int]orchestrator.Expression{},
			expectBody:   orchestrator.ExpressionList{Expressions: []orchestrator.Expression{}},
			expectStatus: 200,
		},
	}

	for _, test := range testCases {
		exps = test.exps
		req := httptest.NewRequest("GET", "http://localhost:8000/api/v1/expressions", nil)
		w := httptest.NewRecorder()

		GetExpressions(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Test %s: error reading response body: %v", test.name, err)
		}

		var result orchestrator.ExpressionList
		err = json.Unmarshal(body, &result)
		if err != nil {
			t.Fatalf("Test %s: error unmarshalling response body: %v", test.name, err)
		}

		if !reflect.DeepEqual(result, test.expectBody) {
			t.Errorf("test %s failed: result: %+v, expected: %+v", test.name, result, test.expectBody)
		} else if resp.StatusCode != test.expectStatus {
			t.Errorf("test %s failed: result: %d, expected: %d", test.name, resp.StatusCode, test.expectStatus)
		} else {
			t.Logf("Test %s success", test.name)
		}
	}
}

func TestGetExpressionById(t *testing.T) {
	testCases := []struct {
		name         string
		path         string
		exps         map[int]orchestrator.Expression
		expectBody   orchestrator.Expression
		expectStatus int
	}{
		{
			name: "200, Valid expression ID",
			path: "/api/v1/expressions/1",
			exps: map[int]orchestrator.Expression{
				1: {Id: 1, Status: "Success", Result: 0},
			},
			expectBody:   orchestrator.Expression{Id: 1, Status: "Success", Result: 0},
			expectStatus: 200,
		},
		{
			name: "404, Invalid expression ID",
			path: "/api/v1/expressions/99",
			exps: map[int]orchestrator.Expression{
				1: {Id: 1, Status: "Success", Result: 0},
			},
			expectBody:   orchestrator.Expression{},
			expectStatus: 404,
		},
		{
			name:         "500, Malformed URL",
			path:         "/api/v1/expressions/abc",
			exps:         map[int]orchestrator.Expression{},
			expectBody:   orchestrator.Expression{},
			expectStatus: 500,
		},
	}

	for _, test := range testCases {
		exps = test.exps
		req := httptest.NewRequest("GET", "http://localhost:8000"+test.path, nil)
		w := httptest.NewRecorder()

		GetExpressionById(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Test %s: error reading response body: %v", test.name, err)
		}

		var result orchestrator.Expression
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
