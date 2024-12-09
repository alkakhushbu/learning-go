package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

func main() {
	f, err := OpenFile("test.txt")

	if err != nil {
		log.Println(err)
		return
	}

	// trying to read file stats and printing the fileName
	info, _ := f.Stat()
	fmt.Println(info.Name())
}

func OpenFile(name string) (*os.File, error) {
	f, err := os.Open(name)

	if err != nil {
		//errors.Is can check if an error was wrapped inside the chain or not
		//  if an error was found in the chain, you now know what exactly went wrong
		// you might want to take some actions to fix the issue
		//or maybe just log the additional details
		if errors.Is(err, os.ErrNotExist) {
			// attempting to create a file
			f, err := os.Create(name)
			// if it still fails, we will return the original error
			if err != nil {
				return nil, err
			}

			// success case // errors.Is succeeded in identifying the root cause of the problems
			return f, nil
		}
		return nil, err
	}

	return f, nil
}
