package rpc

import (
	"context"
	"errors"
	"time"

	"github.com/child6yo/y-lms-discalc/orchestrator"
	"github.com/child6yo/y-lms-discalc/orchestrator/internal/app/processor"
	pb "github.com/child6yo/y-lms-discalc/orchestrator/proto"
)

// OrchestratorServiceServer определяет интерфейс gRPC сервера оркестратора.
type OrchestratorServiceServer interface {
	// GetTask возвращает задачу на вычисление арифметического выражения, если она имеется.
	GetTask(context.Context, *pb.Empty) (*pb.TaskRequest, error)

	// TakeResult принимает результат вычисления арифметического выражения.
	TakeResult(context.Context, *pb.ResultResponse) (*pb.Empty, error)
	mustEmbedUnimplementedOrchestratorServiceServer()
}

// Server реализует интерфейс gRPC сервера оркестратора.
type Server struct {
	pb.OrchestratorServiceServer
	taskChan *chan *orchestrator.Task
}

// NewServer создает новый экземпляр сервера оркестратора.
//
// Параметры:
//   - taskChan: указатель на канал передачи задач на вычисление арифметических выражений.
func NewServer(taskChan *chan *orchestrator.Task) *Server {
	return &Server{taskChan: taskChan}
}

// GetTask возвращает задачу на вычисление арифметического выражения, если она имеется.
func (s *Server) GetTask(ctx context.Context, _ *pb.Empty) (*pb.TaskRequest, error) {
	select {
	case task := <-*s.taskChan:
		return &pb.TaskRequest{Id: task.ID,
			Arg1:          float32(task.Arg1),
			Arg2:          float32(task.Arg2),
			Operation:     task.Operation,
			OperationTime: int64(task.OperationTime)}, nil
	case <-time.After(3 * time.Second):
		return nil, errors.New("time limit exceeded")
	}
}

// TakeResult принимает результат вычисления арифметического выражения.
func (s *Server) TakeResult(ctx context.Context, result *pb.ResultResponse) (*pb.Empty, error) {
	chInterface, ok := processor.TaskResultChannels.Load(result.Id)
	if !ok {
		return nil, errors.New("task not found or already processed")
	}

	resultChan, ok := chInterface.(chan orchestrator.Expression)
	if !ok {
		return nil, errors.New("something went wrong")
	}

	res := orchestrator.Expression{ID: result.Id, Result: float64(result.Result), Status: result.Error}
	resultChan <- res

	return nil, nil
}
