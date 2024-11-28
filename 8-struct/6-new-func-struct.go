package main

import (
	"fmt"
	"log"
)

type Conf struct {
	db string
}

// New functions are used to intialize struct with some config values,
// See proj-struct folder for proper implementation

func NewConf(conn string) Conf {
	if conn == "" {
		// avoid in production until and unless you want your app to stop working
		// this will crash the program
		log.Fatal("empty connection string")
	}
	// try to open the connection, and if it is successful, return the connection
	return Conf{db: conn}
}

func (c Conf) AddToDb() {
	fmt.Println("adding to db", c.db)

}

func main() {
	c := NewConf("mysql://root:123456@localhost/test")
	c.AddToDb()
	log.Println(c.db)
}
