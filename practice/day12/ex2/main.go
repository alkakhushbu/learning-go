package main

import (
	"fmt"
	"sync"
)

func main() {
	for i := 0; i < 10; i++ {
		goroutine()
		fmt.Println(i)
	}
	// goroutine()
	// fmt.Println()
}

func goroutine() {
	wg := new(sync.WaitGroup)
	wgWorker := new(sync.WaitGroup)
	ch := make(chan int)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 5; i++ {
			// fan out pattern, spinning up n number of goroutines, for n number of task
			wgWorker.Add(1)
			go func(j int) {
				defer wgWorker.Done()
				ch <- i
			}(i)
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			wgWorker.Wait()
			close(ch)
		}()
	}()

	// wg.Add(1)
	// go func() {
	// 		defer wg.Done()
	// 		wgWorker.Wait() // until workers are not finished, we would wait
	// 		// close the channel if workers are done sending
	// 		close(ch)
	// }()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// time.Sleep(time.Second * 1)
		// range gives a guarantee everything would be received
		// a, ok := <-ch
		// if ok {
		// 	fmt.Println(a)
		// }
		for v, ok := <-ch; ok; {
			fmt.Print(v)
			v, ok = <-ch
		}
	}()

	wg.Wait()
	fmt.Print("Done")

}
