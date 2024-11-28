package main

import "fmt"

func main() {
	//Arrays size is fixed, we can't grow them
	var i [10]int
	i[0] = 100
	//i[11] = 1100 // at compiler level it would fail

	b := [5]int{10, 20}
	c := [...]int{10, 20, 30, 50} // ... would create the array size according to the number of values passed
	// if three values are passed, then size would three, but after creation we cant grow the array
	fmt.Println(c)

	fmt.Println(i, "len", len(i))
	fmt.Println(b, "len", len(b))

	//range would traverse the data structures like map, arrays, slice
	for index, value := range b {
		fmt.Println(index, "value", value)
	}
	CompareArrays([5]int{10, 20}, b)
}

func CompareArrays(a [5]int, b [5]int) {
	// if you want to compare an array, the size of both arrays must be same
	if a == b {
		fmt.Println("array is same")
		return
	}
	fmt.Println("array is not same")
}
