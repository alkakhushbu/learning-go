package main

import "fmt"

func main() {

	// any can store any type of value
	// it is an empty interface with 0 methods defined
	// any is an alias of interface{}
	var a any // var i interface{}
	a = 10
	a = struct{ a int }{}
	a = true
	var b bool
	a = "hello"
	_, _ = a, b
	//b, ok := a.(bool) // type assertion to fetch concrete value from any type
	//if !ok {
	//	fmt.Println("not of bool type")
	//	return
	//}
	//fmt.Println(b)
	display(10, 20, 30, 40, 50, 50)
}

// display is a variadic function, i is a variadic parameter, it means i can pass any number of values to the i parameter
func display(a int, i ...any) { // variadic parameter should be the last parameter in the function
	fmt.Printf("%T\n", i) // variadic parameter is a slice under the hood
	fmt.Println(i)
}
