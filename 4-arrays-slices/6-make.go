package main

import (
	"fmt"
	"genspark/slice"
)

func main() {
	// make preallocates a backing array
	// very helpful if we have idea about the number of request ,
	//we can create size of array according to that
	// in this case append doesn't have to allocate the memory each time
	i := make([]int, 0, 50) // type, len, cap
	slice.Inspect("i", i)
	i = append(i, 100)
	fmt.Println(i)
}
