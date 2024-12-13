package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type configuration struct {
	DBHost        string `json:"database_host"`
	DBPort        int    `json:"database_port"`
	DBUser        string `json:"database_username"`
	DBPassword    string `json:"-"`
	ServerPort    int    `json:"server_port"`
	ServerDebug   bool   `json:"server_debug"`
	ServerTimeout int    `json:"server_timeout"`
}

// read the json, written at data.json file
// use json.Unmarshal to convert the byte data to a struct
// os.ReadFile, Scanner, (os.OpenFile -> f.Read)
// struct fields must be exported, so json package can work on it

func main() {
	filepath := "read/data.json"
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var dataBytes []byte
	for scanner.Scan() {
		dataBytes = append(dataBytes, scanner.Bytes()...)
		// fmt.Println(string(dataBytes))
	}
	if err = scanner.Err(); err != nil {
		fmt.Println("Error scanning data:", err)
		return
	}
	var config configuration
	err = json.Unmarshal(dataBytes, &config)
	if err != nil {
		fmt.Println("Error unmarshaling data:", err)
		return
	}
	fmt.Println(config)
}
