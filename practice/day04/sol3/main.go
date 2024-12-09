package main

import (
	"fmt"
	"strings"
)

type operation func(string) string

func stringManipulation(fn operation, s string) string {
	return fn(s)
}

func trimSpace(s string) string {
	return strings.TrimSpace(s)
}

func toUpper(s string) string {
	return strings.ToUpper(s)
}

func greet(s string) string {
	return "Hello " + s + "!"
}

func main() {
	a := stringManipulation(trimSpace, "          Have a nice day!           ")
	fmt.Println(a)

	b := stringManipulation(toUpper, "Hello World")
	fmt.Println(b)

	c := stringManipulation(greet, "Alka")
	fmt.Println(c)
}
