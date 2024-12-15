package midware

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"task-mgmt/db"

	"github.com/go-chi/chi"
)

type ContextKey string

var TaskId ContextKey = "taskid"

// func TraceId(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		id := uuid.NewString()
// 		ctx := context.Background()
// 		ctx = context.WithValue(ctx, TaskId, id)
// 		next.ServeHTTP(w, r)
// 	})
// }

func ValidateTaskBody(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Found error in reading request body data byte:%s\n", err.Error())
			http.Error(w, "Found error in reading request body data byte", http.StatusBadRequest)
			return
		}
		task := db.Task{}
		err = json.Unmarshal(data, &task)
		if err != nil {
			log.Printf("Invalid request body, cannot convert to json:%s\n", err.Error())
			http.Error(w, "Invalid request body, cannot convert to json", http.StatusBadRequest)
			return
		}
		log.Println("Task in middleware ValidateTaskBody:", task)
		next(w, r)
	}
}

func ValidateTaskId(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		id := chi.URLParam(r, "id")
		log.Printf("Reading task for id: %s\n", id)
		_, err := strconv.Atoi(id)
		if err != nil {
			msg := fmt.Sprintf("Invalid task id: %s, %s", id, err.Error())
			http.Error(w, msg, http.StatusNoContent)
			return
		}
		log.Printf("Task found with id:%s", id)
		next(w, r)
	}
}
func Logging(http.Handler) http.Handler {
	panic("Not implemented")
}
