package handlers

import (
	"book-store/middlewares"
	"book-store/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log/slog"
)

type Handler struct {
	c        *models.Conn
	validate *validator.Validate
}

func NewConf(c *models.Conn, validate *validator.Validate) *Handler {
	return &Handler{c: c, validate: validate}
}
func SetupGINRoutes(c *models.Conn) *gin.Engine {
	//r := gin.Default()
	r := gin.New()
	//applying middleware to all the routes
	r.Use(gin.Recovery(), middlewares.Logger())
	slog.Info("Setting up routes")
	h := NewConf(c, validator.New())
	r.GET("/ping", h.Ping)
	r.POST("/create", h.CreateBook)
	r.PATCH("/update/:book_id", h.UpdateBookByID)
	return r
}
