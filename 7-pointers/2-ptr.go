package main

import "fmt"

func main() {
	x := 10
	update(&x)
	fmt.Println(x)
}

func update(p *int) {
	*p = 20 // changing the memory address directly of x variable from the main function
}
