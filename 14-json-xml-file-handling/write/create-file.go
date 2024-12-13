package main

import (
	"fmt"
	"os"
)

func main() {
	filepath := "write/data.txt"
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Println("Error in creatinf file:", filepath, " Error:", err)
	}
	defer file.Close()
	file.WriteString("Hello World!\n")
	file.WriteString("My name is Alka.")
	// file.WriteAt()
	// fmt.Fscan()
	fmt.Println("File created successfully!")
}
