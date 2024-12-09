package main

import (
	"log"
	"os"
)

func main() {
	f, err := os.Open("file.txt")
	// after calling the function the next step must always be error handling
	// no business logic until error is handled

	if err != nil {
		log.Println(err)
		return
	}

	// next step is to cleanup resources in defer
	defer func() {
		// do this when you want to make sure your resource cleaning is not failing
		err := f.Close() // when the main function ends file would be closed
		if err != nil {
			log.Println(err)
			return
		}
	}()

	// continue writing the business logic

	//db, err := sql.Open("sqlite3", "file.db")
	//defer db.Close()
}
