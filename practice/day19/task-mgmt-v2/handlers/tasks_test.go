package handlers

import (
	"bytes"
	"context"
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
	task := models.Task{}

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
		code:     http.StatusCreated,
		response: "",
		mockStore: func(m *mockmodels.MockService) {
			m.EXPECT().CreateTask(gomock.Any(), gomock.Eq(newTask)).Return(gomock.Eq(task), nil).Times(1)
		},
	}}

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
			tc.mockStore(mockDb)

			req := httptest.NewRequestWithContext(ctx, http.MethodPost, "/api/v1/tasks", bytes.NewReader(tc.body))
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			require.Equal(t, http.StatusCreated, rec.Code)
		})
	}
}

func TestUpdateTask(t *testing.T) {
	// newTask := models.NewTask{
	// 	Name:      "Middleware for Validating request body",
	// 	Status:    "Registered",
	// 	ManagedBy: "Charles",
	// }
	// task := models.Task{}

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
			m.EXPECT().UpdateTask(gomock.Eq(ctx), gomock.Any(), gomock.Any()).Return(gomock.Any(), nil).Times(1)
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

	router.POST("/api/v1/tasks/:id", h.UpdateTaskById)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockStore(mockDb)

			req := httptest.NewRequestWithContext(ctx, http.MethodPut, "/api/v1/tasks/1", bytes.NewReader(tc.body))
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)
		})
	}
}
