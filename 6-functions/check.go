package main

import "fmt"

func main() {
	func(x, y int) int {
		return x + y
	}(10, 20) // () this is to invoke the anonymous function

	//fmt.Println(greet)
	f := greet()

	fmt.Println(f())

}

func greet() func() string {
	return func() string {
		return "hello"
	}
}
