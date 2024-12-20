package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {

	router := gin.Default()
	router.Use(otelgin.Middleware("book-store"))
	router.POST("/books", AddBook)
	panic(router.Run(":8086"))
}

func AddBook(c *gin.Context) {
	ctx := c.Request.Context()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8089/books", nil)
}
