package main

import "fmt"

func main() {
	x := 10

	updateVal(&x)
	fmt.Println(x)

}
func updateVal(px *int) {
	var abc = 20
	// we have changed the reference of px to store a new variable abc
	px = &abc
	// this would not update the main function x variable, it would update abc
	*px = 100
}
