package main

import (
	"fmt"
	"sync"
)

/*
q1. create a function work that takes a work id and print work {id} is going on
    In the main function run a loop to run work function 10 times
    make the work function call concurrent
    Make sure your program waits for work function to finish gracefully
*/

func main() {
	wg := new(sync.WaitGroup)
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go work(i, wg)
	}
	wg.Wait()
}

func work(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(id)
}
