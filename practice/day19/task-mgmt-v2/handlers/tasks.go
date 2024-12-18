package handlers

import (
	"log"
	"net/http"
	"strconv"
	"task-mgmt-v2/models"
	"task-mgmt-v2/pkg"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func (h *Handler) CreateTask(c *gin.Context) {
	// ctx := c.Request.Context()
	// traceId := pkg.GetTraceId(ctx)

	// Start a new span for the handler
	tracer := otel.Tracer("user-micro")
	ctx, span := tracer.Start(c.Request.Context(), "Handler-CreateTask")
	defer span.End()

	//get traceId from tracer span
	traceId := span.SpanContext().TraceID().String()

	nt := models.NewTask{}
	err := c.ShouldBindWith(&nt, binding.JSON)
	if err != nil {
		slog.Error("Error in task binding:",
			slog.String("Error:", err.Error()),
			slog.String("TraceId", traceId))

		// Handle and record any errors in the span
		span.SetAttributes(semconv.HTTPResponseStatusCodeKey.Int(http.StatusBadRequest)) // HTTP 400
		span.SetAttributes(attribute.String("traceId", traceId))
		span.SetStatus(codes.Error, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	//validate request body
	err = h.validate.Struct(&nt)
	if err != nil {
		slog.Error("Invalid request body, please provide required fields:",
			slog.String("Error:", err.Error()),
			slog.String("TraceId", traceId))

		// Handle and record any errors in the span
		span.SetAttributes(semconv.HTTPResponseStatusCodeKey.Int(http.StatusBadRequest)) // HTTP 400
		span.SetAttributes(attribute.String("traceId", traceId))
		span.SetStatus(codes.Error, err.Error())

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	//create task in db
	task, err := h.service.CreateTask(ctx, nt)

	// validation post db layer call
	if err != nil {
		slog.Error("Error in creating new task:",
			slog.String("Error:", err.Error()),
			slog.String("TraceId", traceId))

		// Handle and record any errors in the span
		span.SetAttributes(semconv.HTTPResponseStatusCodeKey.Int(http.StatusBadRequest)) // HTTP 400
		span.SetAttributes(attribute.String("traceId", traceId))
		span.SetStatus(codes.Error, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error in creating new task"})
		return
	}
	log.Println("New task created:", task)
	span.SetAttributes(semconv.HTTPResponseStatusCodeKey.Int(http.StatusCreated)) // HTTP 201 Created
	c.JSON(http.StatusCreated, task)
}

func (h *Handler) getTaskById(c *gin.Context) {
	// w.Header().Add("Content-Type", "application/json")
	// c.Param("id")
	ctx := c.Request.Context()
	traceId := pkg.GetTraceId(ctx)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// msg := fmt.Sprintf("Invalid task id: %d, %s", id, err.Error())
		// http.Error(w, msg, http.StatusNoContent)
		slog.Error("Invalid task id",
			slog.String("Error", err.Error()),
			slog.String("TraceId", traceId))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "Invalid task id"})
		return
	}
	slog.Info("Task found", "id", id)
	task, err := h.service.GetTaskById(ctx, id)
	if err != nil {
		slog.Error("Invalid task id",
			slog.String("Error", err.Error()),
			slog.String("TraceId", traceId))
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"Error": "Could not find task with the given id"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *Handler) getAllTasks(c *gin.Context) {
	ctx := c.Request.Context()
	traceId := pkg.GetTraceId(ctx)
	tasks, err := h.service.GetAllTasks(ctx)
	if err != nil {
		// log.Printf("Found error in fetching all tasks.. Error: %s", err.Error())
		// http.Error(w, "Found error in fetching all tasks..", http.StatusNoContent)
		slog.Error("Found error in fetching all tasks",
			slog.String("Error", err.Error()),
			slog.String("TraceId", traceId))
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"Error": "Found error in fetching all tasks"})
		return
	}
	log.Println("All tasks received from DB:", tasks)
	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) UpdateTaskById(c *gin.Context) {
	//validating id
	ctx := c.Request.Context()
	traceId := pkg.GetTraceId(ctx)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		slog.Error("Invalid task id",
			slog.String("Error", err.Error()),
			slog.String("TraceId", traceId))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "Invalid task id"})
		return
	}
	slog.Info("Task found", " id:", id)

	// validating body
	var at = models.AlterTask{}
	err = h.validate.Struct(&at)
	if err != nil {
		slog.Error("Invalid request body, please provide required fields:",
			slog.String("Error:", err.Error()),
			slog.String("TraceId", traceId))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	err = c.ShouldBindJSON(&at)
	if err != nil {
		slog.Error("Invalid request body",
			slog.String("Error", err.Error()),
			slog.String("TraceId", traceId))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}

	//db layer call
	task, err := h.service.UpdateTask(ctx, id, at)

	//validation post db layer call
	if err != nil {
		slog.Error("Error in update task",
			slog.String("Error", "Found error in updating task with id"),
			slog.String("TraceId", traceId))
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"Error": "Found error in updating task"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// func deleteTaskById(c *gin.Context) {
// 	w.Header().Add("Content-Type", "application/json")
// 	id, err := strconv.Atoi(chi.URLParam(r, "id"))
// 	if err != nil {
// 		msg := fmt.Sprintf("Invalid task id: %d, %s", id, err.Error())
// 		http.Error(w, msg, http.StatusNoContent)
// 		return
// 	}
// 	log.Printf("Task found with id:%d", id)
// 	err = db.DeleteTask(r.Context(), id)
// 	if err != nil {
// 		log.Println("Found error in deleting task with id:", id, err)
// 		http.Error(w, "Found error in deleting the task:"+err.Error(), http.StatusNoContent)
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("Deleted task"))
// }

// func updateTaskStatus(c *gin.Context) {
// 	//validating id
// 	w.Header().Add("Content-Type", "application/json")
// 	id, err := strconv.Atoi(chi.URLParam(r, "id"))
// 	if err != nil {
// 		msg := fmt.Sprintf("Invalid task id: %d, %s", id, err.Error())
// 		http.Error(w, msg, http.StatusNoContent)
// 		return
// 	}
// 	log.Printf("Task found with id:%d", id)

// 	//validating body
// 	data, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		log.Printf("Found error in reading request body data byte:%s\n", err.Error())
// 		http.Error(w, "Found error in reading request body data byte", http.StatusBadRequest)
// 		return
// 	}
// 	newTask := &db.NewTask{}
// 	err = json.Unmarshal(data, newTask)
// 	if err != nil {
// 		log.Printf("Invalid request body, cannot convert to json:%s\n", err.Error())
// 		http.Error(w, "Invalid request body, cannot convert to json", http.StatusBadRequest)
// 		return
// 	}
// 	log.Println("Task in middleware ValidateTaskBody:", newTask)

// 	//db layer call
// 	task, err := db.UpdateTaskStatus(r.Context(), id, newTask)

// 	//validation post db layer call
// 	if err != nil {
// 		log.Println("Found error in updating task with id:", id, err)
// 		http.Error(w, "Found error in updating the task", http.StatusInternalServerError)
// 		return
// 	}
// 	data, err = json.Marshal(task)
// 	if err != nil {
// 		log.Println("Found error in json Marshal for task")
// 		http.Error(w, "Found error in json Marshal for task", http.StatusBadRequest)
// 		return
// 	}
// 	//write response
// 	log.Println("Task updated with id:", id)
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(data)
// }
