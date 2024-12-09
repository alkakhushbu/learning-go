package main

import (
	"fmt"
	"sync"
)

//******* Note:- This program has a bug*********///

func main() {

	wgWorker := new(sync.WaitGroup)
	wg := new(sync.WaitGroup)

	// select is used when we want to listen or send values to over a multiple channel
	// select should not be used with select if done channel pattern is being used,
	// a lot of times it is possible that workers finished first and done channel is closed
	// which would lead to quit the select early without receiving al the values

	// Solution, spin individual goroutines to range over the different channel
	// don't forget to close the channel where you are sending the values
	get := make(chan string, 1)
	post := make(chan string, 1)
	put := make(chan string, 2)

	// empty struct doesn't take any memory
	done := make(chan struct{}) // datatype doesn't matter

	wgWorker.Add(1)
	go func() {
		defer wgWorker.Done()
		get <- "get"
		close(get)
	}()

	wgWorker.Add(1)
	go func() {
		defer wgWorker.Done()
		//time.Sleep(50 * time.Millisecond)
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

	wg.Add(1)
	go func() {
		defer wg.Done()
		//time.Sleep(5 * time.Second)
		for {
			select {

			// whichever case is not blocking exec that first
			//whichever case is ready first, exec that.
			// possible cases are chan recv , send , default

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

	go func() {
		defer wg.Done()
		wgWorker.Wait()
		close(done) // close is a send signal, and select can recv it
	}()

	wg.Wait()

}
