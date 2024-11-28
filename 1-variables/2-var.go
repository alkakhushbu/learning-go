package main

import (
	"fmt"
)

func main() {
	// if we don't provide any values to the variable, they get initialized with there default values
	// var block
	var (
		name  = "ajay"
		age   = 30
		marks = 80.5
	)
	const key = "some key"
	fmt.Println(name, age, marks, key)

	// peek into it for design pattern
	//os.OpenFile()
	//os.O_RDONLY
	//time.Second
	//http.StatusNotFound

	//log.New(os.Stdout, log.LstdFlags)
	//l.Println("hello")
	//
	//math.MaxUint8
}
