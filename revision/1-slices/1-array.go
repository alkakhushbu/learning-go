package main

import "fmt"

func main() {
	//Arrays size is fixed, we can't grow them
	var a [5]int
	a[0] = 1
	//a[5] = 2 // at compiler level it would fail
	b := [5]int{10, 20}           // 10,20,0,0
	c := [...]int{10, 20, 30, 50} // ... would create the array size according to the number of values passed
	// if three values are passed, then size would three, but after creation we cant grow the array
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)

}
