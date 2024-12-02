package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

/*
q1. Create a function that converts string to an integer

	    if any alphabets are passed, wrap strconv error and ErrStringValue error (create ErrStringValue error)

	    ErrStringValue contains a message that 'value is of string type' and return the wrapped errors
	    otherwise return the original error

	    use the regex to check if value is of string type or not
	    Hint: regexp.MatchString(`^[a-zA-Z]`, s)
	    fmt.Errorf("%w %w") // to wrap error

	    In main function check if ErrStringValue error was wrapped in the chain or not
	    If yes, log a message 'value must be of int type not string' and log original error message alongside as well
		*
*/
var ErrStringValue = errors.New("value is of string type")

func main() {
	// val, err := convert("123")
	val, err := convert("A123")
	// val, err := convert("$12")
	if err != nil {
		if errors.Is(err, ErrStringValue) {
			fmt.Println("Custom Error found:", err)
			return
		}
		fmt.Println("Error found: ", err)
		return
	}
	fmt.Println("Converted string to integer: ", val)

}

func convert(s string) (int, error) {
	val, err := strconv.Atoi(s)
	if err != nil {
		alphabet, regexError := regexp.MatchString(`^[a-zA-Z]`, s)
		if regexError != nil {
			return 0, fmt.Errorf("%w %w", regexError, err)
		}
		if alphabet {
			return 0, fmt.Errorf("%w %w", ErrStringValue, err)
		}
		return 0, fmt.Errorf("%w", err)
	}
	return val, nil
}
