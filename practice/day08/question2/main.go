package main

import (
	"fmt"
	"log"
)

type user struct {
	name  string
	email string
}

// Write implements io.Writer.
func (u user) Write(p []byte) (n int, err error) {
	fmt.Println(u)
	return len(u.name), nil
}

// func (u *user) Write(p []byte) (n int, err error) {
// 	fmt.Println(*u)
// 	return len(u.name), nil
// }

func main() {
	u := user{"raj", "raj@email.com"}
	// l := log.New(&u, "log: ", log.LstdFlags)
	l := log.New(u, "log: ", log.LstdFlags)

	l.Println("Hello, log")
}
