package main

import "fmt"

func main() {
	var i []int // default value of slice is nil
	i[0] = 100  // this is only for update
	fmt.Println(i)
}
