package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

// DB as global var not recommended
var DB *sql.DB

// this is used to initialize the state for the current package
// not recommend to be used most of the times.
// hard to test, hard to know when it runs
func init() {
	const (
		host     = "localhost"
		port     = "5433"
		user     = "postgres"
		password = "postgres"
		dbname   = "postgres"
	)
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("pgx", psqlInfo)
	if err != nil {
		panic(err)
	}
}

func main() {
	err := DB.Ping()
	if err != nil {
		// if db is not connected, no point to continue
		panic(err)
	}

}

func UpdateAuthor() {
	//BeginTx would start the transaction
	tx, err := DB.BeginTx(context.Background(), nil)
	if err != nil {
		log.Println(err)
		return
	}

	// calling rollback multiple times have no effect after commit
	// rollback would roll back any changes if function return early without commit
	defer func() {
		err := tx.Rollback()
		if err != nil {
			log.Println(err)
			return
		}
	}()
	updateQuery := `UPDATE author
					SET name = $1
					WHERE email = $2;`

	_, err = tx.Exec(updateQuery, "ABC", "john.doe@example.com")
	if err != nil {
		log.Println(err)
		return
	}

	_, err = tx.Exec(updateQuery, "John1", "john.doe@example.com")
	if err != nil {
		log.Println(err)
		return
	}

	// only if both transaction finishes then only we would commit
	// All or None concept
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return
	}
}
