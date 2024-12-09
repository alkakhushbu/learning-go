package main

import (
	"fmt"
)

func main() {

	// defer statements execute when surrounding function returns or when your function ends
	defer fmt.Println(1) // defer maintains stack // first in last out
	defer fmt.Println(2)
	fmt.Println(3)
	panic("some panic message")
	//return
	// os.Exit // don't use until you want to quit the program immediately
	//log.Fatal("some fatal message") // don't use until you want to quit the program immediately
	fmt.Println(4)
}
