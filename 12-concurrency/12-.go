package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := new(sync.WaitGroup)
	// this waitgroup would track the number of worker goroutines spawned
	wgWorker := new(sync.WaitGroup)
	ch := make(chan int, 5)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 5; i++ {
			// running n number of task , all goroutines are pushing to the same channel

			wgWorker.Add(1)
			go func(i int) {
				defer wgWorker.Done()

				fmt.Println("Worker", i)
				ch <- i

			}(i)
		}

		//we need to block our goroutine before closing the channel
		//because we want to make sure all the work
		// is done and finished // closing a channel will stop the for range loop
		//sending is finished over the channel ch

		// this goroutine would only run, when correct counter is added to the wgworker waitgroup
		// so this version works
		go func() {
			defer wg.Done()
			// it would wait until counter not zero
			wgWorker.Wait() // waiting until the worker goroutines are not finished
			close(ch)
		}()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range ch {

			fmt.Println("Received", v)
			//ranging until the channel is not closed
			//range would receive all the remaining values even after the channel is closed
			fmt.Println(v)
		}
	}()

	wg.Wait()

}
