package main

import "genspark/slice"

// https://go.dev/ref/spec#Appending_and_copying_slices
func main() {
	// len: number of elements your slice is storing,
	// or number of elems slice is referring to in backing array

	// cap: number of elems your slice can accommodate
	i := []int{1, 2, 3, 4, 5}

	slice.Inspect("i", i)

	i = append(i, 6)
	slice.Inspect("i", i)

	i = append(i, 7)
	slice.Inspect("i", i)
}

// https://go.dev/ref/spec#Appending_and_copying_slices
/*
	append func working

	i := []int{10, 20, 30, 40, 50 } // len = 5 , cap =5
	append(i,60) // not enough cap so allocation is going to happen

//  sufficiently large underlying array.
	underlying array -> [10 20 30 40 50,60,{},{}] len =6 cap = 8

append(i,70,90,300) // allocation would happen as we don't have enough cap to fit three values
	underlying array -> [10 20 30 40 50,60,70,80,90, , , , ] len =9 cap = 13

	If the capacity of s is not large enough to fit the additional values, append allocates a new,
    sufficiently large underlying array that fits both the existing slice elements and the additional values.
    Otherwise, append re-uses the underlying array.
*/
