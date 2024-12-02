package main

import (
	"app/stores"
	"app/stores/models"
	"app/stores/mysql"
	"app/stores/postgress"
	"fmt"
)

// Refer question q3 from day 7 and q1 from day8

func main() {
	//Call postgres and mysql package methods using interface variable
	//var i stores.DataBase

	u1 := models.User{Id: 1, Name: "sandra"}
	u2 := models.User{Id: 2, Name: "komal"}

	m := mysql.NewConn()
	i := stores.NewStore(m)

	fmt.Printf("---------------- Printing type of interface %T ---------------------------------\n", i)

	u, ok := i.Create(u1)
	if !ok {
		fmt.Println("failed to create user")
		return
	}
	fmt.Println("user created in mysql", u)

	allUsers, ok := i.FetchAll()
	if !ok {
		fmt.Println("failed to fetch users from mysql")
		return
	}
	fmt.Println("users fetched from mysql", allUsers)

	p := postgress.NewConn()
	i = stores.NewStore(p)

	fmt.Printf("---------------- Printing type of interface %T ---------------------------------\n", i)
	u, ok = i.Create(u2)
	if !ok {
		fmt.Println("failed to create user postgres")
		return
	}
	fmt.Println("user created in postgres", u)

	allUsers, ok = i.FetchAll()
	if !ok {
		fmt.Println("failed to fetch users from postgres")
		return
	}
	fmt.Println("users fetched from postgres", allUsers)
}
