package main

import "fmt"

/**
q2.     second-proj-day-2/
        ├── main.go
        ├── go.mod
        ├── auth/
        │   └── auth.go
        └── user/
            └── user.go
    In auth package create two functions
    1. Authenticate
    Authenticate function simply prints a message, authenticating user
    2. Name
    This function prints the Name of the user.
    Note:- to print the name of the user,
    use the user package to know who is the current user

    In user package create one global variable, and one func named as AddToDb
    1. AddToDb
    This function accepts database name as string
    It calls the Authenticate function from auth package
    At last it prints a msg, Adding to db DatabaseName [var which was accepted in the parameter]

    Global Variable
	**/

func main() {
	var (
		a = 100
		b = true
		c = "World"
		// d string
	)
	fmt.Printf("%d %v %s \n", a, b, c)
	fmt.Println("Hello World!")
}
