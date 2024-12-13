package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var users = map[int]string{1: "Jon", 2: "Jane", 3: "Jack", 4: "Jacob"}
var posts = map[int]string{1: "Hello", 2: "Welcome", 3: "RSVP"}

// add traceId middleware to trace the request to user routes
func main() {
	mux := chi.NewRouter()

	// set global middleware
	// mux.Use(middleware.Logger, middleware.Recoverer, traceId)
	mux.Use(middleware.Logger, traceId)

	// localhost:8080/v1/users/123
	mux.Route("/v1/users", func(r chi.Router) {
		// get user
		r.Get("/", getUsers)

		//get user by id
		r.Get("/{id}", getUserById)

		// create one user
		r.Post("/create", func(w http.ResponseWriter, r *http.Request) {})
	})

	// localhost:8080/v1/posts/123
	mux.Route("/v1/posts", func(r chi.Router) {
		// r.Use(middleware.Logger, middleware.Recoverer)
		// fetch all posts
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {})
		// fetch post by id
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {})
		// create post
		r.Post("/create", createPost)
	})
	http.ListenAndServe(":8084", mux)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "error in json marshal users to slice of bytes", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "error in writing databytes by responsewriter", http.StatusInternalServerError)
		return
	}
}

func getUserById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid user id ", http.StatusNoContent)
		return
	}
	user, ok := users[id]
	if !ok {
		http.Error(w, "invalid user id ", http.StatusNoContent)
		return
	}
	data, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "error in json marshal for user", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error in reading request body", http.StatusBadRequest)
		return
	}
	var m map[string]string
	err = json.Unmarshal(data, &m)
	if err != nil {
		fmt.Printf("Data after marshal:%v\n", m)
		http.Error(w, "error in unmarshal bytedata", http.StatusBadRequest)
		return
	}
	for k, v := range m {
		key, err := strconv.Atoi(k)
		if err != nil {
			continue
		}
		posts[key] = v
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func traceId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Inside traceid:\n\n")
		// http can recover from panic
		// but if a goroutine that we run manually panics, service fails
		// go panic("Panic happened")
		// panic("Panic happened")

		log.Fatal("Fatal message")
		next.ServeHTTP(w, r)
	})
}
