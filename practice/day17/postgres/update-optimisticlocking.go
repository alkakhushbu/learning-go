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

type Movie struct {
	Title       string
	ReleaseYear int
	Genre       string
	Director    string
	Rating      float32
	Version     int
}

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

func Update(pgx *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	selectQuery := `SELECT movie_id, version FROM movies WHERE title = $1 for update;`
	updateQuery := `UPDATE movies 
					SET rating = $1, version = version + 1 
					WHERE title = $2 and version = $3 
					RETURNING movie_id, version;`
	rating := 9.1
	title := "Inception"
	var movie_id, version int
	err := pgx.QueryRow(ctx, selectQuery, title).Scan(&movie_id, &version)
	if err != nil {
		return fmt.Errorf("error in select query: %w", err)
	}
	// time.Sleep(time.Second * 10)
	fmt.Printf("movieid: %d, version: %d\n", movie_id, version)
	err = pgx.QueryRow(ctx, updateQuery, rating, title, version).Scan(&movie_id, &version)
	if err != nil {
		return fmt.Errorf("error in update query: %w", err)
	}
	selectQuery = `SELECT movie_id, version FROM movies WHERE title = $1 for update;`
	err = pgx.QueryRow(ctx, selectQuery, title).Scan(&movie_id, &version)
	if err != nil {
		return fmt.Errorf("error in update query: %w", err)
	}
	fmt.Println("Executed successful update")
	return nil
}
func main() {
	pgx, err := CreateConnection()
	if err != nil {
		msg := fmt.Sprintf("Cannot create connection:%w", err)
		panic(msg)
	}
	err = Update(pgx)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("success in update records through optimistic locking")
}
