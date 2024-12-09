package main

import (
	"fmt"
	"strings"
)

/*
q4. Modify the above function to perform the following action
    stringManipulation(trimSpace(), "\ngfdngbk \n"))
    Instead of passing trimSpace you need to call the trimSpace function and make the program work
    Hint: you need to return a function with the signature of what stringManipulation accepts as first parameter
*/

type operation func(string) string

func stringManipulation(fn operation, s string) string {
	return fn(s)
}

func trimSpace() operation {
	return func(s string) string {
		return strings.TrimSpace(s)
	}
}

func toUpper() operation {
	return func(s string) string {
		return strings.ToUpper(s)
	}
}

func main() {
	a := stringManipulation(trimSpace(), "\ngfdngbk \n")
	fmt.Println(a)

	b := stringManipulation(toUpper(), "Hello World! I will try to make you a better place!")
	fmt.Println(b)
}
