package main

import (
	"fmt"
	"log"
	"strconv"
)

func main() {
	// anonymous functions are functions without name
	// they can be called only once
	x, err := func(s string) (int, error) {
		i, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return i, nil
	}("abc") // () this is calling of the function

	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(x)

}
