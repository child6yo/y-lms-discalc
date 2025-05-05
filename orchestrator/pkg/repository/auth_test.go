package repository

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/child6yo/y-lms-discalc/orchestrator"
)

func TestCreateUser(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("unable to make mock db: %s", err)
	}
	defer mockDB.Close()

	r := mainDatabase{db: mockDB}

	type args struct {
		user orchestrator.User
	}
	type mockBehavior func(args args, id int)

	tests := []struct {
		name     string
		mockFunc mockBehavior
		input    args
		expect   int
		wantErr  bool
	}{
		{
			name: "successfull case",
			mockFunc: func(args args, id int) {
				mock.ExpectExec("INSERT INTO user").
					WithArgs(
						args.user.Login, args.user.Password,
					).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			input: args{
				orchestrator.User{Login: "test", Password: "hash"},
			},
			expect:  1,
			wantErr: false,
		},
		{
			name: "error handling",
			mockFunc: func(args args, id int) {
				mock.ExpectExec("INSERT INTO user").
					WithArgs(
						args.user.Login, args.user.Password,
					).WillReturnError(errors.New("insert error"))
			},
			input: args{
				orchestrator.User{Login: "test", Password: "hash"},
			},
			expect:  0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.input, tt.expect)

			got, err := r.CreateUser(tt.input.user)
			if err == nil && tt.wantErr {
				t.Errorf("expected error, got id %d", got)
			}

			if got != tt.expect {
				t.Errorf("expected: %d, got: %d", tt.expect, got)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("unable to make mock db: %s", err)
	}
	defer mockDB.Close()

	r := mainDatabase{db: mockDB}

	type args struct {
		login    string
		password string
	}
	type mockBehavior func(args args, id int, login string, password string)

	tests := []struct {
		name     string
		mockFunc mockBehavior
		input    args
		expect   orchestrator.User
		wantErr  bool
	}{
		{
			name: "successfull case",
			mockFunc: func(args args, id int, login string, password string) {
				rows := sqlmock.NewRows([]string{"id", "login", "password"}).AddRow(id, login, password)

				query := regexp.QuoteMeta("SELECT * FROM user WHERE login=$1 AND password=$2")
				mock.ExpectQuery(query).
					WithArgs(
						args.login, args.password,
					).WillReturnRows(rows)
			},
			input: args{
				login:    "login",
				password: "hash",
			},
			expect:  orchestrator.User{ID: 1, Login: "login", Password: "hash"},
			wantErr: false,
		},
		{
			name: "user not found",
			mockFunc: func(args args, id int, login string, password string) {
				mock.ExpectQuery("SELECT \\* FROM user WHERE login=\\$1 AND password=\\$2").
					WithArgs(
						args.login, args.password,
					).WillReturnError(errors.New("user not found"))
			},
			input: args{
				login:    "login",
				password: "hash",
			},
			expect:  orchestrator.User{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.input, tt.expect.ID, tt.expect.Login, tt.expect.Password)

			got, err := r.GetUser(tt.input.login, tt.input.password)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got: %v", got)
				}
			} else {
				if got == nil {
					t.Error("expected non-nil result, got nil")
				} else if *got != tt.expect {
					t.Errorf("expected: %v, got: %v", tt.expect, got)
				}
			}
		})
	}
}
