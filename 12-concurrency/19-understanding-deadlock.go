package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)
	//wg := sync.WaitGroup{}
	//c := make(chan int)
	//wg.Add(1)
	go func() {
		//defer wg.Done()
		for {
			time.Sleep(100 * time.Millisecond)
			fmt.Println("I'm sleeping")
		}

	}()
	//wg.Wait()
	// this program would not deadlock because one goroutine is running forever
	// the program deadlocks when go is sure that no more values would be coming to the channel
	// and there is no more point of waiting
	<-ch
}
