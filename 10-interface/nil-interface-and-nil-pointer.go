package main

import (
	"database/sql"
	"fmt"
)

// ***** nil pointers and nil interface are not same in Go*******

type DB interface {
	ReadAll() (string, bool)
}

type conn struct {
	db *sql.DB
}

func (c *conn) ReadAll() (string, bool) {
	if c.db == nil {
		return "", false
	}
	return "all values from the db", true
}

func main() {
	// nil pointer:- it means it is not pointing to a memory address or storing reference of datatype
	var c *conn = nil

	// d is nil // because it is not storing any type
	//var d DB
	//if d == nil {
	//	fmt.Println("interface value not set")
	//	return
	//}

	// if interface is holding any concrete type then it is not nil, even the type value itself could be nil
	// nil interface means no concrete type is present inside the interface
	var db DB = c /// i[conn]

	fmt.Printf("%#v", db)
	fmt.Printf("%#v", c)

	// db is not nil // storing a conn type in it
	if db == nil { // if interface is nil or not
		fmt.Println("conf is not set")
		return
	}

	//type assertion to get conn struct from interface
	c, ok := db.(*conn) // use ok variant to avoid panics
	if !ok {
		fmt.Println("interface is not holding conn struct")
		return
	}
	if c == nil {
		fmt.Println("no connection is set in conn")
		return
	}

	fmt.Println(db.ReadAll())

}
