package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]
	if len(args) != 3 {
		fmt.Println("3 parameters should be passed in command line arguments")
	}
	operator, ok := args[0]
	first, ok := strconv.Atoi(args[1])
	if !ok {
		fmt.Println("Invalid first integer. Please enter a valid integer.")
	}
	second, ok := strconv.Atoi(args[2])
	if !ok {
		fmt.Println("Invalid second integer. Please enter a valid integer.")
	}
	

	fmt.Println(os.Args)
}
