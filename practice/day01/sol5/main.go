package main

import (
	"fmt"
)

func main() {
	jsonString := "{\n \"name\": \"Alice\", \n \"age\": 25\n}"
	rawString := `{
    "name": "Alice",
    "age": 25
  }`

	rawFilePath := `C:\Users\Alice\Documents\example.txt`
	filePath := "C:\\Users\\alkakhushbu\\Documents\\example.txt"

	fmt.Printf("Json String: %s\n", jsonString)
	fmt.Printf("File Path: %s\n", filePath)

	fmt.Printf("\n\n Raw String: %s\n", rawString)
	fmt.Printf("Raw File Path: %s\n", rawFilePath)

}
