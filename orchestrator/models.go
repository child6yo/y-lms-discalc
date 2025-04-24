package orchestrator

import "time"

type (
	ExpressionInput struct {
		Expression string `json:"expression"`
	}

	ExpressionId struct {
		Id int `json:"id"`
	}

	ExpAndId struct {
		Id         int
		Expression []string
	}

	Expression struct {
		Id     int     `json:"id"`
		Status string  `json:"status"`
		Result float64 `json:"result"`
	}

	ExpressionOutput struct {
		Expression Expression `json:"expression"`
	}

	ExpressionList struct {
		Expressions []Expression `json:"expressions"`
	}

	Task struct {
		Id            string        `json:"id"`
		Arg1          float64       `json:"arg1"`
		Arg2          float64       `json:"arg2"`
		Operation     string        `json:"operation"`
		OperationTime time.Duration `json:"operation_time"`
	}

	Result struct {
		Id         string  `json:"id"`
		Result     float64 `json:"result"`
		Expression string  `json:"expression"`
		Status      string  `json:"error"`
	}

	ErrorModel struct {
		Error string `json:"error"`
	}

	User struct {
		Id       int    `json:"id" db:"id"`
		Login    string `json:"login" db:"login"`
		Password string `json:"password" db:"password"`
	}
)
