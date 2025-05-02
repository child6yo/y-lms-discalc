package repository

import (
	"database/sql/driver"
	"errors"
	"regexp"
	"slices"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/child6yo/y-lms-discalc/orchestrator"
)

func TestAddExpression(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("unable to make mock db: %s", err)
	}
	defer mockDB.Close()

	r := mainDatabase{db: mockDB}

	type args struct {
		userId     int
		expression *orchestrator.Expression
	}
	type mockBehavior func(args args)

	tests := []struct {
		name     string
		mockFunc mockBehavior
		input    args
		expect   int
		wantErr  bool
	}{
		{
			name: "successfull case",
			mockFunc: func(args args) {
				mock.ExpectExec("INSERT INTO expression").
					WithArgs(
						args.userId, args.expression.Expression,
						args.expression.Result, args.expression.Status,
					).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			input: args{
				userId: 1,
				expression: &orchestrator.Expression{
					Expression: "2+2",
					Result:     4,
					Status:     "Success",
				},
			},
			expect:  1,
			wantErr: false,
		},
		{
			name: "error handling",
			mockFunc: func(args args) {
				mock.ExpectExec("INSERT INTO expression").
					WithArgs(
						args.userId, args.expression.Expression,
						args.expression.Result, args.expression.Status,
					).WillReturnError(errors.New("insert error"))
			},
			input: args{
				userId:     0,
				expression: &orchestrator.Expression{},
			},
			expect:  0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.input)

			got, err := r.AddExpression(tt.input.userId, tt.input.expression)
			if err == nil && tt.wantErr {
				t.Errorf("expected error, got id %d", got)
			}

			if got != tt.expect {
				t.Errorf("expected: %d, got: %d", tt.expect, got)
			}
		})
	}
}

func TestGetExpressionById(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("unable to make mock db: %s", err)
	}
	defer mockDB.Close()

	r := mainDatabase{db: mockDB}

	type args struct {
		expId  int
		userId int
	}
	type mockBehavior func(args args)

	tests := []struct {
		name     string
		mockFunc mockBehavior
		input    args
		expect   orchestrator.Expression
		wantErr  bool
	}{
		{
			name: "successfull case",
			mockFunc: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "result", "exp", "status"}).
					AddRow(1, 4, "2+2", "Success")

				query := regexp.QuoteMeta("SELECT id, result, exp, status FROM expression WHERE user_id=$1 AND id=$2")
				mock.ExpectQuery(query).
					WithArgs(
						args.userId, args.expId,
					).
					WillReturnRows(rows)
			},
			input: args{
				expId:  1,
				userId: 1,
			},
			expect:  orchestrator.Expression{Id: "1", Result: 4, Expression: "2+2", Status: "Success"},
			wantErr: false,
		},
		{
			name: "error handling",
			mockFunc: func(args args) {
				query := regexp.QuoteMeta("SELECT id, result, exp, status FROM expression WHERE user_id=$1 AND id=$2")
				mock.ExpectQuery(query).
					WithArgs(
						args.userId, args.expId,
					).
					WillReturnError(errors.New("selecting error"))
			},
			input: args{
				expId:  1,
				userId: 1,
			},
			expect:  orchestrator.Expression{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.input)

			got, err := r.GetExpressionById(tt.input.expId, tt.input.userId)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got: %v", got)
				}
			} else {
				if got == nil {
					t.Error("expected non-nil result, got nil, error: ", err)
				} else if *got != tt.expect {
					t.Errorf("expected: %v, got: %v", tt.expect, got)
				}
			}
		})
	}
}

func TestGetExpressions(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("unable to make mock db: %s", err)
	}
	defer mockDB.Close()

	r := mainDatabase{db: mockDB}

	type args struct {
		userId int
	}
	type mockBehavior func(args args)

	tests := []struct {
		name     string
		mockFunc mockBehavior
		input    args
		expect   *[]orchestrator.Expression
		wantErr  bool
	}{
		{
			name: "successfull case",
			mockFunc: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "result", "exp", "status"}).
					AddRows(
						[]driver.Value{1, 4, "2+2", "Success"},
						[]driver.Value{2, 0, "2+", "ERROR"},
					)
				
				query := regexp.QuoteMeta("SELECT id, result, exp, status FROM expression WHERE user_id=$1")
				mock.ExpectQuery(query).WithArgs(args.userId).WillReturnRows(rows)
			},
			input: args{
				userId: 1,
			},
			expect: &[]orchestrator.Expression{
				{
					Id: "1",
					Result: 4,
					Expression: "2+2",
					Status: "Success",
				},
				{
					Id: "2",
					Result: 0,
					Expression: "2+",
					Status: "ERROR",
				},
			},
			wantErr: false,
		},
		{
			name: "error handling",
			mockFunc: func(args args) {
				query := regexp.QuoteMeta("SELECT id, result, exp, status FROM expression WHERE user_id=$1")
				mock.ExpectQuery(query).WithArgs(args.userId).WillReturnError(errors.New("selecting error"))
			},
			input: args{
				userId: 0,
			},
			expect: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(tt.input)

			got, err := r.GetExpressions(tt.input.userId)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got: %v", got)
				}
			} else {
				if got == nil {
					t.Error("expected non-nil result, got nil, error: ", err)
				} else if !slices.Equal(*got, *tt.expect) {
					t.Errorf("expected: %v, got: %v", tt.expect, got)
				}
			}
		})
	}
}