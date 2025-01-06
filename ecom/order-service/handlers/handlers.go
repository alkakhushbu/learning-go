package handlers

import (
	"order-service/internal/auth"
	"order-service/internal/carts"
	"order-service/internal/orders"
	"order-service/internal/stores/kafka"
	"order-service/middleware"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	consulapi "github.com/hashicorp/consul/api"
)

type Handler struct {
	client    *consulapi.Client
	dbConf    *orders.Conf
	kafkaConf *kafka.Conf
	c         *carts.Conf
	validate  *validator.Validate
}

func NewHandler(client *consulapi.Client, dbConf *orders.Conf,
	kafkaConf *kafka.Conf, cartConf *carts.Conf) *Handler {
	return &Handler{
		client:    client,
		dbConf:    dbConf,
		kafkaConf: kafkaConf,
		c:         cartConf,
		validate:  validator.New(),
	}
}

func API(endpointPrefix string, k *auth.Keys,
	client *consulapi.Client, dbConf *orders.Conf,
	kafkaConf *kafka.Conf, cartConf *carts.Conf) *gin.Engine {
	r := gin.New()
	mode := os.Getenv("GIN_MODE")
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	m, err := middleware.NewMid(k)
	if err != nil {
		panic(err)
	}

	h := NewHandler(client, dbConf, kafkaConf, cartConf)
	r.Use(middleware.Logger(), gin.Recovery())

	r.GET("/ping", HealthCheck)
	v1 := r.Group(endpointPrefix)
	{
		v1.POST("/webhook", h.Webhook)
		v1.Use(m.Authentication())
		v1.POST("/checkout/:productID", h.ProductCheckout)

		v1.POST("/carts/checkout", h.CartCheckout)
		v1.POST("/carts/add", h.AddToCart)
		v1.POST("/carts/remove", h.RemoveFromCart)
		v1.GET("/carts", h.GetAllCartItems)
		v1.GET("/ping", HealthCheck)

	}

	return r
}

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
