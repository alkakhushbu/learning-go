package main

import (
	"context"
	"log"
	"task-mgmt/db"
	"task-mgmt/handler"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	connPool, err := pgxpool.NewWithConfig(context.Background(), db.Config())
	if err != nil {
		log.Printf("Error while creating connection to the database:%s\n", err.Error())
		panic("Error while creating connection to the database -_-")
	}
	connection, err := connPool.Acquire(context.Background())
	if err != nil {
		log.Printf("Error while acquiring connection from the database pool:%s\n", err.Error())
		panic("Error while acquiring connection from the database pool -_-")
	}
	defer connection.Release()
	err = connection.Ping(context.Background())
	if err != nil {
		log.Printf("Could not ping database:%s\n", err.Error())
		panic("Could not ping database -_-")
	}
	log.Println("Ping successful......Connected to DB........")

	log.Println("Starting Task Management Service")
	handler.StartHandlerService()
}
