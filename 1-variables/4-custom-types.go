package main

import (
	"fmt"
)

// this should be created when you need to add two things
// 1. Readability
// 2. Behaviour using methods
type money int // this is a new type

// this should be created when you need to add one thing
// 1. Readability
type otp = int // = creates an alias

func main() {
	var r money = 10
	var d money = 20
	var i int = int(r) // r is of type money, we can't assign it to int directly
	fmt.Println(r, d, i)

	var x otp = 100
	var y int = x // x is of type otp, which is an alias, could be assigned anywhere where the underlying type is same

	//var b []byte // byte is an alias
	fmt.Println(y)

	//time.Duration
}

func AddDatatoFile(u uint8) {

}

func AddToFile(b []byte) {

}
