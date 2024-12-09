package main

import (
	"fmt"
	"sync"
	"time"
)

//var wg = new(sync.WaitGroup)

func work(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	//time.Sleep(3 * time.Second)

	var wgTwo = new(sync.WaitGroup)

	wgTwo.Add(1)
	go func() {
		defer wgTwo.Done() // decrement the counter
		fmt.Println("this is anonymous function, sleeping for 100ms ", id)
		time.Sleep(1 * time.Second)
	}()

	fmt.Printf("Work %d is going on\n", id)
	wgTwo.Wait()
}

func main() {
	//var wg sync.WaitGroup

	var wg = new(sync.WaitGroup)

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go work(i, wg) // spinning 3 goroutines // each function call spins up a goroutine
	}

	wg.Wait()
	fmt.Println("All work completed.")
}
