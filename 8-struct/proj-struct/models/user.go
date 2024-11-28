package models

import (
	"fmt"
	"proj-struct/database"
)

//create a struct for models, and make createUser as a method, and access the conf value using that

func CreateUser(name string, c database.Conf) {
	fmt.Println("adding to db", c)
	fmt.Println("creating the user", name)
}
