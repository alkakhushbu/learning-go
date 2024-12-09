package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	wg := new(sync.WaitGroup)
	ch := make(chan int) // unbuf chan
	wg.Add(2)
	go func() {
		defer wg.Done()

		fmt.Println("sender picked up")
		ch <- 1
		// in our case sender would be picked up first , because recv is blocked by a sleep
		// but when go scheduler try to run this routine, this would also block,
		// because there is no recv to send the values
		// so this goroutine goes in the blocked state until recv is up
		// once recv routine runs this routine would be unblocked, and it will send the value and move on
	}()
	go func() {
		defer wg.Done()
		time.Sleep(time.Second)
		fmt.Println("recv picked up")
		fmt.Println(<-ch)
		// recv would be read after a second
		//and it would wait for sender to send the values
	}()
	wg.Wait()
}
