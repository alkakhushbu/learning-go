package main

import (
	"fmt"
	"strconv"
)

func main() {
	msg, ok := hello("Raj", 20, 100)
	if !ok { // !ok is like ok==false //ok == true
		fmt.Println("Invalid input")
		return
	}
	fmt.Println(msg)
}

func hello(name string, age, marks int) (string, bool) {
	if name == "" {
		return "", false
	}
	if age == 0 {
		return "", false
	}
	if marks == 0 {
		return "", false
	}
	return "Hello " + name + " marks: " + strconv.Itoa(marks) + " age: " + strconv.Itoa(age), true

}
