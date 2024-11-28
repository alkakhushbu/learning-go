package main

import "fmt"

type user struct {
	name string
	age  int
}

type users map[int]*user

func main() {
	//u := make(users)
	//u["john"] = user{"john", 20}
	u := users{1: &user{"john", 20}, 2: &user{"jane", 21}}

	usr, ok := u[2] // user, ok (true,false) if that value exist or not

	if !ok {
		fmt.Println("user not found")
		return
	}
	fmt.Println(usr)
}
