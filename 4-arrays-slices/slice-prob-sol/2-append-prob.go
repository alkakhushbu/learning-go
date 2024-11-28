package main

import "genspark/slice"

func main() {
	x := []int{10, 20, 30, 40, 50, 60}

	// cap of b calculation: - it starts from the 0 index of b till the end of the backing array, in this case 5 is cap
	b := x[1:5]
	slice.Inspect("b", b)

	//below line would change x, b refers to the same backing array,
	//there is room to add one more element to the existing backing array
	//so it would add the value to the backing array refer  by x
	b = append(b, 888)
	// but if two values would be added to the backing array, b= append(b,99,88)
	// no additional cap is left so it would create a new backing array, so it would not effect the x backing array
	b[0] = 1000
	slice.Inspect("b", b)
	slice.Inspect("x", x)

}
