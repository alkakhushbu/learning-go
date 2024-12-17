package models

import "time"

//todo: move this with db tasks
type Task struct {
	Id             int       `json:"id"`
	Name           string    `json:"name"`
	Status         string    `json:"status"`
	ManagedBy      string    `json:"managedBy"`
	StartTime      time.Time `json:"startTime,omitempty"`
	CompletionTime time.Time `json:"completionTime,omitempty"`
}

type NewTask struct {
	Name           string    `json:"name" validate:"required,min=3,max=100"`
	Status         string    `json:"status" validate:"required,min=3,max=20"`
	ManagedBy      string    `json:"managedBy" validate:"required,min=3,max=100"`
	StartTime      time.Time `json:"startTime,omitempty"` //todo: remove this
	CompletionTime time.Time `json:"completionTime,omitempty"`
}

type AlterTask struct {
	Name           string    `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Status         string    `json:"status,omitempty" validate:"omitempty,min=3,max=20"`
	ManagedBy      string    `json:"managedBy,omitempty" validate:"omitempty,min=3,max=100"`
	StartTime      time.Time `json:"startTime,omitempty"` //todo: remove this
	CompletionTime time.Time `json:"completionTime,omitempty"`
}
