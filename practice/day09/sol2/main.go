package main

import (
	"errors"
	"fmt"
	"runtime/debug"
)

/*
q2. Create 3 functions f1, f2, f3
    f1() call f2(), f2() call f3()
    each layer would return the error, wrap the error from each layer
    print stack trace using debug.Stack to get a complete stack trace
*/

func main() {
	_, err := f1()
	if err != nil {
		fmt.Println("Error found: ", err)
		return
	}
	fmt.Println("End of main")
}

func f1() (string, error) {
	// errors.Unwrap()
	val, err := f2()
	if err != nil {
		return "", fmt.Errorf("%w, %w", errors.New("Layer - 1"), err)
	}
	return val, nil

}

func f2() (string, error) {
	val, err := f3()
	if err != nil {
		return "", fmt.Errorf("%w, %w", errors.New("Layer - 2"), err)
	}
	return val, nil
}

func f3() (string, error) {
	fmt.Println("Error found: ", string(debug.Stack()))
	return "", errors.New("Layer - 3")
}
