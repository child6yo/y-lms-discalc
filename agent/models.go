package agent

import "time"

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
