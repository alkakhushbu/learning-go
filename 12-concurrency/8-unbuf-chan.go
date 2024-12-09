/* q1. Create 4 functions
   Add(int,int),Sub(int,int),Divide(int,int), CollectResults()
   Add,Sub,Divide do their operations and push value to an unbuffered channel

   CollectResult() -> It would receive the values from the channel and print it
*/

package main

import (
	"fmt"
	"sync"
)

func Add(a, b int, ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	result := a + b
	ch <- result
}

func Sub(a, b int, ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	result := a - b
	ch <- result
}

func Divide(a, b int, ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	if b == 0 {
		ch <- 0
		return
	}
	result := a / b
	ch <- result
}

func CollectResults(ch chan int, numResults int, wg *sync.WaitGroup) {
	// defer wg.Done()
	fmt.Println("Result:", <-ch) // channel recieve is blockimh call untill there is a value
	fmt.Println("Result:", <-ch)
	fmt.Println("Result:", <-ch)
}

func main() {
	var wg sync.WaitGroup
	ch := make(chan int)
	wg.Add(3)

	go Add(10, 5, ch, &wg)
	go Sub(10, 5, ch, &wg)
	go Divide(10, 5, ch, &wg)

	CollectResults(ch, 3, &wg)
	wg.Wait()
}
