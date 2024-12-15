package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pgxConn *pgxpool.Pool

func GetPGXConn() *pgxpool.Pool {
	return pgxConn
}

func CreateConnection() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	config(ctx)
	log.Println("New database configuration successful......")

	ping(ctx)
	log.Println("Ping to database successful......")
	log.Println("Connection to DB successful........")
}

func config(ctx context.Context) {

	// Your own Database URL
	const DATABASE_URL string = "postgres://postgres:postgres@localhost:5433/postgres?pool_max_conns=10"

	pgxConfig, err := pgxpool.ParseConfig(DATABASE_URL)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	pgxConfig.MaxConns = 4
	pgxConfig.MinConns = 2
	pgxConfig.MaxConnLifetime = time.Hour
	pgxConfig.MaxConnIdleTime = time.Minute * 30
	pgxConfig.HealthCheckPeriod = time.Minute * 5
	pgxConfig.ConnConfig.ConnectTimeout = time.Second * 20
	pgx, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		log.Fatal("Cannot create pgx connection pool to db")
	}
	pgxConn = pgx
}

func ping(ctx context.Context) {
	err := pgxConn.Ping(ctx)
	if err != nil {
		log.Fatal("Could not ping database:%w", err)
	}
}

func CreateTableTask() {
	createTableQuery := `CREATE TABLE IF NOT EXISTS tasks (
				id SERIAL PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				start_time TIMESTAMP,
				completion_time TIMESTAMP,
				status VARCHAR(50),
				managedby VARCHAR(255));`
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	cmdTag, err := pgxConn.Exec(ctx, createTableQuery)
	if err != nil {
		log.Println("Create table task failed:", err)
		return
	}
	log.Println("Create table task successful:", cmdTag)
}
