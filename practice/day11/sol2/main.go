package main

import (
	"fmt"
	"sync"
	"time"
)

/*
q2. Follow up to the previous question

	Spin up one anonymous goroutine in the work function
	This goroutine prints some stuff on the screen and sleeps for 100ms
	Make sure you wait for every goroutine to finish and end everything gracefully
*/
func main() {
	wg := new(sync.WaitGroup)
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go work(i, wg)
	}
	wg.Wait()
}

func work(id int, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println(id)
		time.Sleep(time.Millisecond * 100)

	}()
	defer wg.Done()
}
