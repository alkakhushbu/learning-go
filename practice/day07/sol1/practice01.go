package main

import "fmt"

type book struct {
	title       string
	author      string
	pagesToRead int
}

func main() {
	b := book{title: "Learning Go", author: "Jon Bodner", pagesToRead: 496}
	for b.pagesToRead != 0 {
		b.read(100)
		fmt.Println(b)
	}
}

func (b *book) read(readCount int) {
	if b.pagesToRead > readCount {
		b.pagesToRead -= readCount
	} else {
		b.pagesToRead = 0
	}
}
