package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	greet()
	fmt.Println("main function done")
}

func greet() {
	data := os.Args[1:]
	if len(data) != 3 {
		log.Println("please provide, name, age, marks")
		return
	}
	name := data[0]
	ageString := data[1]
	marksString := data[2]
	var err error // default value of error is nil // it indicates no error

	age, err := strconv.Atoi(ageString)
	if err != nil {
		log.Println("please provide valid age", err)
		return
	}
	marks, err := strconv.Atoi(marksString)
	if err != nil {
		log.Println("please provide valid marks", err)
		return
	}

	fmt.Println("hello", name, "you are", age, "and you have", marks, "marks")

}
