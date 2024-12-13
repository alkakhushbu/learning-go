package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

//orm
//https://gorm.io/docs/query.html

// driver
//https://github.com/jackc/pgx

// whatever moduel we do go get is stored inside our gopath
// go env GOPATH

// go get moduleName (to get an external lib)
// github.com/jackc/pgx/v5 (don't forget to include the version number if there is a major version in the module name)

// automatically resolves all the dependecies required for the project
// go mod tidy // first command to run when importing any project

func CreateConnection() (*pgxpool.Pool, error) {

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
	return db, nil
}

func Ping(db *pgxpool.Pool) {
	// pinging the connection if it is alive or not
	err := db.Ping(context.Background())
	if err != nil {
		panic(err)
	}
}

// Three methods to execute queries on the database
// Exec -> when query does not return anything
// QueryRow -> returns exactly one row
// Query -> returns many rows

func main() {
	db, err := CreateConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	Ping(db)

	CreateTable(db)
	u := User{name: "alka", email: "alka@gmail.com", age: 22}
	// InsertRecord(db, u)

	u.name = "khushbu"
	UpdateRecord(db, u, 1)
}

func CreateTable(db *pgxpool.Pool) {
	query := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100),
        email VARCHAR(100) UNIQUE NOT NULL,
        age INT
    );`
	res, err := db.Exec(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("rows affected: %d\n", res.RowsAffected())
}

// Create two function one to insert one record and another to update the record

type User struct {
	name  string
	email string
	age   int
}

func InsertRecord(db *pgxpool.Pool, u User) {
	query := `Insert into users(name, email, age) values($1, $2, $3);`

	res, err := db.Exec(context.Background(), query, u.name, u.email, u.age)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("rows affected: %d\n", res.RowsAffected())
}

func UpdateRecord(db *pgxpool.Pool, u User, id int) {
	query := `update users set name=$1, email=$2, age=$3 where id=$4;`
	res, err := db.Exec(context.Background(), query, u.name, u.email, u.age, id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("rows affected: %d\n", res.RowsAffected())
}
