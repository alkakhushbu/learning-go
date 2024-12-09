package main

import (
	"fmt"
	"sync"
	"time"
)

//******

// Don't use this solution in real life, this is wrong one

// ********
func main() {

	wg := new(sync.WaitGroup)
	wgWorker := new(sync.WaitGroup)
	ch := make(chan int, 0)
	//m := &sync.Mutex{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 1; i <= 5; i++ {
			wgWorker.Add(1)

			// fan out pattern, spinning up n number of goroutines, for n number of task
			go func(i int) {
				defer wgWorker.Done()
				ch <- i

			}(i)
			time.Sleep(1 * time.Second)
		}

	}()

	// separate goroutine in main should not be created to tack workers
	// possibility of wgworker value is 0
	// maybe workers never got chace to run
	wg.Add(1)
	go func() {
		defer wg.Done()
		wgWorker.Wait()
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
