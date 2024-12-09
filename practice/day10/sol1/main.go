package main

import (
	"errors"
	"fmt"
	"sol1/store"
)

func main() {
	err := store.AddBook("Learning go", -1)
	if err != nil {
		var errBookStore *store.BookStoreError
		if errors.As(err, &errBookStore) {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Book added")
	}

	
	count, err := store.FetchBookCounter("Learning Go")
	if err != nil {
		var errBookStore *store.BookStoreError
		if errors.As(err, &errBookStore) {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Book count: ", count)
	}

}
