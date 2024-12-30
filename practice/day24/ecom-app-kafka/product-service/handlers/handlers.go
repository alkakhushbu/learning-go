package handlers

import (
	"fmt"
	"net/http"
	"os"
	"product-service/internal/products"
	"product-service/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	conf     *products.Conf
	validate *validator.Validate
}

func NewHandler(conf *products.Conf) Handler {
	return Handler{conf: conf, validate: validator.New()}
}

func API(conf *products.Conf) *gin.Engine {
	r := gin.Default()
	mode := os.Getenv("GIN_MODE")
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	prefix := os.Getenv("SERVICE_ENDPOINT_PREFIX")
	if prefix == "" {
		panic("SERVICE_ENDPOINT_PREFIX is not set")
	}

	h := NewHandler(conf)

	r.GET("/ping", healthCheck)
	v1 := r.Group(prefix)
	v1.Use(middleware.Logger())
	{
		v1.POST("/create", h.CreateProduct)
	}

	return r
}

func healthCheck(c *gin.Context) {
	fmt.Println("routine health check : GET /ping endpoint call")
	//JSON serializes the given struct as JSON into the response body. It also sets the Content-Type as "application/json".
	c.JSON(http.StatusOK, gin.H{"status": "ok"})

}
