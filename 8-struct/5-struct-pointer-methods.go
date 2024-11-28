package main

import "fmt"

// https://go.dev/doc/faq#methods_on_values_or_pointers:~:text=Should%20I%20define%20methods%20on%20values%20or%20pointers%3F%C2%B6
type author struct {
	name string
	age  int
}

func main() {
	a := author{name: "zhangsan", age: 18}
	a.UpdateName("Jack")
	a.PrintName()
}

func (a *author) UpdateName(name string) {
	a.name = name

}

func (a *author) PrintName() {
	fmt.Println(a.name)
}
