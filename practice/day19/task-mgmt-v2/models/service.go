package models

import "context"

type Service interface {
	Ping(context.Context)
	CreateTask(context.Context, *NewTask) (*Task, error)
	GetTaskById(context.Context, int) (*Task, error)
	GetAllTasks(context.Context) ([]Task, error)
	UpdateTask(context.Context, int, *AlterTask) (*Task, error)
}
