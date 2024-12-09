package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := new(sync.WaitGroup)
	ch := make(chan int)
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 1; i <= 5; i++ {
			ch <- i
		}
		// close signal range that no more values be sent and it can stop after receiving remaining values
		// close the channel when sending is finished

		// we can't send more values after a channel is closed
		//ch <- 6 // panic: send on closed channel // channel is closed
		close(ch)
	}()

	go func() {
		defer wg.Done()
		for v := range ch {
			fmt.Println(v)
		}
	}()

	wg.Wait()
}
