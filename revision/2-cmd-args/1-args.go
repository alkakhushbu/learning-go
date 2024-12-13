package main

import (
	"fmt"
	"os"
)

func main() {
	// start from 1st index, till the end
	fmt.Println(os.Args[1:])

}
