package main

import (
	"log"
	"strconv"
)

//Create a function that converts a string to float64, if it is successful it returns the value otherwise an error

func main() {

	input := "20.33"
	fltval, err := convStrToFloat64(input)
	if err != nil {
		//log.error - info
		log.Println("There was an error with input : ", input, " ,error message:", err)
		return
	}
	log.Println("The converted float value is ", fltval)
}

func convStrToFloat64(val string) (float64, error) { // err must be the last value to be returned

	flt, err := strconv.ParseFloat(val, 64)
	if err != nil {
		//log statement here not a good idea, repeated logs problem
		return 0, err // whenever err happens set other values to default
	}
	// success case should be clearly visible
	return flt, nil // don't write err , even err is going to be nil
}
