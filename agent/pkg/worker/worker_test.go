package worker

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/child6yo/y-lms-discalc/agent"
	pb "github.com/child6yo/y-lms-discalc/shared/proto"
	"google.golang.org/grpc"
)

type mockOrchestratorClient struct {
	getTaskFunc    func(context.Context, *pb.Empty, ...grpc.CallOption) (*pb.TaskRequest, error)
	takeResultFunc func(context.Context, *pb.ResultResponse, ...grpc.CallOption) (*pb.Empty, error)
}

func (m *mockOrchestratorClient) GetTask(ctx context.Context, in *pb.Empty, opts ...grpc.CallOption) (*pb.TaskRequest, error) {
	return m.getTaskFunc(ctx, in, opts...)
}

func (m *mockOrchestratorClient) TakeResult(ctx context.Context, in *pb.ResultResponse, opts ...grpc.CallOption) (*pb.Empty, error) {
	return m.takeResultFunc(ctx, in, opts...)
}

type MockEvaluator struct {
	EvaluateFunc func(task agent.Task) agent.Result
}

func (m *MockEvaluator) PostfixEvaluate(task agent.Task) agent.Result {
	return m.EvaluateFunc(task)
}

func TestWorker(t *testing.T) {
	tests := []struct {
		name          string
		getTask       func(ctx context.Context, in *pb.Empty, opts ...grpc.CallOption) (*pb.TaskRequest, error)
		evaluate      func(task agent.Task) agent.Result
		takeResult    func(ctx context.Context, in *pb.ResultResponse, opts ...grpc.CallOption) (*pb.Empty, error)
		expectSuccess bool
	}{
		{
			name: "Successful task processing",
			getTask: func(ctx context.Context, in *pb.Empty, opts ...grpc.CallOption) (*pb.TaskRequest, error) {
				return &pb.TaskRequest{
					Id:            "1",
					Arg1:          10,
					Arg2:          5,
					Operation:     "+",
					OperationTime: 1,
				}, nil
			},
			evaluate: func(task agent.Task) agent.Result {
				return agent.Result{Id: task.Id, Result: 15}
			},
			takeResult: func(ctx context.Context, in *pb.ResultResponse, opts ...grpc.CallOption) (*pb.Empty, error) {
				if in.Id != "1" || in.Result != 15 {
					t.Errorf("Unexpected result: %v", in)
				}
				return &pb.Empty{}, nil
			},
			expectSuccess: true,
		},
		{
			name: "Task evaluation error",
			getTask: func(ctx context.Context, in *pb.Empty, opts ...grpc.CallOption) (*pb.TaskRequest, error) {
				return &pb.TaskRequest{
					Id:            "2",
					Arg1:          10,
					Arg2:          0,
					Operation:     "/",
					OperationTime: 1,
				}, nil
			},
			evaluate: func(task agent.Task) agent.Result {
				return agent.Result{Id: task.Id, Error: "division by zero"}
			},
			takeResult: func(ctx context.Context, in *pb.ResultResponse, opts ...grpc.CallOption) (*pb.Empty, error) {
				t.Error("TakeResult should not be called on evaluation error")
				return nil, nil
			},
			expectSuccess: false,
		},
		{
			name: "GetTask network error",
			getTask: func(ctx context.Context, in *pb.Empty, opts ...grpc.CallOption) (*pb.TaskRequest, error) {
				return nil, context.DeadlineExceeded
			},
			evaluate: func(task agent.Task) agent.Result {
				t.Error("Evaluate should not be called if GetTask fails")
				return agent.Result{}
			},
			takeResult: func(ctx context.Context, in *pb.ResultResponse, opts ...grpc.CallOption) (*pb.Empty, error) {
				t.Error("TakeResult should not be called if GetTask fails")
				return nil, nil
			},
			expectSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var once sync.Once
			done := make(chan struct{})

			grpcClient := &mockOrchestratorClient{
				getTaskFunc: func(ctx context.Context, in *pb.Empty, opts ...grpc.CallOption) (*pb.TaskRequest, error) {
					once.Do(func() { close(done) })
					return tt.getTask(ctx, in, opts...)
				},
				takeResultFunc: tt.takeResult,
			}
			evaluator := &MockEvaluator{EvaluateFunc: tt.evaluate}

			go Worker(1, grpcClient, evaluator)

			select {
			case <-done:
			case <-time.After(500 * time.Millisecond):
				t.Error("Worker did not process task in time")
			}
		})
	}
}
