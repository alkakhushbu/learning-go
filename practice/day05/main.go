package main

import "fmt"

/**
In main function create a slice of names.
   Add two elements in it.

   Create a function AddNames which appends new names to the slice
   Use double pointer concept to make AddNames function work
*/
func main() {
	names := []string{}
	names = append(names, "John")
	names = append(names, "Ronald")
	names = AddName(names, "Peter")
	fmt.Println(names)
}

func AddName(names []string, name string) []string {
	names = append(names, name)
	return names
}
