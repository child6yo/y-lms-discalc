package orchestrator

import "time"

type (
	// Expression - арифметическое выражение. Может также содержать результат выражения и статус.
	Expression struct {
		ID         string  `json:"id"`
		Result     float64 `json:"result"`
		Expression string  `json:"expression"`
		Status     string  `json:"error"`
	}

	// ExpressionInput - арифметическое выражение, приходящее по HTTP.
	ExpressionInput struct {
		Expression string `json:"expression"`
	}

	// ExpressionID - целочисленное айди арифметического выражения.
	ExpressionID struct {
		ID int `json:"id"`
	}

	// ExpressionOutput - арифметическое выражение, передающееся клиенту по HTTP.
	ExpressionOutput struct {
		Expression Expression `json:"expression"`
	}

	// ExpressionListOutput - список арифметических выражений, передающихся клиенту по HTTP.
	ExpressionListOutput struct {
		Expressions []Expression `json:"expressions"`
	}

	// Task - задача на обработку участка арифметического выражения.
	Task struct {
		ID            string        `json:"id"`
		Arg1          float64       `json:"arg1"`
		Arg2          float64       `json:"arg2"`
		Operation     string        `json:"operation"`
		OperationTime time.Duration `json:"operation_time"`
	}

	// ErrorModel - ошибка, передающаяся по HTTP.
	ErrorModel struct {
		Error string `json:"error"`
	}

	// User - пользователь, проходящий регистрацию или аутентификацию.
	User struct {
		ID       int    `json:"id" db:"id"`
		Login    string `json:"login" db:"login"`
		Password string `json:"password" db:"password"`
	}
)
