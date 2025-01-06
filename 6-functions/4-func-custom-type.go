package main

import "fmt"

type money int

// defining a custom type for the function,
// type typeName functionSignature
type operation func(int, int) int

func main() {
	// add would be passed to op parameter, 10 and 20 would be passed to x and y variables
	operate(add, 10, 20)
	operate(sub, 30, 10)

}

// operate func can accept function in op parameter,
// the function signature we are passing should match to op parameter type
func operate(op operation, x, y int) {
	sum := op(x, y)
	fmt.Println(sum)

}

// just like what we do in gin handlerFunc
func testOp(next operation, x, y, z int) int {
	return next(next(x, y), z)
}

func add(a, b int) int {
	return a + b
}

func sub(a, b int) int {
	return a - b
}
