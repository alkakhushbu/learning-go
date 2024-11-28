package main

import (
	"fmt"
	"genspark/slice"
)

func main() {
	x := make([]int, 0, 5)
	x = append(x, 1, 2, 3)
	update(x)

	fmt.Println(x)
}

func update(s []int) {
	// if you want to update the value , then passing the reference of existing backing array is fine
	s[0] = 100

	// if you want to append
	// if you are appending the values to the slice, always return the slice
	// returning the slice would ensure that caller has updated reference
	s = append(s, 111)

	slice.Inspect("after s", s)

	//and slice is not returned back, then code review should stop

}
