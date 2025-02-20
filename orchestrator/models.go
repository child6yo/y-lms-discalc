package orchestrator

import "time"

type ExpressionInput struct {
	Expression string `json:"expression"`
}

type ExpressionId struct {
	Id int `json:"id"`
}

type ExpAndId struct {
	Id         int
	Expression []string
}

type Expression struct {
	Id     int     `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}

type ExpressionOutput struct {
	Expression Expression `json:"expression"`
}

type ExpressionList struct {
	Expressions []Expression `json:"expressions"`
}

type Task struct {
	Id            string        `json:"id"`
	Arg1          float64       `json:"arg1"`
	Arg2          float64       `json:"arg2"`
	Operation     string        `json:"operation"`
	OperationTime time.Duration `json:"operation_time"`
}

type Result struct {
	Id     string  `json:"id"`
	Result float64 `json:"result"`
	Error  string  `json:"error"`
}

type ErrorModel struct {
	Error string `json:"error"`
}
