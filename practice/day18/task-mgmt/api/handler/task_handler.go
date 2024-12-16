package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"task-mgmt/api/midware"
	"task-mgmt/db"

	"github.com/go-chi/chi"
)

func StartHandlerService() {
	log.Println("Inside StartHandlerService function")
	mux := chi.NewRouter()
	// mux.Use(midware.Log)
	mux.Route("/api/v1", func(r chi.Router) {
		r.Post("/tasks", midware.ReqIdMid(createTask))
		r.Get("/tasks/{id}", midware.ReqIdMid(getTaskById))
		r.Get("/tasks", midware.ReqIdMid(getAllTasks))
		r.Put("/tasks/{id}", midware.ReqIdMid(updateTaskById))
		r.Delete("/tasks/{id}", midware.ReqIdMid(deleteTaskById))
		r.Patch("/tasks/{id}/status", midware.ReqIdMid(updateTaskStatus))
	})
	http.ListenAndServe(":8084", mux)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Found error in reading request body data byte:%s\n", err.Error())
		http.Error(w, "Found error in reading request body data byte", http.StatusBadRequest)
		return
	}
	newTask := &db.NewTask{}
	err = json.Unmarshal(data, newTask)
	if err != nil {
		log.Printf("Invalid request body, cannot convert to json:%s\n", err.Error())
		http.Error(w, "Invalid request body, cannot convert to json", http.StatusBadRequest)
		return
	}
	log.Println("Task in middleware ValidateTaskBody:", newTask)

	task, err := db.CreateTask(r.Context(), newTask)

	if err != nil {
		msg := fmt.Sprintf("Task creation failed: %s", err.Error())
		http.Error(w, msg, http.StatusNoContent)
		return
	}
	data, err = json.Marshal(task)
	if err != nil {
		log.Println("Found error in json Marshal for task")
		http.Error(w, "Found error in json Marshal for task", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func getTaskById(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		msg := fmt.Sprintf("Invalid task id: %d, %s", id, err.Error())
		http.Error(w, msg, http.StatusNoContent)
		return
	}
	log.Printf("Task found with id:%d", id)
	task, err := db.GetTaskById(r.Context(), id)
	if err != nil {
		msg := fmt.Sprintf("Invalid task id: %d, %s", id, err.Error())
		http.Error(w, msg, http.StatusNoContent)
		return
	}
	data, err := json.Marshal(task)
	if err != nil {
		log.Printf("Found error in marshaling the task object: %v, error: %s", task, err.Error())
		msg := fmt.Sprintf("Found error in marshaling the task object: %v", task)
		http.Error(w, msg, http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}
func getAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	tasks, err := db.GetAllTasks(r.Context())
	if err != nil {
		log.Printf("Found error in fetching all tasks.. Error: %s", err.Error())
		http.Error(w, "Found error in fetching all tasks..", http.StatusNoContent)
		return
	}
	log.Println("All tasks received from DB:", tasks)
	data, err := json.Marshal(tasks)
	if err != nil {
		log.Printf("Found error in marshaling the slice of tasks.. error: %s", err.Error())
		http.Error(w, "Found error in marshaling the slice of tasks..", http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}

func updateTaskById(w http.ResponseWriter, r *http.Request) {
	//validating id
	w.Header().Add("Content-Type", "application/json")
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		msg := fmt.Sprintf("Invalid task id: %d, %s", id, err.Error())
		http.Error(w, msg, http.StatusNoContent)
		return
	}
	log.Printf("Task found with id:%d", id)

	//validating body
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Found error in reading request body data byte:%s\n", err.Error())
		http.Error(w, "Found error in reading request body data byte", http.StatusBadRequest)
		return
	}
	newTask := &db.NewTask{}
	err = json.Unmarshal(data, newTask)
	if err != nil {
		log.Printf("Invalid request body, cannot convert to json:%s\n", err.Error())
		http.Error(w, "Invalid request body, cannot convert to json", http.StatusBadRequest)
		return
	}
	log.Println("Task in middleware ValidateTaskBody:", newTask)

	//db layer call
	task, err := db.UpdateTask(r.Context(), id, newTask)

	//validation post db layer call
	if err != nil {
		log.Println("Found error in updating task with id:", id, err)
		http.Error(w, "Found error in updating the task", http.StatusInternalServerError)
		return
	}
	data, err = json.Marshal(task)
	if err != nil {
		log.Println("Found error in json Marshal for task")
		http.Error(w, "Found error in json Marshal for task", http.StatusBadRequest)
		return
	}

	//write response
	log.Println("Task updated with id:", id)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func deleteTaskById(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		msg := fmt.Sprintf("Invalid task id: %d, %s", id, err.Error())
		http.Error(w, msg, http.StatusNoContent)
		return
	}
	log.Printf("Task found with id:%d", id)
	err = db.DeleteTask(r.Context(), id)
	if err != nil {
		log.Println("Found error in deleting task with id:", id, err)
		http.Error(w, "Found error in deleting the task:"+err.Error(), http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deleted task"))
}

func updateTaskStatus(w http.ResponseWriter, r *http.Request) {
	//validating id
	w.Header().Add("Content-Type", "application/json")
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		msg := fmt.Sprintf("Invalid task id: %d, %s", id, err.Error())
		http.Error(w, msg, http.StatusNoContent)
		return
	}
	log.Printf("Task found with id:%d", id)

	//validating body
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Found error in reading request body data byte:%s\n", err.Error())
		http.Error(w, "Found error in reading request body data byte", http.StatusBadRequest)
		return
	}
	newTask := &db.NewTask{}
	err = json.Unmarshal(data, newTask)
	if err != nil {
		log.Printf("Invalid request body, cannot convert to json:%s\n", err.Error())
		http.Error(w, "Invalid request body, cannot convert to json", http.StatusBadRequest)
		return
	}
	log.Println("Task in middleware ValidateTaskBody:", newTask)

	//db layer call
	task, err := db.UpdateTaskStatus(r.Context(), id, newTask)

	//validation post db layer call
	if err != nil {
		log.Println("Found error in updating task with id:", id, err)
		http.Error(w, "Found error in updating the task", http.StatusInternalServerError)
		return
	}
	data, err = json.Marshal(task)
	if err != nil {
		log.Println("Found error in json Marshal for task")
		http.Error(w, "Found error in json Marshal for task", http.StatusBadRequest)
		return
	}
	//write response
	log.Println("Task updated with id:", id)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
