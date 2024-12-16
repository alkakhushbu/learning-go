package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	//r := gin.Default() //  Default returns an Engine instance with the Logger and Recovery middleware already attached.
	r := gin.New() //start fresh
	r.Use(gin.Logger(), gin.Recovery())

	// JSON Response
	r.GET("/json", sendJson)

	// Route Parameters
	// name is going to be a parameter
	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name") // fetching the param
		c.String(200, "Hello, %s! (Gin)", name)
	})

	// Query Parameters
	r.GET("/welcome", func(c *gin.Context) {
		// this is setting the default value if not set
		firstName := c.DefaultQuery("firstName", "Guest")
		lastName := c.Query("lastName") // fetching the query
		c.String(200, "Hello, %s %s! (Gin)", firstName, lastName)
	})

	// Grouping Routes
	v1 := r.Group("/v1")
	{

		v1.GET("/users", func(c *gin.Context) {
			c.String(200, "Users v1 (Gin)")
		})
		v1.POST("/posts", func(c *gin.Context) {
			c.String(200, "Posts v1 (Gin)")
		})
	}

	r.GET("/custom-error", func(c *gin.Context) {
		err := errors.New("custom error message")
		//AbortWithStatusJSON would not quit the request, we need to return manually
		//AbortWithStatusJSON gives a clear signal that error happened
		// gin.H is a map, useful to send quick responses
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	})

	err := r.Run(":8083")
	if err != nil {
		panic(err)
	}

}

func sendJson(c *gin.Context) {
	u := struct {
		Message string `json:"message"`
	}{
		Message: "Hello, World!",
	}

	// convert the struct to json and send the response as well
	c.JSON(200, u)

}

// middleware for gin
func mid(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("middleware")
		next(c)
	}
}

/*
# GET /
curl http://localhost:8081/

# GET /json
curl http://localhost:8081/json

# GET /user/:name (example with 'name' as 'John')
curl http://localhost:8081/user/John

# GET /welcome with query parameters
curl "http://localhost:8081/welcome?firstName=John&lastName=Doe" // try without firstName

# GET /v1/users
curl http://localhost:8081/v1/users



# GET /custom-error
curl http://localhost:8081/custom-error
*/
