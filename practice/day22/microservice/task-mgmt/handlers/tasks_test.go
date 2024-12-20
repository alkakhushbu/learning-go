package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"task-mgmt/midware"
	"task-mgmt/models"
	"task-mgmt/models/mockmodels"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateTask(t *testing.T) {
	newTask := models.NewTask{
		Name:      "Middleware for Validating request body",
		Status:    "Registered",
		ManagedBy: "Charles",
	}
	task := models.Task{
		Id:        1,
		Name:      "Middleware for Validating request body",
		Status:    "Registered",
		ManagedBy: "Charles",
	}
	err := errors.New("Found error in creating task")
	ctx := context.WithValue(context.Background(), midware.TraceId, "Test-Trace-ID")

	// create test table, an array.
	tt := [...]struct {
		name      string
		body      []byte
		code      int
		response  string
		mockStore func(*mockmodels.MockService)
	}{
		{
			name: "Ok",
			body: []byte(`{
				"name": "Middleware for Validating request body",
				"status": "Registered",
				"managedBy": "Charles"}`),
			code:     http.StatusCreated,
			response: `{"id":1,"name":"Middleware for Validating request body","status":"Registered","managedBy":"Charles","startTime":"0001-01-01T00:00:00Z","completionTime":"0001-01-01T00:00:00Z"}`,
			// setting expectations for the mock call
			mockStore: func(m *mockmodels.MockService) {
				m.EXPECT().CreateTask(gomock.Any(), gomock.Eq(newTask)).Return(task, nil).Times(1)
			},
		},
		{
			name:     "JSON Binding failed",
			body:     []byte("Not a JSON"),
			code:     http.StatusBadRequest,
			response: `{"error":"Error in JSON binding"}`,
			mockStore: func(m *mockmodels.MockService) {
				// mocking CreateTask method would cause the test case to fail since the request is never reaching here
				// m.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Return(models.Task{}, err).Times(1)
			},
		},
		{
			name: "JSON Validation failed",
			body: []byte(`{
				"name": "Random ",
				"description": "Valid description"
				}`),
			code:     http.StatusExpectationFailed,
			response: `{"error":"Validation failed, please provide required fields"}`,
			mockStore: func(m *mockmodels.MockService) {
				// mocking CreateTask method would cause the test case to fail since the request is never reaching here
				// m.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Return(models.Task{}, err).Times(1)
			},
		},
		{
			name: "DB Service Down",
			body: []byte(`{
				"name": "Middleware for Validating request body",
				"status": "Registered",
				"managedBy": "Charles"}`),
			code:     http.StatusInternalServerError,
			response: `{"error":"Error in creating new task"}`,
			mockStore: func(m *mockmodels.MockService) {
				m.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Return(models.Task{}, err).Times(1)
			},
		},
	}

	ctrl := gomock.NewController(t)

	// NewMockService would give us the implementation of the
	// interface that we can set in handlers struct
	mockDb := mockmodels.NewMockService(ctrl)

	// Creating the handler with the mocked service and validator
	h := Handler{
		service:  mockDb,          //passing the mocked service
		validate: validator.New(), //initializing the validator for input validation
	}

	// create gin engine
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/tasks", h.CreateTask)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockStore(mockDb)

			// constructing the request
			req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/api/v1/tasks", bytes.NewReader(tc.body))

			// Create a response recorder to capture the response
			// it is an implementation of ResponseWriter
			rec := httptest.NewRecorder()

			//invoke handler serves for the matching request(method, api)
			router.ServeHTTP(rec, req)

			// checking if expected output matches or not
			require.Equal(t, tc.code, rec.Code)
			require.Equal(t, tc.response, rec.Body.String())
		})
	}
}

func TestUpdateTask(t *testing.T) {
	alterTask := models.AlterTask{
		Name:      "Middleware for Validating request body",
		Status:    "Registered",
		ManagedBy: "Charles",
	}
	task := models.Task{
		Name:      "Middleware for Validating request body",
		Status:    "Registered",
		ManagedBy: "Charles",
	}
	err := errors.New("Invalid task id")
	ctx := context.WithValue(context.Background(), midware.TraceId, "Test-Trace-ID")

	// create test table, an array
	tt := [...]struct {
		name      string
		url       string
		body      []byte
		code      int
		response  string
		mockStore func(*mockmodels.MockService)
	}{
		{
			name: "Ok",
			url:  "/api/v1/tasks/1",
			body: []byte(`{
				"name": "Middleware for Validating request body",
				"status": "Registered",
				"managedBy": "Charles"}`),
			code:     http.StatusOK,
			response: "",
			mockStore: func(m *mockmodels.MockService) {
				m.EXPECT().UpdateTask(gomock.Eq(ctx), gomock.Any(), gomock.Eq(alterTask)).Return(task, nil).Times(1)
			},
		},
		{
			name: "Invalid Id non Integer",
			url:  "/api/v1/tasks/abs",
			body: []byte(`{
				"name": "Middleware for Validating request body",
				"status": "Registered",
				"managedBy": "Charles"}`),
			code:     http.StatusBadRequest,
			response: `"Error":"Invalid task id"`,
			mockStore: func(m *mockmodels.MockService) {
				// mocking UpdateTask method would cause the test case to fail since the request is never reaching here
				// m.EXPECT().UpdateTask(gomock.Eq(ctx), gomock.Any(), gomock.Eq(alterTask)).Return(task, nil).Times(1)
			},
		},
		{
			name: "Invalid Id Integer",
			url:  "/api/v1/tasks/123",
			body: []byte(`{
				"name": "Middleware for Validating request body",
				"status": "Registered",
				"managedBy": "Charles"}`),
			code:     http.StatusInternalServerError,
			response: `"Error":"Invalid task id"`,
			mockStore: func(m *mockmodels.MockService) {
				// setting the expectation for mocked service
				m.EXPECT().UpdateTask(gomock.Eq(ctx), gomock.Any(), gomock.Eq(alterTask)).Return(task, err).Times(1)
			},
		},
		{
			name:     "JSON Binding failed",
			url:      "/api/v1/tasks/1",
			body:     []byte("Not a JSON"),
			code:     http.StatusBadRequest,
			response: `{"error":"Error in JSON binding"}`,
			mockStore: func(m *mockmodels.MockService) {
			},
		},
	}

	// create gin engine
	gin.SetMode(gin.TestMode)
	router := gin.New()

	ctrl := gomock.NewController(t)
	mockDb := mockmodels.NewMockService(ctrl)

	h := Handler{
		service:  mockDb,
		validate: validator.New(),
	}

	router.PUT("/api/v1/tasks/:id", h.UpdateTaskById)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockStore(mockDb)

			req := httptest.NewRequestWithContext(ctx, http.MethodPut, tc.url, bytes.NewReader(tc.body))
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			require.Equal(t, tc.code, rec.Code)
		})
	}
}
