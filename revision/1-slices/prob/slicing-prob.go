package main

import (
	"fmt"
	"revision/slice"
)

func main() {
	i := []int{10, 20, 30, 40, 50}
	b := i[1:3] // 20,30
	b[0] = 100
	fmt.Println(b)
	fmt.Println(i)

	slice.Inspect("a", i)
	slice.Inspect("b", b)

}

//func slice(b []int) {
//	b[0] = 100
//}
