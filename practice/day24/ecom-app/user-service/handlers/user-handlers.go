package handlers

import (
	"log/slog"
	"net/http"
	"user-service/internal/users"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Signup(c *gin.Context) {
	var newUser users.NewUser
	err := c.ShouldBindJSON(&newUser)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Error in json binding"})
		return
	}
	err = h.validate.Struct(newUser)
	if err != nil {
		// Log the error and respond with a 400 (Bad Request) status code if validation fails.
		slog.Error(
			"Validation failed",
			slog.String("Error", err.Error()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "please provide values in correct format"})
		return
	}

	user, err := h.conf.InsertUser(c.Request.Context(), newUser)
	if err != nil {
		// Log the error and respond with a 400 (Bad Request) status code if validation fails.
		slog.Error(
			"Error in inserting user into database",
			slog.String("Error", err.Error()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error in inserting user into database"})
		return
	}

	c.JSON(http.StatusOK, user)
}
