package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
)

// Error type should end with the word Error
// Error types should not be used for any domain data

type QueryError struct {
	Func  string
	Input string
	Err   error
}

// Error() method is implemented to implement error interface
func (q *QueryError) Error() string {
	// formatting the output string
	return "main." + q.Func + ": " + "input " + q.Input + " " + q.Err.Error()
}

var ErrNotFound = errors.New("not found")

func SearchSomething(id int) (string, error) {
	// assume that search code is written and we need to return an error
	return "", &QueryError{
		Func:  "SearchSomething",
		Input: strconv.Itoa(id),
		Err:   ErrNotFound,
	}
}

func main() {

	user, err := SearchSomething(11)
	if err != nil {
		// creating a nil pointer to QueryError struct
		var qe *QueryError       // nil
		if errors.As(err, &qe) { // checking if struct was present in the chain, pass the reference to &qe
			fmt.Println(qe.Func) // we can access individual fields if needed, or take some specific actions
			fmt.Println(qe.Err)
			fmt.Println(err)
			return
		}
		log.Println(err)
		return
	}

	log.Println(user)

}
