package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	// Context is an interface, Background method returns an implementation of that interface
	ctx := context.Background() // empty container for storing context values
	wg := &sync.WaitGroup{}

	// Note:- if a ctx is already available do not create a new empty context
	// use the existing one
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel() // clean up the resources taken up by the context

	doSomething(ctx, "John", wg)
	//fmt.Println("main function doing further work")
	////time.Sleep(time.Second * 3)
	//fmt.Println("main function done")
	wg.Wait()

}

// ctx must be the first argument in the function,
// ctx should not be part of the struct, but it should be passed to function as an argument
func doSomething(ctx context.Context, name string, wg *sync.WaitGroup) {

	ch := make(chan int)

	wg.Add(1)
	go func() {
		defer wg.Done()
		x := slowFunction(ctx)
		fmt.Println("goroutine:", x)
		select {
		case ch <- x:
			fmt.Println("goroutine: value sent to the channel")
		case <-ctx.Done():
			fmt.Println("goroutine: timeout happened, cant send value to channel")
			fmt.Println(ctx.Err())
			fmt.Println()
			return
		}
		//ch <- x // this would deadlock if receiver is not available because this is unbuffered channel

	}()

	select {
	// selecting over the channel
	// if value is ready inside the channel ch, then we would receive the value
	// if timer is over then ctx.Done case would execute
	case v := <-ch:
		// received value from the channel
		fmt.Println("doSomething: value received from slow function", v)
	case <-ctx.Done():
		// timeout happened, no point of sending the value over the channel, receiver is not present
		fmt.Println("doSomething: timeout happened", "can't receive value from slow function")
		fmt.Println()
		return
	}

	fmt.Println("doSomething done")

}

func slowFunction(ctx context.Context) int {
	time.Sleep(4 * time.Second)
	fmt.Println("slowFunction: slow fn ran and add 100 records to db")
	fmt.Println()
	select {
	//checking if ctx timer was over
	// if yes, then we would rollback the operation
	case <-ctx.Done():
		fmt.Println("slowFunction: timeout happened", "reversing the operation")
		fmt.Println("slowFunction: rollback the operation")
		fmt.Println()
		return 0
	default:
		// default would be true if timeout didn't happen
		// without a default the select would always block
		// with default select default behaviour is changed to unblocking
		return 42

	}

}
