package main

import (
	"fmt"
	"net/http"
)

// http://localhost:8086/search?user_name=a
var s []string = []string{"a", "b", "c"}

func main() {
	// Register the "/search" route to handle HTTP requests using the Search function
	http.HandleFunc("/search", Search)

	// Start the HTTP server on port 8086 and handle incoming requests
	http.ListenAndServe(":8086", nil)
}

// Search is the HTTP handler function for the "/search" endpoint.
// It checks if a user_name passed as a query parameter exists in the slice
func Search(w http.ResponseWriter, r *http.Request) {
	// Log the full query parameters for debugging purposes
	fmt.Println(r.URL.Query())

	// Extract the "user_name" query parameter from the request URL
	userName := r.URL.Query().Get("user_name")

	for _, name := range s {
		if name == userName {
			// If the user is found, respond with an HTTP 200 status code (OK)
			// Write a personalized greeting in the HTTP response body
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello " + userName))
			return
		}
	}

	// If the user is not found in the database, return an HTTP 404 status code (Not Found)
	// Respond with an appropriate error message
	http.Error(w, "user not found in db", http.StatusNotFound)
	return
}
