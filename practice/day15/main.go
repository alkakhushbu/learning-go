package main

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

/*
q1. Create a function that converts string to int
    use time.Sleep to make this function slow
    Pass context to this function with a certain timeout
    If timeout happens this function should report an error back to main
    If timeout didn't happen, this function should return the actual result of Atoi
*/

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	val, err := stringToInt2(ctx, "123")
	if err != nil {
		fmt.Println("main: Cannot convert string into integer:", err)
		return
	}
	fmt.Println("main: String converted to integer:", val)
	// intError := stringToInt(ctx, "12334.5")
	// if intError.err != nil {
	// 	fmt.Println("main: Cannot convert string into integer:", intError.err)
	// 	return
	// }
	// fmt.Println("main: String converted to integer:", intError.val)
}

type IntError struct {
	val int
	err error
}

func stringToInt2(ctx context.Context, s string) (int, error) {
	time.Sleep(time.Second * 5)
	val, err := strconv.Atoi(s)
	select {
	case <-ctx.Done():
		fmt.Println("Cancel convert function call")
		return 0, ctx.Err()
	default:
		if err != nil {
			return 0, err
		}
		return val, nil
	}
}

// func stringToInt(ctx context.Context, s string) IntError {
// 	ch := make(chan IntError)
// 	go func(string) {
// 		time.Sleep(time.Second * 2)
// 		val, err := strconv.Atoi(s)
// 		ch <- IntError{val: val, err: err}
// 		select {
// 		case <-ctx.Done():
// 			fmt.Println("goroutine1: Context timeout, string could not be converted into integer")
// 		default:
// 			fmt.Println("goroutine1: String converted into integer")
// 			return
// 		}
// 	}(s)
// 	select {
// 	case <-ctx.Done():
// 		fmt.Println("stringToInt: Context timeout, string could not be converted into integer")
// 		return IntError{val: 0, err: errors.New("stringToInt: Context timeout, string could not be converted into integer")}
// 	case v := <-ch:
// 		fmt.Println("stringToInt: String converted into integer")
// 		return v
// 	}
// }
