package main

import (
	"fmt"
	"time"
)

func main() {
	// we can't store values returned from the goroutine
	//s := go greet() //spin / spawn

	// this is a hack, we created a goroutine, inside it we called the original func
	// and stored the returned value in a var
	go func() {
		s := greet()
		fmt.Println(s)
	}()

	fmt.Println("some work going on")
	fmt.Println("end of main")
	time.Sleep(time.Second) // worst practice
}

// returning values from the go routine have no use
// returned values can't be stored by the main function in this case
func greet() string {
	return "Hello, World!"
}
