package main

import "fmt"

func main() {
	// var sInterface SquareInterface
	// sInterface = &SquareImpl{}
	//if we defer recoverPanic() here, the control does not reaches to end of function main when panic happens
	var s *SquareImpl
	defer recoverPanic()
	callSquare(s)
	fmt.Println("end of main")
}

type SquareInterface interface {
	square(int)
}

type SquareImpl struct {
	squareVal int
}

func (s *SquareImpl) square(num int) {
	s.squareVal = num * num
}

func callSquare(s SquareInterface) {

	//if we defer recoverPanic() here, the control does not reaches to end of function callSquare when panic happens
	// defer recoverPanic()
	s.square(4)
	//fmt.Println(val)
	fmt.Println("end of callSquare")
}

func recoverPanic() {
	msg := recover()
	if msg != nil {
		fmt.Println("Panic happened:", msg)
		return
	}
}
