package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/set-cookie", func(c *gin.Context) {
		token := "example-token-value" // Your token value here
		SetCookie(c, token)
		c.JSON(200, gin.H{
			"message": "Cookie has been set",
		})
	})

	r.Run(":8080") // Listen and serve on localhost:8080
}

func SetCookie(c *gin.Context, token string) {
	// Setting an HTTP-only cookie using Gin
	c.SetCookie(
		"token", // Name of the cookie
		token,   // Value of the cookie
		3600,    // MaxAge in seconds (1 hour here)
		"/",     // Path (accessible site-wide)
		"",      // Domain (empty means current domain)
		true,    // Secure (send cookie only over HTTPS)
		true,    // HttpOnly (not accessible via JavaScript)
	)
}

func GetCookie(c *gin.Context) {
	// Try to fetch the "token" cookie
	token, err := c.Cookie("token")
	if err != nil {
		// Handle the case where the cookie does not exist or couldn't be fetched
		c.JSON(400, gin.H{
			"error": "Cookie not found",
		})
		return
	}

	// If the cookie is found, do something with its value
	c.JSON(200, gin.H{
		"token":   token,
		"message": "Cookie fetched successfully",
	})
}
