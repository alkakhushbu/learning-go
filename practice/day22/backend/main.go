package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {

	router := gin.Default()
	router.Use(otelgin.Middleware("backend"))
	router.POST("/backend/books", AddBook)

	panic(router.Run(":8089"))
}

func AddBook(c *gin.Context) {
	log.Println("Book added")
	book := struct{ name string }{name: "Go Programming"}
	c.JSON(http.StatusCreated, book)
}
