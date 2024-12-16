package db

import (
	"context"
	"errors"
	"log"
	"time"
)

func CreateTask(ctx context.Context, newTask *NewTask) (*Task, error) {
	insertQuery := `INSERT INTO tasks (name, status, managedby, start_time, completion_time)
					VALUES ($1, $2, $3, $4, $5) returning id,name, status, managedby, start_time, completion_time;`
	// name, status, managedby := "Handlerfunc for Get tasks by id", "In Progress", "Jane"
	startTime := time.Now().UTC()
	completionTime := startTime.AddDate(0, 0, 30)
	var task = &Task{}
	err := pgxConn.QueryRow(ctx, insertQuery,
		newTask.Name, newTask.Status, newTask.ManagedBy, startTime, completionTime).
		Scan(&task.Id, &task.Name, &task.Status, &task.ManagedBy, &task.StartTime, &task.CompletionTime)
	if err != nil {
		log.Println("error in insert query:", err)
		return task, err
	}
	log.Println("Task insert successful with id:", task.Id)
	return task, nil
}

func GetTaskById(ctx context.Context, id int) (*Task, error) {
	selectQuery := `SELECT id, name, status, managedby, start_time, completion_time
	FROM tasks
	WHERE id=$1;`
	var task = &Task{}
	err := pgxConn.QueryRow(ctx, selectQuery, id).
		Scan(&task.Id, &task.Name, &task.Status, &task.ManagedBy, &task.StartTime, &task.CompletionTime)
	if err != nil {
		log.Println("error in scanning query row:", err)
		return nil, err
	}
	log.Println(task)
	return task, nil
}

func GetAllTasks(ctx context.Context) ([]Task, error) {
	selectQuery := `SELECT id, name, status, managedby, start_time, completion_time FROM tasks;`
	tasks := make([]Task, 0)
	rows, err := pgxConn.Query(ctx, selectQuery)
	if err != nil {
		log.Println("Found error in task table query", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		task := &Task{}
		rows.Scan(&task.Id, &task.Name, &task.Status, &task.ManagedBy, &task.StartTime, &task.CompletionTime)
		log.Printf("Scanning row: %v", task)
		tasks = append(tasks, *task)
	}
	// log.Println("Scanned all rows in task table:", tasks)
	return tasks, nil
}

func UpdateTask(ctx context.Context, id int, newTask *NewTask) (*Task, error) {
	if newTask.Name == "" || newTask.Status == "" || newTask.ManagedBy == "" {
		return nil, errors.New("empty task field(s) name,task status, or managedBy")
	}
	updateQuery := `UPDATE tasks
					SET name=$1, status=$2, managedby=$3 WHERE id = $4
					returning id, name, status, managedby, start_time, completion_time;`

	var task = &Task{}
	err := pgxConn.QueryRow(ctx, updateQuery, newTask.Name, newTask.Status, newTask.ManagedBy, id).
		Scan(&task.Id, &task.Name, &task.Status, &task.ManagedBy, &task.StartTime, &task.CompletionTime)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func DeleteTask(ctx context.Context, id int) error {
	deleteQuery := `DELETE FROM tasks WHERE id=$1;`
	log.Println("Found task with id:", id)
	res, err := pgxConn.Exec(ctx, deleteQuery, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("no rows found with id")
	}
	log.Println("Deleted task with id:", id)
	return nil
}

func UpdateTaskStatus(ctx context.Context, id int, newTask *NewTask) (*Task, error) {

	if newTask.Status == "" {
		return nil, errors.New("empty task field: status")
	}
	updateQuery := `UPDATE tasks
					SET status=$1 WHERE id = $2
					returning id, name, status, managedby, start_time, completion_time;`

	var task = &Task{}
	err := pgxConn.QueryRow(ctx, updateQuery, newTask.Status, id).
		Scan(&task.Id, &task.Name, &task.Status, &task.ManagedBy, &task.StartTime, &task.CompletionTime)
	if err != nil {
		return nil, err
	}
	return task, nil
}
