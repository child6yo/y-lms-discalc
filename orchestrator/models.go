package orchestrator

import "time"

type ExpressionInput struct {
	Expression string `json:"expression"`
}

type ExpressionId struct {
	Id int `json:"id"`
}

type Expression struct {
	Id     int    `json:"id"`
	Status string `json:"status"`
	Result string `json:"result"`
}

type ExpressionOutput struct {
	Expression Expression `json:"expression"`
}

type ExpressionList struct {
	Expressions []Expression `json:"expressions"`
}

type Task struct {
	Id            int       `json:"id"`
	Arg1          string    `json:"arg1"`
	Arg2          string    `json:"arg2"`
	Operation     string    `json:"operation"`
	OperationTime time.Time `json:"operation_time"`
}

type Result struct {
	Id     int     `json:"id"`
	Result float64 `json:"result"`
}
