package postgres

import "fmt"

type Conn struct {
	conn string
}

func NewConn(conn string) *Conn {
	fmt.Println("Creating new connection for postgres")
	newConn := &Conn{conn: conn}
	return newConn
}

func (c *Conn) Create(user string) error {
	fmt.Println("Creating new user for postgres")
	return nil
}

func (c *Conn) Update(name string) error {
	fmt.Println("Updating user for postgres")
	return nil
}

func (c *Conn) Delete(id string) error {
	fmt.Println("Deleting user for postgres")
	return nil
}
