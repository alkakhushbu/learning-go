package main

import (
	"fmt"
	"sync"
)

func main() {
	// FizzBuzz()
	FizzBuzz()
}
func FizzBuzz() {
	for i := 1; i <= 10; i++ {
		switch {
		case i%3 == 0:
			fmt.Print("Fizz")
			fallthrough
		case i%5 == 0:
			fmt.Print("Buzz")
		default:
			fmt.Print(i)
		}
		fmt.Println()
	}
}

func unbuffered() {
	ch := make(chan int, 1)
	// a := <-ch
	ch <- 10
	fmt.Println(<-ch)
}
func unbufferedFixed() {
	ch := make(chan int, 1)
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func(){
		defer wg.Done()
		fmt.Println("Sender picked up")
		ch <- 1
	}()
	go func() {
		defer wg.Done()
		fmt.Println("Receiver picked up")
		fmt.Println(<-ch)
	}()
	wg.Wait()
}

// func FizzBuzz() {
// 	for i := 1; i <= 20; i++ {
// 		switch {
// 		case i%3 == 0 && i%5 == 0:
// 			fmt.Print("FizzBuzz")
// 		case i%3 == 0:
// 			fmt.Print(i, "Fizz")
// 		case i%5 == 0:
// 			fmt.Print(i, "Buzz")
// 		default:
// 			fmt.Print(i)
// 		}
// 		fmt.Println()
// 	}
// }

