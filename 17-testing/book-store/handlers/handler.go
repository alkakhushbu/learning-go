package handlers

import (
	"book-store/middlewares"
	"book-store/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log/slog"
)

type Handler struct {

	// handlers now accept an interface,
	//which means that I can provide a mock implementation in test
	service  models.Service
	validate *validator.Validate
}

func NewConf(service models.Service, validate *validator.Validate) *Handler {
	return &Handler{service: service, validate: validate}
}

func SetupGINRoutes(service models.Service) *gin.Engine {
	//r := gin.Default()
	r := gin.New()
	//applying middleware to all the routes
	r.Use(gin.Recovery(), middlewares.Logger())
	slog.Info("Setting up routes")

	h := NewConf(service, validator.New())

	r.GET("/ping", h.Ping)
	r.POST("/create", h.CreateBook)
	r.PATCH("/update/:book_id", h.UpdateBookByID)
	return r
}
