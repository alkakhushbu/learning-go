package handlers

import (
	"log"
	"task-mgmt-v2/midware"
	"task-mgmt-v2/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Handler struct {
	service       models.Service
	validate      *validator.Validate
	traceProvider *trace.TracerProvider
}

func SetupGINRoutes(service models.Service, traceProvider *trace.TracerProvider) *gin.Engine {
	log.Println("Inside StartHandlerService function")
	h := Handler{
		service:       service,
		validate:      validator.New(),
		traceProvider: traceProvider,
	}
	
	route := gin.Default()
	// route.Use(midware.Logger())

	// the midware.Logger() method can also be placed using route.Use(midware.Logger()) method
	api := route.Group("/api/v1", midware.Logger())

	api.POST("/tasks", otelgin.Middleware("task-mgmt-v2"), h.CreateTask)
	api.GET("/tasks/:id", h.getTaskById)
	api.GET("/tasks", h.getAllTasks)

	// todo: make this PATCH
	api.PUT("/tasks/:id", h.UpdateTaskById)

	return route
}
