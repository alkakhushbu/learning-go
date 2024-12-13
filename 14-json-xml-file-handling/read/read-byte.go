package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	filepath := "read/data.txt"
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error in opening file:", filepath, "Error:", err)
	}
	defer file.Close()
	buffer := make([]byte, 256)
	for {
		count, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Reached end of file")
				break
			}
			fmt.Println("Error found in reading file:", err)
			break
		}
		fmt.Println("Byte count:", count, " Buffer:\n", string(buffer))
		fmt.Println("_____________________________________________")
	}
}
