package main

import (
	"fmt"
	"sync"
)

// go run -race 4-races.go // run your program with race detector
// it should be used in dev environment, not in production
var x int = 1

// create a data-race for the map and detect with race detector
var user map[int]string

func main() {
	wg := new(sync.WaitGroup)
	m := new(sync.Mutex)
	//wg.Add(10)
	for i := 1; i <= 5; i++ {
		wg.Add(2)
		go updateX(i, wg, m)
		go UpdateLocal(i, wg)
	}
	wg.Wait()
}

func updateX(val int, wg *sync.WaitGroup, m *sync.Mutex) {
	defer wg.Done()
	// critical section
	// this is the place where we access the shared resource

	// when a goroutine acquires a lock, another goroutine can't access the critical section
	// until the lock is not released

	//if we don't protect critical section with locks then data race can happen
	//if any goroutine tries to write to the shared resource

	// data race situations
	//	- at least one concurrent write and n number of reads
	//	- n number of concurrent writes
	// 	- n number of concurrent writes and n number of concurrent reads
	// 	Note - Data race doesn't happen if there are only concurrent reads
	m.Lock()
	defer m.Unlock() // releasing the lock
	x = val
	fmt.Println("updateX", x)

}

func UpdateLocal(val int, wg *sync.WaitGroup) {
	defer wg.Done()
	var abc int // abc is a local variable //
	// if we run 10 goroutines, then 10 stack frames would be created,
	// every update would happen in the local stack frame, nothing shared with other goroutines
	abc = val
	fmt.Println("UpdateLocal", abc)
}
