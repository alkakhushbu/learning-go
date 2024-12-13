package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// read the json, written at data.json file
// use json.Unmarshal to convert the byte data to a struct
// os.ReadFile, Scanner, (os.OpenFile -> f.Read)
// struct fields must be exported, so json package can work on it

// convert json to struct
// turn inline off to create different types for nested json
// https://mholt.github.io/json-to-go/
type user struct {
	FirstName    string          `json:"first_name"`    // json is a field level tag, used by the json package
	PasswordHash string          `json:"password_hash"` // setting name of the field in the json output
	Perms        map[string]bool `json:"perms"`
}

// 	fmt.Println("Parsed data:")
// 	for _, user := range users {
// 		fmt.Printf("user: %+v\n", user)
// 	}
// }

func main() {
	filepath := "data.json"
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
	var u []user
	err = json.Unmarshal(dataBytes, &u)
	if err != nil {
		fmt.Println("Error unmarshaling data:", err)
		return
	}
	fmt.Println(u)
}
