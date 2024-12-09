package main

import "fmt"

func main() {
	x := []int{0}
	fmt.Println(x, len(x), cap(x))
	for i := 1; i <= 1_000_000; i++ {
		beforeCap := cap(x)
		x = append(x, i)
		afterCap := cap(x)
		if beforeCap != afterCap {
			changeInPercent := float32(afterCap-beforeCap) / float32(beforeCap) * 100
			fmt.Println("beforeCap:", beforeCap, ", afterCap:", afterCap, ", changeInPercent:", changeInPercent)
		}
	}
	for i := 1_000_000 ; i >= 1 ; i-- {
		beforeCap := cap(x)
		x = append(x, i)
		afterCap := cap(x)
		if beforeCap != afterCap {
			changeInPercent := float32(afterCap-beforeCap) / float32(beforeCap) * 100
			fmt.Println("beforeCap:", beforeCap, ", afterCap:", afterCap, ", changeInPercent:", changeInPercent)
		}
	}
	// fmt.Println(x)
}
