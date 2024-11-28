package main

import "fmt"

type user struct {
	name string
	age  int
}

func main() {
	u := user{"jack", 18}
	updateUser(&u) // passing the reference of the u object
	fmt.Println(u)
}

// passing a struct to a function which accepts struct value as a pointer
func updateUser(u *user) {
	u.age = 20 // we don not have to dereference the struct object to access the field
}
