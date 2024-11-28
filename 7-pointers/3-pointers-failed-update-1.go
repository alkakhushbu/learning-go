package main

import "fmt"

var x = 10

func main() {
	// address of p is x80 = and the value it is storing is nil
	var p *int     // nil // default value of pointer is nil
	updateValue(p) //
	// after calling the update value p is still nil, as we never updated the pointer
	fmt.Println(*p) // 10 // access the memory at nil is always going to be a panic state
}

// address of p1 is x90 = nil
func updateValue(p1 *int) {

	p1 = &x // p1 = x120 // x90 = x120
	fmt.Println(*p1)
}
