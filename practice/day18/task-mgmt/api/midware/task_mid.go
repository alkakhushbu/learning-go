package midware

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"task-mgmt/db"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
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

// Validation should not be in middle ware
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
func Log(next http.Handler) http.Handler {
	// return func(w http.ResponseWriter, r *http.Request) {
	// 	t := time.Now().UTC()
	// 	log.Println("Logging started", t)
	// 	log.Println(r.URL.Path, r.Method)
	// 	fn(w, r)
	// 	diff := time.Since(t)
	// 	log.Println("Logging ended", diff)
	// }

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now().UTC()
		log.Println("Logging started", t)
		log.Println(r.URL.Path, r.Method)
		next.ServeHTTP(w, r)
		diff := time.Since(t)
		log.Println("Logging ended", diff)
	})
}

func ReqIdMid(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := uuid.NewString()
		ctx := r.Context()                       // fetching the ctx object from the request
		ctx = context.WithValue(ctx, TaskId, id) // creating an updated ctx with a traceId store in it
		r = r.WithContext(ctx)                   // putting context inside the request object
		next(w, r)                               // calling next thing in the chain

	}
}
