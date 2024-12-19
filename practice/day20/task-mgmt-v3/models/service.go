package models

import "context"

//go generate

// We write go generate to run the below line, it would generate mock implementation of interface
// run "go generate" command from the current directory

// flags
// - source - fileName
// - destination - destination for generated mocks
// - package - package name for mock

//go:generate mockgen -source service.go -destination mockmodels/service_mock.go -package mockmodels
type Service interface {
	Ping(ctx context.Context)
	CreateTask(ctx context.Context, newTask NewTask) (Task, error)
	GetTaskById(ctx context.Context, id int) (Task, error)
	GetAllTasks(ctx context.Context) ([]Task, error)
	UpdateTask(ctx context.Context, id int, alterTask AlterTask) (Task, error)
}
