package main

import (
	"fmt"
	"revision/slice"
)

func main() {
	x := []int{10, 20, 30, 40, 50, 60}

	b := x[1:4] //  20, 30, 40, 50,
	slice.Inspect("b", b)
	slice.Inspect("x", x)
	fmt.Println("after append")
	//b = append(b, 999)
	b = append(b, 999, 888, 777)

	slice.Inspect("b", b)
	slice.Inspect("x", x)

}
