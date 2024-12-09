package store

import (
	"errors"
	"strconv"
)

var books = make(map[string]int)

var ErrInvalidCount = errors.New("invalid book count")
var ErrNotFound = errors.New("book does not exist")

type BookStoreError struct {
	Func  string
	Input string
	Err   error
}

func (e *BookStoreError) Error() string {
	return "Invalid input \"" + e.Input + "\" passed to function : " + e.Func + ". Error message : " + e.Err.Error()
}

func AddBook(title string, counter int) error {
	if counter <= 0 {
		return invalidCountErr(strconv.Itoa(counter))
	}
	count, ok := books[title]
	if ok {
		count += counter
	}
	books[title] = count
	return nil
}

func invalidCountErr(s string) error {
	err := &BookStoreError{Func: "invalidCountErr", Input: s, Err: ErrInvalidCount}
	return err
}

func FetchBookCounter(name string) (int, error) {
	count, ok := books[name]
	if !ok {
		return 0, invalidBookNameErr(name)
	}
	return count, nil
}

func invalidBookNameErr(name string) error {
	err := &ErrorBookStore{Func: "invalidBookName", Input: name, Err: ErrNotFound}
	return err
}
