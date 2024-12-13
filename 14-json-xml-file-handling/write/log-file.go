package main

import (
	"log"
	"os"
)

func main() {
	filepath := "files/log.txt"
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Println("Error in opening file:", filepath, "Error:", err)
		return
	}
	defer file.Close()
	file.WriteString("\n2024/12/10 18:57:41 File saved successfully!!.\n________________________")
	log.Println("File saved successfully!!")
}
