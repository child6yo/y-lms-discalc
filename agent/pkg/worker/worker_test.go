package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/child6yo/y-lms-discalc/agent"
)

// fakeTransport реализует http.RoundTripper и позволяет имитировать различные ответы.
type fakeTransport struct {
	t             *testing.T
	getCount      int
	postCount     int
	iterationChan chan struct{}
}

func (ft *fakeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	switch req.Method {
	case http.MethodGet:
		ft.getCount++
		switch ft.getCount {
		case 1:
			// Итерация 1: имитируем ошибку GET
			return nil, fmt.Errorf("simulated GET error")
		case 2:
			// Итерация 2: возвращаем ответ с не-OK статусом (например, 500)
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewBufferString("error")),
				Header:     make(http.Header),
			}, nil
		case 3:
			// Итерация 3: возвращаем OK, но с невалидным JSON
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("invalid json")),
				Header:     make(http.Header),
			}, nil
		case 4:
			// Итерация 4: возвращаем задание, приводящее к ошибке вычисления (деление на ноль)
			task := agent.Task{
				Id:        "task_div0",
				Operation: "/",
				Arg1:      5,
				Arg2:      0, // вызовет ошибку деления на ноль в EvaluatePostfix
				OperationTime: 3*time.Second,
			}
			data, err := json.Marshal(task)
			if err != nil {
				ft.t.Fatalf("Failed to marshal task: %v", err)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBuffer(data)),
				Header:     make(http.Header),
			}, nil
		case 5:
			// Итерация 5: возвращаем задание для успешной обработки, но POST имитируем ошибку
			task := agent.Task{
				Id:        "task_post_error",
				Operation: "+",
				Arg1:      2,
				Arg2:      3,
				OperationTime: 3*time.Second,
			}
			data, err := json.Marshal(task)
			if err != nil {
				ft.t.Fatalf("Failed to marshal task: %v", err)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBuffer(data)),
				Header:     make(http.Header),
			}, nil
		case 6:
			// Итерация 6: возвращаем задание, которое должно успешно обработаться
			task := agent.Task{
				Id:        "task_success",
				Operation: "+",
				Arg1:      2,
				Arg2:      3,
				OperationTime: 3*time.Second,
			}
			data, err := json.Marshal(task)
			if err != nil {
				ft.t.Fatalf("Failed to marshal task: %v", err)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBuffer(data)),
				Header:     make(http.Header),
			}, nil
		default:
			// Для последующих GET-запросов возвращаем не-OK, чтобы не обрабатывать их дальше
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewBufferString("error")),
				Header:     make(http.Header),
			}, nil
		}
	case http.MethodPost:
		ft.postCount++
		switch ft.postCount {
		case 1:
			// Итерация 5: имитируем ошибку POST
			return nil, fmt.Errorf("simulated POST error")
		case 2:
			// Итерация 6: проверяем содержимое POST и возвращаем успешный ответ
			var res agent.Result
			if err := json.NewDecoder(req.Body).Decode(&res); err != nil {
				ft.t.Errorf("Failed to decode POST payload: %v", err)
			}
			// Ожидаем, что EvaluatePostfix вернул: {Id:"task_success", Result:5, Error:""}
			expected := agent.Result{
				Id:     "task_success",
				Result: 5,
				Error:  "",
			}
			if res != expected {
				ft.t.Errorf("Expected POST payload %+v, got %+v", expected, res)
			}
			// Сигнализируем о завершении успешной итерации
			select {
			case ft.iterationChan <- struct{}{}:
			default:
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("ok")),
				Header:     make(http.Header),
			}, nil
		default:
			// Для последующих POST возвращаем успех
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("ok")),
				Header:     make(http.Header),
			}, nil
		}
	default:
		return nil, fmt.Errorf("unsupported method: %s", req.Method)
	}
}

func TestWorker(t *testing.T) {
	// Сохраним оригинальный транспорт и восстановим его по завершении теста.
	origTransport := http.DefaultClient.Transport
	defer func() {
		http.DefaultClient.Transport = origTransport
	}()

	// Канал, по которому мы получим сигнал успешной итерации (итерация 6).
	iterationChan := make(chan struct{}, 1)
	ft := &fakeTransport{
		t:             t,
		iterationChan: iterationChan,
	}

	// Переопределяем транспорт для http.DefaultClient, который используется в Worker.
	http.DefaultClient.Transport = ft

	// Запускаем Worker в отдельной горутине.
	// Заметим: Worker работает бесконечно, поэтому мы не можем его остановить,
	// но тест завершится, как только будет получен сигнал успешной итерации.
	go func() {
		// Здесь используем локальный адрес для теста
		Worker(1, "http://localhost:8080/task")
	}()

	// Ждём сигнала успешной обработки (итерация 6).
	select {
	case <-iterationChan:
		// Успех: Worker выполнил обработку задачи из итерации 6.
	case <-time.After(15 * time.Second):
		t.Fatal("Timeout waiting for Worker to process a task successfully")
	}
}

