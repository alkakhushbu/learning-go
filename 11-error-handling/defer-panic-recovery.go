package main

import "fmt"

// panic is a runtime exception
// we need to decide at what leve we need to stop panic
// the level we stop panic must stop or that function must stop that is recovering the panic

// if the caller function doesn't depends on the called function you can stop panic propagation back
// by calling the recovery function in defer in the function that is getting called

// defer guarantees to run // so it would recover the panic if it would happen

func main() {
	// RecoverPanic would recover the current function from panic, but the function needs to stop
	// it can't continue executing
	defer recoverPanic()
	s := make([]int, 1, 10)
	updateSlice(s, 1, 10)
	fmt.Println(s)
	fmt.Println("end of main")
}

func recoverPanic() {

	// The built-in `recover` function can stop the process of panicking,
	//if it is called within a deferred function.

	// msg would have the actual panic message if that happened
	msg := recover()
	if msg != nil {
		// if this case is true, panic happened
		// If `recover` captured a panic, it returns the panic value.
		// Here we print it.
		fmt.Println("recover from panic")
		fmt.Println("panic msg", msg)

	}
}

func updateSlice(s []int, index int, val int) {

	s[index] = val
}
