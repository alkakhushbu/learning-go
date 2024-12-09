package main

import (
	"fmt"
	"sync"
	"time"
)

// var // map
var cab int = 1

func main() {
	wg := new(sync.WaitGroup)
	m := new(sync.Mutex)
	names := []string{"a", "b", "c", "d"}
	for _, name := range names {
		wg.Add(1)
		go bookCab(name, wg, m)

	}

	wg.Wait()
}
func bookCab(name string, wg *sync.WaitGroup, m *sync.Mutex) {
	defer wg.Done()
	fmt.Println("welcome to the website", name)
	fmt.Println("some offers for you", name)

	//until the lock is not released
	//any read , write from other goroutines would not be allowed after lock is acquired
	func() {
		m.Lock()
		defer m.Unlock()
		if cab >= 1 {
			//critical section // area where we access shared resource
			// we are protecting critical section from concurrent writes
			fmt.Println("car is available for", name)
			time.Sleep(5 * time.Second) // in real world, this is some kind of latency// network, db calls, or something else
			fmt.Println("booking confirmed", name)
			cab--
		} else {
			fmt.Println("car is not available for", name)
		}
	}()

	// some more parts to the program, which are not locked
}
