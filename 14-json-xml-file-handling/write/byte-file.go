package main

import (
	"fmt"
	"os"
)

func main() {
	filepath := "files/data.bin"
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error in opening file:", filepath, "Error:", err)
		return
	}
	defer file.Close()
	data := []byte{0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x2C, 0x20, 0x57, 0x6F, 0x72, 0x6C, 0x64, 0x21, 0x0A}
	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Error in saving file:", filepath, "Error:", err)
	}
	fmt.Println("File saved successfully!")
}
