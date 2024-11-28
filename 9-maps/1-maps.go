package main

import "fmt"

func main() {
	// keys in map // anytype which can be compared using == can be used as map key
	// value could be of any valid type
	// maps are not concurrency safe
	// provide a len if you already know an estimate of how many values are going to be store in the map

	dictionary := make(map[string]string, 1)
	dictionary["up"] = "above"
	fmt.Println(dictionary)
	fmt.Println(dictionary["up"]) // accessing value using the key

	dictionary["below"] = "down"
	fmt.Println(dictionary)

	// by default map is nil // map is a pointer under the hood
	//var user map[int]string
	//user[1] = "John" // panic because map is not initialized with any memory
	//fmt.Println(user)

	for key, value := range dictionary {
		// Note:- ranging over a map doesn't return values in the same order
		fmt.Println(key, value)
	}

	fmt.Println(len(dictionary))
	delete(dictionary, "down")
	fmt.Println(dictionary)
	modifyMap(dictionary) // maps are pointers
	fmt.Println(dictionary)

}

func modifyMap(dictionary map[string]string) {
	dictionary["up"] = "abc"
	dictionary["random"] = "random" // it will add the values to the original map, as map is a pointer under the hood
}
