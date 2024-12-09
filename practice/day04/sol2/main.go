package main

import (
	"log"
	"os"
	"strconv"
)

//q2. Create a function that converts a string to float64,
// if it is successful it returns the value otherwise an error

func main() {
	data := os.Args[1:]
	num, err := convertStringToFloat(data[0])
	if err != nil {
		log.Println("There is an error with the input : ", num, " error : ", err)
		return
	}
	log.Println(num)
}

func convertStringToFloat(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0, err
	}
	return f, err
}
