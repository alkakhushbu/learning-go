package main

import (
	"app/stores"
	"app/stores/mysql"
	"app/stores/postgres"
)

func main() {
	var conn stores.Database
	conn = mysql.NewConn("mysql")
	conn.Create("alka")
	conn.Update("alka")
	conn.Delete("alka")

	conn = postgres.NewConn("postgres")
	conn.Create("alka")
	conn.Update("alka")
	conn.Delete("alka")

}
