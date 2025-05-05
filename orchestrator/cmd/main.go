package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"

	"github.com/child6yo/y-lms-discalc/orchestrator"
	h "github.com/child6yo/y-lms-discalc/orchestrator/pkg/handler"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/processor"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/repository"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/rpc"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/service"
	pb "github.com/child6yo/y-lms-discalc/shared/proto"
	"google.golang.org/grpc"

	_ "github.com/mattn/go-sqlite3"
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

type httpServer struct {
	server *http.Server
}

func (h *httpServer) startHTTPServer(port string, handler *h.Handler) error {
	http.HandleFunc("/api/v1/calculate", handler.AuthorizeMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CulculateExpression(w, r)
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
			handler.GetExpressionByID(w, r)
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

	h.server = &http.Server{
		Addr: ":" + port,
	}

	if err := h.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (h *httpServer) httpServerShutdown(ctx context.Context) error {
	return h.server.Shutdown(ctx)
}

type grpcServer struct {
	server *grpc.Server
}

func newGRPCServer() *grpcServer {
	return &grpcServer{server: grpc.NewServer()}
}

func (g *grpcServer) startGRPCServer(host, port string, taskChan *chan *orchestrator.Task) error {
	addr := fmt.Sprintf("%s:%s", host, port)
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		return err
	}

	taskServiceServer := rpc.NewServer(taskChan)
	pb.RegisterOrchestratorServiceServer(g.server, taskServiceServer)

	if err := g.server.Serve(lis); err != nil {
		return err
	}
	return nil
}

func (g *grpcServer) grpcServerShutdown() {
	g.server.GracefulStop()
}

func main() {
	// Загрузка конфигураций
	config := map[string]time.Duration{}
	config["+"] = time.Duration(getIntEnv("TIME_ADDITION_MS", 100) * int(time.Millisecond))
	config["-"] = time.Duration(getIntEnv("TIME_SUBTRACTION_MS", 100) * int(time.Millisecond))
	config["*"] = time.Duration(getIntEnv("TIME_MULTIPLICATIONS_MS", 100) * int(time.Millisecond))
	config["/"] = time.Duration(getIntEnv("TIME_DIVISIONS_MS", 100) * int(time.Millisecond))

	httpPort := getEnv("HTTP_PORT", "8000")

	gRPChost := getEnv("GRPC_HOST", "localhost")
	gRPCport := getEnv("GRPC_PORT", "5000")

	// Создание каналов
	expressionInput := make(chan *orchestrator.Expression, 30) // канал передачи арифметических выражений в обработчик
	tasks := make(chan *orchestrator.Task, 30)                 // канал передачи задач на вычисление участков выражений

	// Подключение к локальной базе данных sqlite
	db, err := repository.NewSqliteDb()
	if err != nil {
		log.Fatal("failed to connect sqlight: ", err)
	}

	// создание нового экземпляра репозитория
	repository := repository.NewRepository(db, 100)

	// сроздание нового экземпляра сервиса
	service := service.NewService(repository, &expressionInput)

	// создание нового экземпляра сервиса
	handler := h.NewHandler(service)

	// запуск горутины обработчика арифметических выражений
	go processor.StartExpressionProcessor(&expressionInput, &tasks, config, service)

	// создание нового экземпляра http сервера
	httpSrv := new(httpServer)

	// запуск http сервера
	go func() {
		if err := httpSrv.startHTTPServer(httpPort, handler); err != nil {
			log.Fatal("error starting http server: ", err)
		}
	}()
	log.Println("http server started at port: ", httpPort)

	// создание нового экземпляра gRPC сервера
	grpcSrv := newGRPCServer()

	// запуск gRPC сервера
	go func() {
		if err := grpcSrv.startGRPCServer(gRPChost, gRPCport, &tasks); err != nil {
			log.Fatal("error serving grpc: ", err)
		}
	}()
	log.Println("grpc tcp listener started at port: ", gRPCport)

	log.Println("orchestrator successfully started")

	// получение сигнала на остановку приложения (напр. Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	// graceful shutdown
	log.Print("orchestrator shutting down")

	if err = httpSrv.httpServerShutdown(context.Background()); err != nil {
		log.Print("error occured on server shutting down: ", err)
	}

	grpcSrv.grpcServerShutdown()

	if err = db.Close(); err != nil {
		log.Print("error occured on db connection close: ", err)
	}
}
