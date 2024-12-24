package handlers

import (
	"os"
	"user-service/internal/users"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type Handler struct {
	conf     *users.Conf
	validate *validator.Validate
}

func NewConf(c *users.Conf, validate *validator.Validate) *Handler {
	return &Handler{conf: c, validate: validate}
}

func API(c *users.Conf) *gin.Engine {
	r := gin.New()
	mode := os.Getenv("GIN_MODE")
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	prefix := os.Getenv("SERVICE_ENDPOINT_PREFIX")
	if prefix == "" {
		panic("SERVICE_ENDPOINT_PREFIX is not set")
	}

	//setUpHandler and convert all handlerfunctions of user-handlers into methods of handlers
	h := NewConf(c, validator.New())

	v1 := r.Group(prefix)
	{
		v1.Use(gin.Logger(), gin.Recovery())
		v1.POST("/signup", h.Signup)
	}
	return r
}
