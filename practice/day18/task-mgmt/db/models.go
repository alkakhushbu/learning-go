package db

import "time"

type Task struct {
	Id             int       `json:"id"`
	Name           string    `json:"name"`
	Status         string    `json:"status"`
	ManagedBy      string    `json:"managedBy"`
	StartTime      time.Time `json:"startTime,omitempty"`
	CompletionTime time.Time `json:"completionTime,omitempty"`
}

type NewTask struct {
	Name           string    `json:"name"`
	Status         string    `json:"status"`
	ManagedBy      string    `json:"managedBy"`
	StartTime      time.Time `json:"startTime,omitempty"`
	CompletionTime time.Time `json:"completionTime,omitempty"`
}
