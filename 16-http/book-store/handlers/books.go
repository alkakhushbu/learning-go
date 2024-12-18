package handlers

import (
	"book-store/models"
	"book-store/pkg/ctxmanage"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
)

func (h *Handler) Ping(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}

// CreateBook handles the HTTP request to create a new book record.
// It performs validation on the request, inserts the new book into the database,
// and returns the created book or an error if something goes wrong.
func (h *Handler) CreateBook(c *gin.Context) {
	// Extracting the trace ID for logging (useful for tracking requests).
	traceId := ctxmanage.GetTraceID(c.Request.Context())

	// Declare a variable to hold the incoming request data for creating a book.
	var nb models.NewBook

	// Parse the JSON request body into the variable `nb`.
	err := c.ShouldBindJSON(&nb)
	if err != nil {
		// Log the error and respond with a 400 (Bad Request) status code if the JSON is invalid.
		slog.Error(
			"Invalid JSON request",            // Error description.
			slog.String("TRACE ID", traceId),  // Attach trace ID for tracking purposes.
			slog.String("Error", err.Error()), // Log the specific error message.
		)
		c.JSON(400, gin.H{"error": "invalid request body"}) // Respond to the client.
		return                                              // Exit the function to stop further processing.
	}

	// Validate the parsed request data using a validation library or custom validation rules.
	err = h.validate.Struct(nb)
	if err != nil {
		// Log the error and respond with a 400 (Bad Request) status code if validation fails.
		slog.Error(
			"Validation failed",
			slog.String("TRACE ID", traceId),
			slog.String("Error", err.Error()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "please provide values in correct format"})
		return
	}

	// If validation succeeds, insert the book into the database.
	book, err := h.c.InsertBook(c.Request.Context(), nb)
	if err != nil {
		// Log the error and respond with a 500 (Internal Server Error) if the insertion fails.
		slog.Error(
			"Error inserting book",
			slog.String("TRACE ID", traceId),
			slog.String("Error", err.Error()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "problem inserting book"})
		return
	}

	// Log the successful insertion and respond with a 201 (Created) status code.
	slog.Info("Book inserted", slog.String("TRACE ID", traceId))
	c.JSON(http.StatusCreated, book)
}

// UpdateBookByID handles the HTTP request to update an existing book record by its ID.
// It validates the request, updates the book in the database, and returns the updated book or an error.
func (h *Handler) UpdateBookByID(c *gin.Context) {
	// Extracting the trace ID for logging (similar to CreateBook).
	traceId := ctxmanage.GetTraceID(c)

	// Extract the `book_id` parameter from the URL.
	idParam := c.Param("book_id")

	// Convert the `book_id` from a string to an integer.
	id, err := strconv.Atoi(idParam)
	if err != nil {
		// Log the error and respond with a 400 (Bad Request) status code if the `book_id` is invalid.
		slog.Error(
			"Invalid book id",
			slog.String("TRACE ID", traceId),
			slog.String("Error", err.Error()),
		)
		c.AbortWithStatusJSON(400, gin.H{"error": "invalid book id"}) // Abort request processing.
		return
	}

	// Declare a variable to hold the incoming request data for updating the book.
	var updateBook models.UpdateBook

	// Parse the JSON request body into the variable `updateBook`.
	err = c.ShouldBindJSON(&updateBook)
	if err != nil {
		// Log the error and respond with a 400 (Bad Request) status code if the JSON is invalid.
		slog.Error(
			"Invalid JSON request",
			slog.String("TRACE ID", traceId),
			slog.String("Error", err.Error()),
		)
		c.AbortWithStatusJSON(400, gin.H{"error": "invalid request body"})
		return
	}

	// Update the book in the database with the given ID and updated data.
	updatedBook, err := h.c.Update(c.Request.Context(), id, updateBook)
	if err != nil {
		// Log the error and respond with a 500 (Internal Server Error) if the update fails.
		slog.Error(
			"Error updating book",
			slog.String("TRACE ID", traceId),
			slog.String("Error", err.Error()),
		)
		c.AbortWithStatusJSON(500, gin.H{"error": "problem updating book"})
		return
	}

	// If the update succeeds, respond with the updated book and a 200 (OK) status code.
	c.JSON(http.StatusOK, updatedBook)
}
