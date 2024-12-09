package main

import (
	"fmt"
	"strconv"
)

var ErrInvalidValue = fmt.Errorf("invalid value")

func main() {
	//var err error
	fmt.Println(strconv.ParseInt("1Q", 10, 64))
	fmt.Println(strconv.Atoi("xyz"))
	fmt.Println(strconv.Atoi("abc"))
	fmt.Println(strconv.ParseBool("2W"))
	fmt.Println(strconv.ParseUint("3W", 10, 64))
}
