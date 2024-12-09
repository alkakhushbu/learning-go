package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	wgWorker := new(sync.WaitGroup)
	wg := new(sync.WaitGroup)

	// select is used when we want to listen or send values to over a multiple channel
	get := make(chan string, 1)
	post := make(chan string, 1)
	put := make(chan string, 1)

	// empty struct doesn't take any memory
	done := make(chan struct{}) // datatype doesn't matter

	wgWorker.Add(1)
	go func() {
		defer wgWorker.Done()
		get <- "get"
	}()

	wgWorker.Add(1)
	go func() {
		defer wgWorker.Done()
		time.Sleep(50 * time.Millisecond)
		post <- "post"
	}()

	wgWorker.Add(1)
	go func() {
		defer wgWorker.Done()
		put <- "put"
		put <- "p1"

	}()

	// not efficient // because we have to wait for get even if it is taking long time execute
	//fmt.Println(<-get)
	//fmt.Println(<-post)
	//fmt.Println(<-put)

	// the problem with below loop is if less or more values are sent it would not work
	// and deadlock would happened
	//for i := 0; i < 3; i++ {
	//	// whichever case is not blocking exec that first
	//	//whichever case is ready first, exec that.
	//	// possible cases are chan recv , send , default
	//	select {
	//	case g := <-get:
	//		fmt.Println(g)
	//	case p := <-post:
	//		fmt.Println(p)
	//	case pu := <-put:
	//		fmt.Println(pu)
	//
	//	}
	//}

	wg.Add(1)
	go func() {
		defer wg.Done()
		// time.Sleep(2 * time.Second)
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

	wg.Add(1)
	go func() {
		defer wg.Done()
		wgWorker.Wait()
		close(done) // close is a send signal, and select can recv it
	}()

	wg.Wait()

}
