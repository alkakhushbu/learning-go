package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"log"
)

const (
	driverName = "pgx" // Database driver

	migrationsDir = "migration" // Directory where the migration files are stored
)

func main() {
	const (
		host     = "localhost"
		port     = "5433"
		user     = "postgres"
		password = "postgres"
		dbname   = "postgres"
	)

	//sql.Open(psqlInfo)
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	// Open a connection to the database
	db, err := sql.Open(driverName, psqlInfo)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Set the dialect for Goose (PostgreSQL in this case)
	goose.SetDialect(driverName)

	//// Apply all pending migrations in the directory
	//err = goose.Up(db, migrationsDir)
	//if err != nil {
	//	log.Fatalf("Failed to apply migrations: %v", err)
	//}
	//goose.Down(db, migrationsDir)
	log.Println("Database migrations applied successfully!")
	err = goose.Up(db, migrationsDir)

	//err = goose.UpTo(db, migrationsDir, 1)
	if err != nil {

		log.Fatalf("Failed to apply migrations: %v", err)
	}

}
