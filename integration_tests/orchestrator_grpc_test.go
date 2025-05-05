package integration_test

import (
	"context"
	"strings"
	"testing"

	pb "github.com/child6yo/y-lms-discalc/shared/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Данные тесты проверяют пингуются ли функции
// без полноценных проверок их функционала

const srvGrpcUrl = "127.0.0.1:5000"

func TestOrchestratorGetTask(t *testing.T) {
	conn, err := grpc.NewClient(srvGrpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal("could not connect to grpc server: ", err)
	}

	client := pb.NewOrchestratorServiceClient(conn)

	_, err = client.GetTask(context.TODO(), nil)

	if strings.Contains(err.Error(), "time limit exceeded") {
		return
	}
	t.Fatal("unexpected server behavior. ", err)
}

func TestOrchestratorTakeResult(t *testing.T) {
	conn, err := grpc.NewClient(srvGrpcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal("could not connect to grpc server: ", err)
	}

	client := pb.NewOrchestratorServiceClient(conn)

	result := &pb.ResultResponse{Id: "1", Result: 0, Error: ""}
	_, err = client.TakeResult(context.TODO(), result)

	if strings.Contains(err.Error(), "task not found or already processed") {
		return
	}
	t.Fatal("unexpected server behavior. ", err)
}