package main

import "fmt"

func main() {
	users := []string{"Darcy", "Bingley", "Jane", "Lizzie"}
	emp := users[1:3]
	inspect(users)
	inspect(emp)

	users = append(users, "Bennet", "Lydia", "Lucas")
	inspect(users)
	inspect(emp)

	emp = append(emp, "Joe", "Wright")
	inspect(users)
	inspect(emp)
}

func inspect(users []string) {
	fmt.Println("len: ", len(users), "capacity: ", cap(users))
	fmt.Println(users, "\n")
}
