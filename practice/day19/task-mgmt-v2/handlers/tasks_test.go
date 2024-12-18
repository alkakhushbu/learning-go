package handlers

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"task-mgmt-v2/midware"
	"task-mgmt-v2/models"
	"task-mgmt-v2/models/mockmodels"
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
	// err := errors.New("Found error in creating task")

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
			mockStore: func(m *mockmodels.MockService) {
				m.EXPECT().CreateTask(gomock.Any(), gomock.Eq(newTask)).Return(task, nil).Times(1)
			},
		},
		// {
		// 	name:     "Bad Request",
		// 	body:     []byte(`{"id:":1}`),
		// 	code:     http.StatusBadRequest,
		// 	response: `{"error": "Invalid request body"`,
		// 	mockStore: func(m *mockmodels.MockService) {
		// 		m.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Return(models.Task{}, err).Times(1)
		// 	},
		// },
	}

	ctrl := gomock.NewController(t)
	mockDb := mockmodels.NewMockService(ctrl)

	h := Handler{
		service:  mockDb,
		validate: validator.New(),
	}

	// create gin engine
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/tasks", h.CreateTask)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			log.Println("Before mockstore*****************")
			tc.mockStore(mockDb)
			log.Println("After mockstore*****************")

			req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/api/v1/tasks", bytes.NewReader(tc.body))
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			require.Equal(t, tc.code, rec.Code)
			// require.Equal(t, tc.response, rec.Body.String())
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
	ctx := context.WithValue(context.Background(), midware.TraceId, "Test-Trace-ID")

	// create test table, an array
	tt := [...]struct {
		name      string
		body      []byte
		code      int
		response  string
		mockStore func(*mockmodels.MockService)
	}{{
		name: "Ok",
		body: []byte(`{
				"name": "Middleware for Validating request body",
				"status": "Registered",
				"managedBy": "Charles"}`),
		code:     http.StatusOK,
		response: "",
		mockStore: func(m *mockmodels.MockService) {
			m.EXPECT().UpdateTask(gomock.Eq(ctx), gomock.Any(), gomock.Eq(alterTask)).Return(task, nil).Times(1)
		},
	}}

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

			req := httptest.NewRequestWithContext(ctx, http.MethodPut, "/api/v1/tasks/1", bytes.NewReader(tc.body))
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			require.Equal(t, tc.code, rec.Code)
		})
	}
}
