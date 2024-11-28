package main

import "fmt"

func main() {
	a := 10

	// copy the address of a to p
	p := &a // p is a pointer, it stores the memory address of a
	fmt.Println(&a)
	fmt.Println(*p) // derefrecning the memory // access value at the address
	fmt.Println(p)  // what address p holds
	fmt.Println(&p)
}
