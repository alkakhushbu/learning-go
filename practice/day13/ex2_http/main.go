package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	res, err := http.Get("http://www.google.com")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	bytes = bytes[:256]
	fmt.Println(string(bytes))
}
