package main

import "fmt"

func main() {
	x := []int{10, 20, 30, 40}
	updateSlice(x)
	fmt.Println(x)
	x = AppendSlice(x, 50)
	fmt.Println(x)
}

// if you want to update your slice then pass the slice normally
func updateSlice(s []int) {
	s[0] = 100
	// but never ever append in the update function without returning the slice
	// it would not work as expected
}

func AppendSlice(s []int, n int) []int {
	// if you want to append then always return your slice
	s = append(s, n)
	// without returning, you would not update the slice in the main function
	return s
}
