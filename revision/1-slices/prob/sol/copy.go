package main

import "fmt"

func main() {
	x := []int{10, 20, 30, 40, 50, 60, 70}

	// creating a seperate backing array for b, so we can copy elems from x
	b := make([]int, len(x[1:4]), 100) // type , len , cap

	// dest,src
	copy(b, x[1:4])

	b[0] = 100
	fmt.Println(x)
	fmt.Println(b)
}
