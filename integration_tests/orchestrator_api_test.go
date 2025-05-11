package integration_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"testing"
)

const srvHTTPtest = "http://localhost:8000"

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func registration(user User) *http.Response {
	url := srvHTTPtest + "/api/v1/register"

	requestJSON, _ := json.Marshal(user)
	res, err := http.Post(url, "application/json", bytes.NewBuffer(requestJSON))
	if err != nil {
		log.Println(err)
	}

	return res
}

func auth(user User) *http.Response {
	url := srvHTTPtest + "/api/v1/login"

	requestJSON, _ := json.Marshal(user)
	res, err := http.Post(url, "application/json", bytes.NewBuffer(requestJSON))
	if err != nil {
		log.Println(err)
	}

	return res
}

func TestRegistration(t *testing.T) {
	type Expect struct {
		Id int `json:"id"`
	}

	tests := []struct {
		name      string
		input     User
		expect    int
		expextErr bool
	}{
		{
			name: "ok",
			input: User{
				Login:    "login",
				Password: "123",
			},
			expect:    http.StatusCreated,
			expextErr: false,
		},
		{
			name: "same data",
			input: User{
				Login:    "login",
				Password: "123",
			},
			expect:    http.StatusInternalServerError,
			expextErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := registration(tt.input)

			var decRes Expect

			if res.StatusCode != tt.expect {
				t.Fatalf("test %s unexpected status: expected: %d, got %d", tt.name, http.StatusCreated, res.StatusCode)
			}
			if !tt.expextErr {
				if err := json.NewDecoder(res.Body).Decode(&decRes); err != nil {
					t.Fatalf("test %s failed: %s", tt.name, err.Error())
				}
			}
		})
	}
}

func TestAuth(t *testing.T) {
	type Expect struct {
		JWT string `json:"jwt"`
	}

	registration(User{Login: "abc", Password: "123"})

	tests := []struct {
		name      string
		input     User
		expect    int
		expextErr bool
	}{
		{
			name:  "ok",
			input: User{Login: "abc", Password: "123"},
			expect: 200,
			expextErr: false,
		},
		{
			name:  "invalid user",
			input: User{Login: "abcd", Password: "123"},
			expect: 400,
			expextErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := auth(tt.input)

			var decRes Expect

			if res.StatusCode != tt.expect {
				t.Fatalf("test %s unexpected status: expected: %d, got %d", tt.name, http.StatusCreated, res.StatusCode)
			}
			if !tt.expextErr {
				if err := json.NewDecoder(res.Body).Decode(&decRes); err != nil {
					t.Fatalf("test %s failed: %s", tt.name, err.Error())
				}
			}
		})
	}
}
