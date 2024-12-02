package main

import (
	"fmt"
	"log"
)

type user struct {
	name  string
	email string
}

func (u user) Write(p []byte) (n int, err error) {
	fmt.Printf("sending notification to %s %s\n", u.name, u.email)
	return len(p), nil
}

func main() {
	u := user{"raj", "raj@email.com"}
	// log.New accepts a writer interface
	// user type implements writer interface so we can pass its value to log.New
	l := log.New(u, "log: ", log.LstdFlags)
	l.Println("Hello, log")
}
