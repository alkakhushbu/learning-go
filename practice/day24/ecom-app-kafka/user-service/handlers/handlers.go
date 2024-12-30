package handlers

import (
	"fmt"
	"net/http"
	"os"
	"user-service/internal/auth"
	"user-service/internal/stores/kafka"
	"user-service/internal/users"
	"user-service/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	u        *users.Conf
	validate *validator.Validate
	k        *kafka.Conf
	a        *auth.Keys
}

func NewHandler(u *users.Conf, k *kafka.Conf, a *auth.Keys) *Handler {

	return &Handler{
		u:        u,
		k:        k,
		validate: validator.New(),
		a:        a,
	}
}

func API(u *users.Conf, k *kafka.Conf, a *auth.Keys) *gin.Engine {
	r := gin.New()
	mode := os.Getenv("GIN_MODE")
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	h := NewHandler(u, k, a)

	prefix := os.Getenv("SERVICE_ENDPOINT_PREFIX")
	if prefix == "" {
		panic("SERVICE_ENDPOINT_PREFIX is not set")
	}
	r.GET("/ping", healthCheck)
	v1 := r.Group(prefix)
	v1.Use(middleware.Logger())
	{
		v1.Use(gin.Logger(), gin.Recovery())
		v1.POST("/signup", h.Signup)
		v1.POST("/login", h.Login)

		// this middleware would be applied to the handler functions which are after it
		// it would not apply to the previous one
		v1.Use(middleware.Authentication(a))
		v1.GET("/check", h.AuthCheck)
	}

	return r
}

func healthCheck(c *gin.Context) {
	fmt.Println("routine health check : GET /ping endpoint call")
	//JSON serializes the given struct as JSON into the response body. It also sets the Content-Type as "application/json".
	c.JSON(http.StatusOK, gin.H{"status": "ok"})

}
