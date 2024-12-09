package main

import "fmt"

func main() {
	c := make(chan int)
	go sendData(c)
	receiveData(c)
}

func sendData(c chan<- int) {
	c <- 1

	// <-c // sendData accepts a channel with sending permission only
	close(c) // closing is a send signal,
}

func receiveData(c <-chan int) {
	// we can only receive in this function
	for v := range c {
		fmt.Println(v)
	}
}
