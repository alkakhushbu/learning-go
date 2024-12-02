package main

import "fmt"

type Runner interface {
	Run()
}
type Walker interface {
	Walk()
}

// embedding two interface to compose a new interface out of it
// if a type needs to implement a RunnerWalker interface then it should implement all the methods from
// embedded interfaces

type RunnerWalker interface {
	Runner
	Walker
}

//

type cat struct {
	name string
}

func (c *cat) Run() {
	fmt.Println("cat is running")
}
func (c *cat) Walk() {
	fmt.Println("cat is walking")
}

func main() {
	c := cat{name: "tom"}
	accept(&c)
}

func accept(walker RunnerWalker) {

}
