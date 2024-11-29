package main

import (
	"fmt"
	"question1/stores"
	"question1/stores/models"
	"question1/stores/mysql"
	"question1/stores/postgres"
)

/*
Note: day 7, question number 3 extension
q1. Create a map in postgres.go and mysql.go to store user details in the memory

	Change the signature of methods (Create(user) bool, Update(name) bool, Delete(id) bool)
	Note:- true would indicate operation has been successful.
	Add two new methods FetchAll, FetchUser(id) bool
	This time write concrete functionalities not simple print statements*/

func main() {
	fmt.Println("Hello World!")
	user1 := models.User{Name: "alka"}
	user2 := models.User{Name: "John"}

	// var database stores.Database
	mysql := mysql.NewConnection()
	store := stores.NewStore(mysql)
	database := store.Database
	database.Create(user1)
	database.Create(user2)
	database.Update(1, "khushbu")
	database.Delete(3)
	database.FetchAll()
	database.FetchUser(3)

	postgres := postgres.NewConnection()
	store = stores.NewStore(postgres)
	database = store.Database
	database.Create(user1)
	database.Create(user2)
	database.Update(1, "khushbu")
	database.Delete(3)
	database.FetchAll()
	database.FetchUser(3)
}
