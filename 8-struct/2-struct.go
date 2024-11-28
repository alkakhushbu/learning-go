package main

import "fmt"

func main() {

	type person struct {
		name string
		age  int
	}

	var p person
	p.name = "ajay"
	p.age = 18

	p1 := person{"raj", 18}

	fmt.Printf("%T\n", p.name) // string
	fmt.Printf("%T\n", p1)     // person struct

}
