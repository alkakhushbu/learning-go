package main

import "fmt"

type operation func(int, int) int

func main() {
	//for operate func in this case
	operate(add(), 10, 20)
	//operate(sub, 10, 20)

}

// operate func can accept function in op parameter,
// the function signature we are passing should match to op parameter type
func operate(op operation, x, y int) {
	sum := op(x, y)
	fmt.Println(sum)

}

func add() func(int, int) int {
	fmt.Println("in add function at the top")
	return func(x, y int) int {
		fmt.Println("doing add inside")
		return x + y
	}

}
