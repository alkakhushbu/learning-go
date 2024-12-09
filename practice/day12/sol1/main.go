package main

import (
	"fmt"
	"sync"
	"time"
)

/*
q3. Create 4 functions
    Add(int,int),Sub(int,int),Divide(int,int), CollectResults()
    Add,Sub,Divide do their operations and push value to an unbuffered channel

    CollectResult() -> It would receive the values from the channel and print it
*/

func main() {
	ch := make(chan int)
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		defer wg.Done()
		fmt.Println("Receiver picked up")
		ch <- Add(1, 2)
		ch <- Sub(2, 3)
		ch <- Divide(4, 2)
	}()
	go func() {
		defer wg.Done()
		fmt.Println("Sender picked up")
		CollectResults(<-ch)
		CollectResults(<-ch)
		CollectResults(<-ch)
	}()
	wg.Wait()
}

func Add(a, b int) int {
	time.Sleep(time.Second * 5)
	return a + b
}

func Sub(a, b int) int {
	return a - b
}

func Divide(a, b int) int {
	if b == 0 {
		panic("cannot divide by 0")
	}
	return a / b
}

func CollectResults(a int) {
	fmt.Println(a)
}
