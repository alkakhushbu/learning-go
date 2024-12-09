package main

import (
	"fmt"
	"time"
)

// concurrency is dealing with a lot of things at once
// parallelism is doing multiple things at once

// Concurrency is not Parallelism by Rob Pike
// https://www.youtube.com/watch?v=oV9rvDllKEg&ab_channel=gnbitcom
func main() {
	//panic("panic") // reveals goroutine id and name
	go hello()

	fmt.Println("end of the main")

	// worst case to wait for goroutines
	time.Sleep(time.Second) // sleep would put the main in the blocking state
	// go scheduler would pick up another goroutine if available

}

func hello() {
	fmt.Println("Hello, world!")
}
