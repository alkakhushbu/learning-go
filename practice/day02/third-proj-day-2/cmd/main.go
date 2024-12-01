package main

import (
	"fmt"
	"third-proj-day-2/auth"
)

/*
Use one of the two options to run the code:

	go run .
	go run . main.go
	go build => builds entire project
*/
func main() {
	Setup()
	auth.Authenticate()
	fmt.Println("End of func main")
}
