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
	h "github.com/child6yo/y-lms-discalc/orchestrator/pkg/handler"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/processor"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/repository"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/rpc"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/service"
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

func startHttpServer(port int, handler *h.Handler) {
	http.HandleFunc("/api/v1/calculate", handler.AuthorizeMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CulculateExpression()(w, r)
		} else {
			http.NotFound(w, r)
		}
	}))
	http.HandleFunc("/api/v1/expressions", handler.AuthorizeMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.GetExpressions(w, r)
		} else {
			http.NotFound(w, r)
		}
	}))
	http.HandleFunc("/api/v1/expressions/", handler.AuthorizeMiddleware(func(w http.ResponseWriter, r *http.Request) {
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
	}))
	http.HandleFunc("/api/v1/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateUser(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
	http.HandleFunc("/api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.Auth(w, r)
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

	expressionInput := make(chan *orchestrator.Result, 10)
	tasks := make(chan orchestrator.Task, 30)

	httpPort := getIntEnv("HTTP_PORT", 8000)

	gRPChost := getEnv("GRPC_HOST", "orchestrator")
	gRPCport := getEnv("GRPC_PORT", "5000")

	repository, err := repository.NewRepository(100)
	if err != nil {
		log.Println("failed to connect sqlight: ", err)
		os.Exit(1)
	}
	defer repository.Db.Close()
	service := service.NewService(repository, &expressionInput)
	handler := h.NewHandler(service)

	go processor.StartExpressionProcessor(&expressionInput, tasks, config, service)
	go startHttpServer(httpPort, handler)
	go startGRPCServer(gRPChost, gRPCport, tasks)

	log.Println("orchestrator successfully started")
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
