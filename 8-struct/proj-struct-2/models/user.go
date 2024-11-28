package models

import (
	"fmt"
	"proj-struct/database"
)

//create a struct for models, and make createUser as a method, and access the conf value using that

type Service struct {
	c database.Conf
}

func NewService(c database.Conf) *Service {
	return &Service{c: c}
}

type user struct {
	name string
}

var users []user

func (s Service) CreateUser(name string) {
	fmt.Println("adding to db", s.c)
	fmt.Println("creating the user", name)
	users = append(users, user{name: name})
}
func (s Service) FetchUsers() {
	fmt.Println("getting users")
	fmt.Println(users)
}
