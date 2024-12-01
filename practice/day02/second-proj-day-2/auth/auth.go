package auth

import (
	"fmt"
	"second-proj-day-2/user"
)

func Authenticate() {
	fmt.Println("authenticating user")
}

func Name() {
	currentUser := user.DbUser
	fmt.Println(currentUser)
}
