package main

import "fmt"

func main() {
	i := []int{10, 20, 30, 40, 50}
	// append is used to grow the slice
	i = append(i, 60)
	fmt.Println(i)
	fmt.Println(len(i))
}
