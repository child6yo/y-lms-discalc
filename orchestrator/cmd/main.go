package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/child6yo/y-lms-discalc/orchestrator"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/handler"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/processor"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/rpc"
	pb "github.com/child6yo/y-lms-discalc/orchestrator/proto"
	"google.golang.org/grpc"
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

func startHttpServer(port int, expressionInput chan orchestrator.ExpAndId) {
	http.HandleFunc("/api/v1/calculate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CulculateExpression(expressionInput)(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
	http.HandleFunc("/api/v1/expressions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.GetExpressions(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
	http.HandleFunc("/api/v1/expressions/", func(w http.ResponseWriter, r *http.Request) {
		pattern := `/api/v1/expressions/\d+`
		matched, err := regexp.MatchString(pattern, r.URL.Path)
		if err != nil || !matched {
			http.NotFound(w, r)
			return
		}
		if r.Method == http.MethodGet {
			handler.GetExpressionById(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	http.HandleFunc("/", handler.StaticFileHandler)

	httpPort := fmt.Sprint(":", port)

	log.Println("http server started at port: ", port)
	err := http.ListenAndServe(httpPort, nil)
	if err != nil {
		log.Println("error starting http server: ", err)
		os.Exit(1)
	}
}

func startGRPCServer(host, port string, taskChan chan orchestrator.Task) {
	addr := fmt.Sprintf("%s:%s", host, port)
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		log.Println("error starting tcp listener: ", err)
		os.Exit(1)
	}
	
	log.Println("tcp listener started at port: ", port)
	grpcServer := grpc.NewServer()
	taskServiceServer := rpc.NewServer(taskChan)
	pb.RegisterOrchestratorServiceServer(grpcServer, taskServiceServer)
	
	if err := grpcServer.Serve(lis); err != nil {
		log.Println("error serving grpc: ", err)
		os.Exit(1)
	}
}

func main() {
	config := map[string]time.Duration{}
	config["+"] = time.Duration(getIntEnv("TIME_ADDITION_MS", 100) * int(time.Millisecond))
	config["-"] = time.Duration(getIntEnv("TIME_SUBTRACTION_MS", 100) * int(time.Millisecond))
	config["*"] = time.Duration(getIntEnv("TIME_MULTIPLICATIONS_MS", 100) * int(time.Millisecond))
	config["/"] = time.Duration(getIntEnv("TIME_DIVISIONS_MS", 100) * int(time.Millisecond))

	expressionInput := make(chan orchestrator.ExpAndId, 10)
	expressionsMap := make(chan map[int]orchestrator.Expression, 10)
	tasks := make(chan orchestrator.Task, 30)

	go processor.StartExpressionProcessor(expressionInput, tasks, expressionsMap, config)
	go handler.HandleExpressionsChanel(expressionsMap)

	httpPort := 8000
	go startHttpServer(httpPort, expressionInput)
	
	gRPChost := "orchestrator"
	gRPCport := "5000"
	go startGRPCServer(gRPChost, gRPCport, tasks)

	log.Println("orchestrator successfully started")

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
