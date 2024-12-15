package db

import (
	"context"
	"errors"
	"log"
	"time"
)

func CreateTask(task *Task) (int, error) {
	insertQuery := `INSERT INTO tasks (name, status, managedby, start_time, completion_time)
					VALUES ($1, $2, $3, $4, $5) returning id;`
	// name, status, managedby := "Handlerfunc for Get tasks by id", "In Progress", "Jane"
	startTime := time.Now().UTC()
	completionTime := startTime.AddDate(0, 0, 30)
	var id int
	err := pgxConn.QueryRow(context.Background(), insertQuery,
		task.Name, task.Status, task.ManagedBy, startTime, completionTime).Scan(&id)
	if err != nil {
		log.Println("error in insert query:", err)
		return 0, err
	}
	log.Println("Task insert successful with id:", id)
	return id, nil
}

func GetTaskById(id int) (*Task, error) {
	selectQuery := `SELECT id, name, status, managedby, start_time, completion_time
	FROM tasks
	WHERE id=$1;`
	var task = &Task{}
	err := pgxConn.QueryRow(context.Background(), selectQuery, id).
		Scan(&task.Id, &task.Name, &task.Status, &task.ManagedBy, &task.StartTime, &task.CompletionTime)
	if err != nil {
		log.Println("error in scanning query row:", err)
		return nil, err
	}
	log.Println(task)
	return task, nil
}

func GetAllTasks() ([]Task, error) {
	selectQuery := `SELECT id, name, status, managedby, start_time, completion_time FROM tasks;`
	tasks := make([]Task, 0)
	rows, err := pgxConn.Query(context.Background(), selectQuery)
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

func UpdateTask(id int, task *Task) (*Task, error) {
	if task.Name == "" || task.Status == "" || task.ManagedBy == "" {
		return nil, errors.New("empty task field(s) name,task status, or managedBy")
	}
	updateQuery := `UPDATE tasks
					SET name=$1, status=$2, managedby=$3 WHERE id = $4
					returning id, name, status, managedby, start_time, completion_time;`

	var newTask = &Task{}
	err := pgxConn.QueryRow(context.Background(), updateQuery, task.Name, task.Status, task.ManagedBy, id).
		Scan(&newTask.Id, &newTask.Name, &newTask.Status, &newTask.ManagedBy, &newTask.StartTime, &newTask.CompletionTime)
	if err != nil {
		return nil, err
	}
	return newTask, nil
}

func DeleteTask(id int) error {
	deleteQuery := `DELETE FROM tasks WHERE id=$1;`
	log.Println("Found task with id:", id)
	res, err := pgxConn.Exec(context.Background(), deleteQuery, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("no rows found with id")
	}
	log.Println("Deleted task with id:", id)
	return nil
}

func UpdateTaskStatus(id int, task *Task) (*Task, error) {

	if task.Status == "" {
		return nil, errors.New("empty task field: status")
	}
	updateQuery := `UPDATE tasks
					SET status=$1 WHERE id = $2
					returning id, name, status, managedby, start_time, completion_time;`

	err := pgxConn.QueryRow(context.Background(), updateQuery, task.Status, id).
		Scan(&task.Id, &task.Name, &task.Status, &task.ManagedBy, &task.StartTime, &task.CompletionTime)
	if err != nil {
		return nil, err
	}
	return &task, nil
}
