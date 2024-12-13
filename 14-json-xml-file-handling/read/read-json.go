package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// read the json, written at data.json file
// use json.Unmarshal to convert the byte data to a struct
// os.ReadFile, Scanner, (os.OpenFile -> f.Read)

func main() {
	readFile()
}

type configuration struct {
	DBHost        string `json:"database_host"`
	DBPort        int    `json:"database_port"`
	DBUser        string `json:"database_username"`
	DBPassword    string `json:"-"`
	ServerPort    int    `json:"server_port"`
	ServerDebug   bool   `json:"server_debug"`
	ServerTimeout int    `json:"server_timeout"`
}

func readFile() {
	filepath := "read/data.json"
	var config configuration
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()
	// scanner := bufio.NewScanner(file)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Invalid json format:", err)
		return
	}
	fmt.Println(config)
}
