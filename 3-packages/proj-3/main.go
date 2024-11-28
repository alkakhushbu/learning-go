package main

import (
	"proj-3/database"
	"proj-3/user"
)

func main() {
	database.InitDB("localhost:3306:postgres")
	//database.db = "random" // anyone can change the value of the DB here
	user.AddUser(database.DB)

}
