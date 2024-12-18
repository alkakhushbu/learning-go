// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -source service.go -destination mockmodels/service_mock.go -package mockmodels
//

// Package mockmodels is a generated GoMock package.
package mockmodels

import (
	context "context"
	reflect "reflect"
	models "task-mgmt-v2/models"

	gomock "go.uber.org/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
	isgomock struct{}
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// CreateTask mocks base method.
func (m *MockService) CreateTask(ctx context.Context, newTask models.NewTask) (models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTask", ctx, newTask)
	ret0, _ := ret[0].(models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTask indicates an expected call of CreateTask.
func (mr *MockServiceMockRecorder) CreateTask(ctx, newTask any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTask", reflect.TypeOf((*MockService)(nil).CreateTask), ctx, newTask)
}

// GetAllTasks mocks base method.
func (m *MockService) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllTasks", ctx)
	ret0, _ := ret[0].([]models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllTasks indicates an expected call of GetAllTasks.
func (mr *MockServiceMockRecorder) GetAllTasks(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllTasks", reflect.TypeOf((*MockService)(nil).GetAllTasks), ctx)
}

// GetTaskById mocks base method.
func (m *MockService) GetTaskById(ctx context.Context, id int) (models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTaskById", ctx, id)
	ret0, _ := ret[0].(models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTaskById indicates an expected call of GetTaskById.
func (mr *MockServiceMockRecorder) GetTaskById(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTaskById", reflect.TypeOf((*MockService)(nil).GetTaskById), ctx, id)
}

// Ping mocks base method.
func (m *MockService) Ping(ctx context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Ping", ctx)
}

// Ping indicates an expected call of Ping.
func (mr *MockServiceMockRecorder) Ping(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockService)(nil).Ping), ctx)
}

// UpdateTask mocks base method.
func (m *MockService) UpdateTask(ctx context.Context, id int, alterTask models.AlterTask) (models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTask", ctx, id, alterTask)
	ret0, _ := ret[0].(models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTask indicates an expected call of UpdateTask.
func (mr *MockServiceMockRecorder) UpdateTask(ctx, id, alterTask any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTask", reflect.TypeOf((*MockService)(nil).UpdateTask), ctx, id, alterTask)
}
