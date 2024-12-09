package main

import "fmt"

func main() {
	// var sInterface SquareInterface
	// sInterface = &SquareImpl{}
	//if we defer recoverPanic() here, the control does not reaches to end of function main when panic happens
	var sImpl *SquareImpl = nil
	defer recoverPanic()
	callSquare(sImpl)
	fmt.Println("end of main")
}

type SquareInterface interface {
	square(int) int
}

type SquareImpl struct {
	squareNum int
}

func (s *SquareImpl) square(num int) int {
	s.squareNum = num * num
	return num
}

// a pointer of interface type points to concrete type and not the object
// so when we pass a pointer of the concrete type pointing to nil,
// it works fine until it tries to access object fields
func callSquare(s SquareInterface) {
	//if we defer recoverPanic() here, the control does not reaches to end of function callSquare when panic happens
	// defer recoverPanic()
	val := s.square(4)
	fmt.Println(val)
	fmt.Println("end of callSquare")
}

func recoverPanic() {
	msg := recover()
	if msg != nil {
		fmt.Println("Panic happened:", msg)
		return
	}
}
