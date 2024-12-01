package user

import (
	"fmt"
	// "second-proj-day-2/auth"
)

var DbUser string = "alka"

func AddToDb(name string) {
	// auth.Authenticate() // this won't work, solution will be in the next week
	fmt.Printf("Adding to db DatabaseName, %s\n", name)
}

func currentUser() {
	fmt.Println(DbUser)
}
