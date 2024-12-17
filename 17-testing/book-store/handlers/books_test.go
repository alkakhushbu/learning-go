package handlers

import (
	"book-store/middlewares"
	"book-store/models"
	"book-store/models/mockmodels"
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateBook(t *testing.T) {

	// Mock data representing a book returned by the mocked service
	mockBook := models.Book{
		ID:          1,
		Title:       "Go Programming",
		AuthorEmail: "John@doe.com",
		AuthorName:  "John Doe",
		Description: "Learn Go programming with this book",
		Price:       29.99,
		Stock:       1,
		Category:    "Fiction",
	}
	newBook := models.NewBook{
		Title:       "Go Programming",
		AuthorName:  "John Doe",
		AuthorEmail: "John@doe.com",
		Description: "Learn Go programming with this book and enhance your skills.",
		Category:    "Programming",
		Price:       29.99,
		Stock:       10,
	}
	// we need a traceId in the create book, creating a fake one and putting it in the context
	traceId := "fake-trace-id"
	ctx := context.WithValue(context.Background(), middlewares.TraceIdKey, traceId)

	tt := [...]struct {
		name             string
		body             []byte // Body to send to the request
		expectedStatus   int
		expectedResponse string
		//Function to set up the mock behavior // setting expectations
		MockStore func(m *mockmodels.MockService)
	}{
		{
			name: "OK",
			// The JSON payload being sent in the request
			body: []byte(`{
   					 "title": "Go Programming",
   					 "author_name": "John Doe",
   					 "author_email": "John@doe.com",
   					 "description": "Learn Go programming with this book and enhance your skills.",
   					 "category": "Programming",
   					 "price": 29.99,
   					 "stock": 10
				}`),
			expectedStatus:   http.StatusCreated,
			expectedResponse: `{"id":1,"title":"Go Programming","author_Name":"John Doe","author_email":"John@doe.com","description":"Learn Go programming with this book","category":"Fiction","price":29.99}`,
			MockStore: func(m *mockmodels.MockService) {
				// setting the expectations for the mock call
				m.EXPECT().InsertBook(gomock.Eq(ctx), gomock.Eq(newBook)).Return(mockBook, nil).Times(1)
			},
		},
		// input
		// dependencies // it would be setup once for the test function
		// output
	}

	// Creating a new Gin router and setting it to test mode
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Initializing the Gomock controller required for mocking
	ctrl := gomock.NewController(t)

	// NewMockService would give us the implementation of the
	// interface that we can set in handlers struct
	mockDb := mockmodels.NewMockService(ctrl)

	// Creating the handler with the mocked service and validator
	h := Handler{
		service:  mockDb,          // Passing the mocked service
		validate: validator.New(), // Initializing the validator for input validation
	}

	// Registering the CreateBook handler to the POST /create route
	router.POST("/create", h.CreateBook)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// passing the mockDB to MockStore to set expectation for the mocking calls
			tc.MockStore(mockDb)

			//creating the request with context
			req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/create", bytes.NewReader(tc.body))

			// creating a rw implementation
			rec := httptest.NewRecorder()

			// this would call the /create endpoint
			router.ServeHTTP(rec, req)

			require.Equal(t, tc.expectedStatus, rec.Code)
			require.Equal(t, tc.expectedResponse, rec.Body.String())

		})
	}
}
