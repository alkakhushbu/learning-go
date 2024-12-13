package main

import (
	"revision/slice"
)

func main() {
	i := []int{10, 20, 30, 40, 50}
	slice.Inspect("i", i)
	// append is used to grow the slice
	i = append(i, 60, 70, 80)
	slice.Inspect("i", i)

	i = append(i, 90)
	slice.Inspect("i", i)

	//cap(i)

	//b := []int{100, 200, 300}
	//i = append(i, b...) // ... , here it means we are unpacking the slice
	//fmt.Println(i)
	//fmt.Println(len(i))
}
