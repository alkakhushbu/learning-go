package main

import "fmt"

type person struct {
	name string
}

func main() {
	degree := float32(32.0)
	fmt.Printf("%T\n", degree)
	farenheit := convertTemp(degree)
	fmt.Println(farenheit)
}

func convertTemp(degree float32) float32 {
	farenheit := degree*(9/5) + 32
	return farenheit
}
