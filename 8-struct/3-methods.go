package main

import "fmt"

type Student struct {
	name string
	age  int
}

//func (receiver) funcSignature {//body}
// methods could be only called using struct variable

func (s Student) SayHello() {
	fmt.Println("Hello, my name is", s.name)
}

func main() {
	s := Student{"Jack", 20}
	s.SayHello()

}
