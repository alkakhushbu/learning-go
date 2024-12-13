package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func StartHandlerService() {
	mux := chi.NewRouter()
	// mux.Use(service.TraceId, service.Logging)
	mux.Route("/api/v1", func(r chi.Router) {
		r.Post("/tasks", createTask)
		r.Get("/tasks/{id}", getTaskById)
		r.Get("/tasks", getAllTasks)
		r.Put("/tasks/{id}", updateTaskById)
		r.Delete("/tasks/{id}", deleteTaskById)
		r.Patch("/tasks/{id}/status", updateTaskStatus)
	})
	http.ListenAndServe(":8084", mux)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	panic("Not implemented")
}

func getTaskById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	log.Printf("Got id: %s\n", id)
	// panic("Not implemented")
}
func getAllTasks(w http.ResponseWriter, r *http.Request) {
	panic("Not implemented")
}
func updateTaskById(w http.ResponseWriter, r *http.Request) {
	panic("Not implemented")
}
func deleteTaskById(w http.ResponseWriter, r *http.Request) {
	panic("Not implemented")
}
func updateTaskStatus(w http.ResponseWriter, r *http.Request) {
	panic("Not implemented")
}
