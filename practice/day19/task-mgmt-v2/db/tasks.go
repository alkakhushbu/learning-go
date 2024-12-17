package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"task-mgmt-v2/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Conn struct {
	db *pgxpool.Pool
}

func NewConn() (*Conn, error) {
	const (
		host     = "localhost"
		port     = "5433"
		user     = "postgres"
		password = "postgres"
		dbname   = "postgres"
	)
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// ParseConfig takes the connection string to generate a config
	config, err := pgxpool.ParseConfig(psqlInfo)
	if err != nil {
		return nil, err
	}

	// MinConns is the minimum number of connections kept open by the pool.
	// The pool will not proactively create this many connections, but once this many have been established,
	// it will not close idle connections unless the total number exceeds MaxConns.
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	// MaxConns is the maximum number of connections that can be opened to PostgreSQL.
	// This limit can be used to prevent overwhelming the PostgreSQL server with too many concurrent connections.
	config.MaxConns = 30

	config.HealthCheckPeriod = 5 * time.Minute

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return &Conn{db: db}, nil
}

func (c *Conn) Ping(ctx context.Context) {
	err := c.db.Ping(ctx)
	if err != nil {
		log.Fatal("Could not ping database:%w", err)
	}
	log.Println("Ping to db is successful")
}

func (c *Conn) CreateTask(ctx context.Context, newTask *models.NewTask) (*models.Task, error) {
	insertQuery := `INSERT INTO tasks (name, status, managedby, start_time, completion_time)
					VALUES ($1, $2, $3, $4, $5) returning id,name, status, managedby, start_time, completion_time;`
	// name, status, managedby := "Handlerfunc for Get tasks by id", "In Progress", "Jane"
	startTime := time.Now().UTC()
	completionTime := startTime.AddDate(0, 0, 30)
	var task = &models.Task{}
	err := c.db.QueryRow(ctx, insertQuery,
		newTask.Name, newTask.Status, newTask.ManagedBy, startTime, completionTime).
		Scan(&task.Id, &task.Name, &task.Status, &task.ManagedBy, &task.StartTime, &task.CompletionTime)
	if err != nil {
		log.Println("error in insert query:", err)
		return task, err
	}
	log.Println("Task insert successful with id:", task.Id)
	return task, nil
}

func (c *Conn) GetTaskById(ctx context.Context, id int) (*models.Task, error) {
	selectQuery := `SELECT id, name, status, managedby, start_time, completion_time
	FROM tasks
	WHERE id=$1;`
	var task = &models.Task{}
	err := c.db.QueryRow(ctx, selectQuery, id).
		Scan(&task.Id, &task.Name, &task.Status, &task.ManagedBy, &task.StartTime, &task.CompletionTime)
	if err != nil {
		log.Println("error in scanning query row:", err)
		return nil, err
	}
	log.Println(task)
	return task, nil
}

func (c *Conn) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	selectQuery := `SELECT id, name, status, managedby, start_time, completion_time FROM tasks;`
	tasks := make([]models.Task, 0)
	rows, err := c.db.Query(ctx, selectQuery)
	if err != nil {
		slog.Error("Error", "Found error in task table query", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		task := &models.Task{}
		rows.Scan(&task.Id, &task.Name, &task.Status, &task.ManagedBy, &task.StartTime, &task.CompletionTime)
		slog.Info("Scanning row", slog.Any("Task:", task))
		tasks = append(tasks, *task)
	}
	// log.Println("Scanned all rows in task table:", tasks)
	return tasks, nil
}

// todo: update *models.AlterTask into models.AlterTask
func (c *Conn) UpdateTask(ctx context.Context, id int, alterTask *models.AlterTask) (*models.Task, error) {
	if alterTask.Name == "" && alterTask.Status == "" && alterTask.ManagedBy == "" {
		return nil, fmt.Errorf("all task field(s) empty: name, task status, and managedBy")
	}
	tx, err := c.db.Begin(ctx)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer tx.Rollback(ctx)
	selectQuery := `SELECT id, name, status, managedby, start_time, completion_time FROM tasks
					 WHERE id = $1
					 FOR UPDATE;`
	var dbTask models.Task
	//var dbTask = models.Task{}
	err = tx.QueryRow(ctx, selectQuery, id).Scan(
		&dbTask.Id,
		&dbTask.Name,
		&dbTask.Status,
		&dbTask.ManagedBy,
		&dbTask.StartTime,
		&dbTask.CompletionTime)
	if err != nil {
		return nil, fmt.Errorf("error in select query:%w", err)
	}
	slog.Info("Task", slog.Any("Existing task", dbTask))

	data, err := json.Marshal(alterTask)
	if err != nil {
		return nil, fmt.Errorf("error in marshaling data:%w", err)
	}

	err = json.Unmarshal(data, &dbTask)
	if err != nil {
		return nil, fmt.Errorf("error in unmarshaling data:%w", err)
	}
	slog.Info("Task", slog.Any("Alter db task", dbTask))

	updateQuery := `UPDATE tasks
					SET name=$1, status=$2, managedby=$3 WHERE id = $4
					returning id;`

	// var task = &models.Task{}
	err = tx.QueryRow(ctx, updateQuery, dbTask.Name, dbTask.Status, dbTask.ManagedBy, id).
		Scan(&dbTask.Id) //, &task.Name, &task.Status, &task.ManagedBy, &task.StartTime, &task.CompletionTime)
	if err != nil {
		return nil, fmt.Errorf("error in update query:%w", err)
	}
	// time.Sleep(time.Second * 10)
	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("error in transaction commit:%w", err)
	}
	return &dbTask, nil
}

// func (c *Conn) DeleteTask(ctx context.Context, id int) error {
// 	deleteQuery := `DELETE FROM tasks WHERE id=$1;`
// 	log.Println("Found task with id:", id)
// 	res, err := c.db.Exec(ctx, deleteQuery, id)
// 	if err != nil {
// 		return err
// 	}
// 	if res.RowsAffected() == 0 {
// 		return errors.New("no rows found with id")
// 	}
// 	log.Println("Deleted task with id:", id)
// 	return nil
// }

// func (c *Conn) UpdateTaskStatus(ctx context.Context, id int, newTask *NewTask) (*Task, error) {

// 	if newTask.Status == "" {
// 		return nil, errors.New("empty task field: status")
// 	}
// 	updateQuery := `UPDATE tasks
// 					SET status=$1 WHERE id = $2
// 					returning id, name, status, managedby, start_time, completion_time;`

// 	var task = &Task{}
// 	err := c.db.QueryRow(ctx, updateQuery, newTask.Status, id).
// 		Scan(&task.Id, &task.Name, &task.Status, &task.ManagedBy, &task.StartTime, &task.CompletionTime)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return task, nil
// }
