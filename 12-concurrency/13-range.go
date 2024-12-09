package main

import (
	"fmt"
	"sync"
)

func main() {

	wg := new(sync.WaitGroup)
	wgWorker := new(sync.WaitGroup)
	ch := make(chan int)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 1; i <= 5; i++ {
			//
			wgWorker.Add(1)
			// fan out pattern, spinning up n number of goroutines, for n number of task
			go func(i int) {
				// pass a local variable if you are using go version anything before 1.23
				// in older version go would create reference of i variable
				//and reuses the same memory of i for all the goroutine
				defer wgWorker.Done()
				ch <- i
			}(i)

		}

		// the wait value would always be correct here
		// because the loop would run first, and it would add the correct value in the wait group
		wgWorker.Wait()
		// close the channel when all workers are finished working
		close(ch)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// range gives a guarantee everything would be received
		for v := range ch {
			fmt.Println(v)
		}
	}()

	wg.Wait()

}
