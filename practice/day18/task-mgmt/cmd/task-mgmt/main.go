package main

import (
	"log"
	"task-mgmt/api/handler"
	"task-mgmt/db"
)

func main() {
	log.Println("Starting to connect to DB.....")
	db.CreateConnection()

	log.Println("Creating table task.....")
	db.CreateTableTask()

	log.Println("Starting Task Management Service")
	handler.StartHandlerService()

}
