package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/child6yo/y-lms-discalc/agent/pkg/service"
	"github.com/child6yo/y-lms-discalc/agent/pkg/worker"

	pb "github.com/child6yo/y-lms-discalc/shared/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Print("Failed to load env")
		return defaultValue
	}
	return value
}

func getIntEnv(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Print("Failed to load env: ", err)
		return defaultValue
	}
	return value
}

func main() {
	gRPChost := getEnv("GRPC_HOST", "localhost")
	gRPCport := getEnv("GRPC_PORT", "5000")
	addr := fmt.Sprintf("%s:%s", gRPChost, gRPCport)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println("could not connect to grpc server: ", err)
		os.Exit(1)
	}
	defer conn.Close()

	grpcClient := pb.NewOrchestratorServiceClient(conn)

	computingPower := getIntEnv("COMPUTING_POWER", 10)
	evaluater := new(service.EvaluateService)

	var wg sync.WaitGroup

	for w := 1; w <= computingPower; w++ {
		go func() {
			defer wg.Done()
			worker.Worker(w, grpcClient, evaluater)
		}()
		time.Sleep(1 * time.Second)
		wg.Add(1)
	}

	log.Print("Agent successfully started")

	wg.Wait()
}
