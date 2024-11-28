package database

// avoid exported global variable most of the times
// anyone can change the state of it

var DB string

func InitDB(conn string) {
	DB = conn
}
