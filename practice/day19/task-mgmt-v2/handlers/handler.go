package handlers

import (
	"log"
	"task-mgmt-v2/midware"
	"task-mgmt-v2/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	conn     *models.Conn
	validate *validator.Validate
}

func SetupGINRoutes(conn *models.Conn) *gin.Engine {
	log.Println("Inside StartHandlerService function")
	h := Handler{conn: conn, validate: validator.New()}
	route := gin.Default()
	// route.Use(midware.Logger())

	// the midware.Logger() method can also be placed using route.Use(midware.Logger()) method
	api := route.Group("/api/v1", midware.Logger())

	api.POST("/tasks", h.createTask)
	api.GET("/tasks/:id", h.getTaskById)
	api.GET("/tasks", h.getAllTasks)

	// todo: make this PATCH
	api.PUT("/tasks/:id", h.updateTaskById)

	return route
}
