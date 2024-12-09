package main

// example1
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
		worker := new(sync.WaitGroup)
		// defer close(ch)
		// defer wg2.Wait()
		for i := 1; i <= 5; i++ {
			worker.Add(1)
			// fan out pattern, spinning up n number of goroutines, for n number of task
			go func() {
				defer worker.Done()
				ch <- i
			}()
		}
		worker.Wait()
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
