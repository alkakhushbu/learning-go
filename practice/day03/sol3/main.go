package main

import "fmt"

/*
Q3. Create a function that takes a list of users
    This func can append new values to the list or change the existing elems
    But Make sure this function can't modify the original slice
    that was created in the main function
*/
func main() {
	users := []string{"Darcy", "Bingley", "Jane", "Lizzie"}
	newUsers := appendUsers(users, "Lucas", "Mary", "Bennet")
	fmt.Println(newUsers)
}

func appendUsers(list []string, data ...string) []string {
	m := len(list)
	n := len(data)
	totalLen := m + n + 1
	newList := make([]string, totalLen)
	copy(newList, list)
	copy(newList[len(list):], data)
	return newList
}
