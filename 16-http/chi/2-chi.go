package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

// add traceId middleware to trace the request to user routes
func main() {
	mux := chi.NewRouter()

	// set global middleware
	mux.Use(middleware.Logger, middleware.Recoverer)

	// localhost:8080/v1/users/123
	mux.Route("/v1/users", func(r chi.Router) {
		// get user
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {})

		//get user by id
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {})

		// create one user
		r.Post("/create", func(w http.ResponseWriter, r *http.Request) {})
	})

	// localhost:8080/v1/posts/123
	mux.Route("/v1/posts", func(r chi.Router) {
		r.Use(middleware.Logger, middleware.Recoverer)
		// fetch all posts
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {})
		// fetch post by id
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {})
		// create post
		r.Post("/create", func(w http.ResponseWriter, r *http.Request) {})
	})
}
