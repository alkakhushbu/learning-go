package main

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDoubleHandler(t *testing.T) {
	// Name of the test
	// Input
	// output
	tt := [...]struct {
		name           string // Name of the test
		queryParam     string // Query parameter to pass
		expectedStatus int    // Expected HTTP status code
		expectedBody   string // Expected response body
	}{
		{
			name:           "OK",
			queryParam:     "10",
			expectedStatus: http.StatusOK,
			expectedBody:   "20",
		},
		{
			name:           "Fail_MissingValue",
			queryParam:     "", // Missing `v` parameter
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "missing value",
		},
		{
			name:           "Fail_NotANumber",
			queryParam:     "abc", // `v` is not a number
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "not a number: abc",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			// Create a response recorder to capture the response
			// it is an implementation of ResponseWriter
			rec := httptest.NewRecorder()

			// constructing the request
			req := httptest.NewRequest(http.MethodGet, "/double?v="+tc.queryParam, nil)

			// calling the actual handler function
			doubleHandler(rec, req)

			// checking if expected output matches or not
			require.Equal(t, tc.expectedStatus, rec.Code)

			body := rec.Body.String()
			body = strings.TrimSpace(body)

			require.Equal(t, tc.expectedBody, body)

		})

	}
}

//*******************************************************************//
// Benchmark
//*******************************************************************//

func BenchmarkDoubleHandler(b *testing.B) {

}
