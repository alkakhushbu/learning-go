package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
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

func ExecuteWithTx(pgx *pgxpool.Pool, fn func(context.Context, pgx.Tx) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	tx, err := pgx.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			fmt.Println("error in rollback:", err.Error())
			return
		}

	}()
	err = fn(ctx, tx)
	if err != nil {
		return fmt.Errorf("error in executing transaction queries: %w", err)
	}
	fmt.Println("Query executed successfully")
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error commiting the transaction: %w", err)
	}
	return nil
}

func InsertManyWithTx(ctx context.Context, tx pgx.Tx) error {
	insertQuery := `INSERT INTO movies (title, release_year, genre, director, rating)
				VALUES ($1, $2, $3, $4, $5)`
	m := Movie{Title: "The Godfather", ReleaseYear: 1972, Genre: "Crime", Director: "Francis Ford Coppola", Rating: 9.2}
	_, err := tx.Exec(ctx, insertQuery, m.Title, m.ReleaseYear, m.Genre, m.Director, m.Rating)

	if err != nil {
		return fmt.Errorf("error inserting record: %w", err)
	}
	// time.Sleep(time.Second * 10)
	// return fmt.Errorf("some error happened")

	m = Movie{Title: "Parasite", ReleaseYear: 2019, Genre: "Thriller", Director: "Bong Joon Ho", Rating: 8.6}
	_, err = tx.Exec(ctx, insertQuery, m.Title, m.ReleaseYear, m.Genre, m.Director, m.Rating)

	if err != nil {
		return fmt.Errorf("error inserting record: %w", err)
	}
	return nil
}

func main() {
	pgx, err := CreateConnection()
	if err != nil {
		msg := fmt.Sprintf("Cannot create connection:%w", err)
		panic(msg)
	}
	err = ExecuteWithTx(pgx, InsertManyWithTx)
	if err != nil {
		fmt.Println("err inserting many records:", err)
		return
	}
	fmt.Println("success in inserting many records through transaction")
}
