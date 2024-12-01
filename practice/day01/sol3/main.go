package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("Hello World!!")
	l := log.New(os.Stdout, "Prefix: ", log.LUTC)
	l.Println("hello this is first log message")
	log.Println("hello this is first log message")
}

// q3. Use log.New() and print a log message
//     use shortfile flag and stdFlags for flag argument (bit of Google required)
// // Hint:- first value to this function could be os.Stdout
// // Printing: l.Println("hello this is first log message")

// q4. Print default values and Type names of variables from question 2 using printf
// // Quick Tip, Use %v if not sure about what verb should be used,
// // but don't use it in this question :)
// // but generally using %v should be fine
