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
	fmt.Println("main function doing further work")
	//time.Sleep(time.Second * 3)
	fmt.Println("main function done")
	wg.Wait()

}

// ctx must be the first argument in the function,
// ctx should not be part of the struct, but it should be passed to function as an argument
func doSomething(ctx context.Context, name string, wg *sync.WaitGroup) {

	ch := make(chan int, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		x := slowFunction()
		fmt.Println(x)
		ch <- x

	}()

	select {
	case v := <-ch:
		fmt.Println("value received from slow function", v)
	case <-ctx.Done():
		fmt.Println(ctx.Err())
		return
	}

	fmt.Println("doSomething done")

}

func slowFunction() int {
	time.Sleep(4 * time.Second)
	fmt.Println("slow fn ran and add 100 records to db")
	fmt.Println("receiver should process it")
	return 42
}
