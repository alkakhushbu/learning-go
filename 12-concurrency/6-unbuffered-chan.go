package main

import "fmt"

// https://go.dev/ref/spec#Send_statements
// A send on an unbuffered channel can proceed if a receiver is ready.
// send will block until there is no recv
func main() {

	// no size specified, which means it is an unbuffered chan
	ch := make(chan int, 0) // you need to make sure to allocate memory to the channel using make

	// we never send and recv in the same goroutine
	// channels are only meant to be used in concurrent programming
	ch <- 10  // send signal // would block forever // receiver gets read in next line which would never happen
	a := <-ch // blocking forever no sender    // recv

	fmt.Println(a)
}
