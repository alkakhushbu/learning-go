package main

import "fmt"

func main() {
	x := []int{10, 20, 30, 40, 50, 60, 70}

	// creating a seperate backing array for b, so we can copy elems from x
	b := make([]int, len(x[1:4]), 100)

	// make sure enough length is present in the backing array,
	// copy only cares about length
	// dest,src
	copy(b, x[1:4])
	b[0] = 1000 // different copy of , so no changes in x
	b = append(b, 99)
	fmt.Println(b)
	fmt.Println(x)
}
