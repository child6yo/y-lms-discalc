package agent

import "time"

// Task - задача, возвращающаяся при вызове GetTask в оркестраторе через gRPC.
type Task struct {
	ID            string        `json:"id"`
	Arg1          float64       `json:"arg1"`
	Arg2          float64       `json:"arg2"`
	Operation     string        `json:"operation"`
	OperationTime time.Duration `json:"operation_time"`
}

// Result - результат отработки задачи,
// передающийся через вызов функции TakeResult оркестратора через gRPC.
type Result struct {
	ID     string  `json:"id"`
	Result float64 `json:"result"`
	Error  string  `json:"error"`
}
