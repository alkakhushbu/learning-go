package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	filepath := "read/data.txt"
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error found:", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	for _, line := range lines {
		fmt.Println(line)
	}
}
