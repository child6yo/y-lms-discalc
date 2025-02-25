package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/child6yo/y-lms-discalc/orchestrator"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/processor"
)

func GetTask(output chan orchestrator.Task) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		select {
		case task := <-output:
			responseData, err := json.MarshalIndent(task, "", " ")
			if err != nil {
				httpNewError(w, 500, "Internal server error")
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(responseData)
		case <-time.After(2 * time.Second):
			w.WriteHeader(404)
			return
		}
	}
}

func Result() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req orchestrator.Result

		data, err := io.ReadAll(r.Body)
		if err != nil || len(data) == 0 {
			httpNewError(w, 500, "Internal server error")
			return
		}
		defer r.Body.Close()

		err = json.Unmarshal(data, &req)
		if err != nil {
			httpNewError(w, 422, "Invalid data")
			return
		}

		chInterface, ok := processor.TaskResultChannels.Load(req.Id)
		if !ok {
			httpNewError(w, 404, "Task not found or already processed")
			return
		}

		resultChan, ok := chInterface.(chan orchestrator.Result)
		if !ok {
			httpNewError(w, 500, "Internal server error")
			return
		}

		resultChan <- req

		w.WriteHeader(http.StatusOK)
	}
}
