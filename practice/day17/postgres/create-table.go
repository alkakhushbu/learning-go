package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

/*
q1. Create docker compose file to run postgres container
    Connect to postgres using pgx
    Create movies table using Go program
    Insert two records within transaction to movies table
    Update one record using optimistic locking
*/

func CreateConnection() (*pgxpool.Pool, error) {
	const (
		host     = "localhost"
		port     = "5434"
		user     = "postgres"
		password = "postgres"
		dbname   = "postgres"
	)
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		host, port, user, password, dbname)
	config, err := pgxpool.ParseConfig(psqlInfo)

	if err != nil {
		return nil, err
	}
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	config.MaxConns = 30

	config.HealthCheckPeriod = 5 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Ping(db *pgxpool.Pool) {
	// pinging the connection if it is alive or not
	err := db.Ping(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println("Ping successful")
}

func CreateTable(pgx *pgxpool.Pool) {
	query := `CREATE TABLE IF NOT EXISTS movies (
			movie_id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			release_year INT,
			genre VARCHAR(100),
			director VARCHAR(255),
			rating DECIMAL(2, 1),
			version INT NOT NULL DEFAULT 1
			);`
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	status, err := pgx.Exec(ctx, query)
	if err != nil {
		fmt.Println("Error in creating table movies")
		return
	}
	fmt.Println("Table movies created successfully:", status)
}

func main() {
	pgx, err := CreateConnection()
	if err != nil {
		msg := fmt.Sprintf("Cannot create connection:%w", err)
		panic(msg)
	}
	Ping(pgx)
	CreateTable(pgx)
}
