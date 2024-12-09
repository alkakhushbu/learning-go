package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	// select is used when we want to listen or send values to over a multiple channel
	get := make(chan string)
	post := make(chan string)
	put := make(chan string)

	done := make(chan struct{}) // datatype doesn't matter

	wg.Add(1)
	go func() {
		defer wg.Done()
		// time.Sleep(1 * time.Second)
		get <- "get"
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// time.Sleep(50 * time.Millisecond)
		post <- "post"
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		put <- "put"
		put <- "p1"
	}() // close is a send signal, and select can recv it

	// didn't work
	// wg.Wait()
	// close(done)

	var wg2 sync.WaitGroup
	wg2.Add(1)
	go func() {
		defer wg2.Done()
		time.Sleep(time.Second * 5)
		for {
			select {
			case g := <-get:
				fmt.Println(g)
			case p := <-post:
				fmt.Println(p)
			case pu := <-put:
				fmt.Println(pu)
			case <-done:
				fmt.Println("all values are received")
				return
			}
		}
	}()

	wg.Wait()
	close(done)
	wg2.Wait()
}
