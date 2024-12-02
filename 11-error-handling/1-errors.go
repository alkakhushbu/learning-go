package main

import (
	"errors"
	"fmt"
	"log"
)

var user = make(map[int]string)

// variable which would be used to store errors should start with word Err

var ErrNotFound = errors.New("not found")

func main() {
	name, err := FetchRecord(1)
	if err != nil {
		// log.Println + os.Exit() // -> log.Fatal
		//log.Fatal(err) // would quit the app // should only be used when critical parts are failing usually
		// at startup and there is no point of continuing
		log.Println("FetchRecord:", err)
		return // return most of the times
	}
	fmt.Println("user name:", name)

}

// error must be the last value to be returned from function

func FetchRecord(id int) (string, error) {
	name, ok := user[id]
	if !ok {
		//return "", ErrNotFound // whenever error happens, set other values to their defaults
		return "", fmt.Errorf("user %d not found", id)
	}
	return name, nil
}
