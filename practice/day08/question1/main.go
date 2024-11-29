package main

import (
	"fmt"
	"question1/stores"
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
	user1 := stores.User{Name: "alka"}
	user2 := stores.User{Name: "John"}
	// users := []stores.User{
	// 	{Name: "alka"},
	// 	{Name: "alka"},
	// }
	var database stores.Database
	database = mysql.NewConnection()
	database.Create(user1)
	database.Create(user2)
	database.Update(1, "khushbu")
	database.FetchAll()
	database.FetchUser(3)

	database = postgres.NewConnection()
	database.Create(user1)
	database.Create(user2)
	database.Update(1, "khushbu")
	database.FetchAll()
	database.FetchUser(3)
}
