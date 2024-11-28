package main

// package main is also known as a binary package,
// this is the only package that can run

import (
	"first-proj-day-2/calc" // moduleName/packageName
	"first-proj-day-2/some_pack"
	"fmt"
)

// func main is required to run a binary package
// https://google.github.io/styleguide/go/best-practices#util-packages
func main() {
	fmt.Println("running the app")
	calc.Add(1, 2)
	calc.Subtract(2, 44)

	some_pack.NotRecommended() // don't use _ in package names

	// don't use util packages in naming
	//ioutil.ReadFile()
	//ioutil.ReadAll()
	//os.ReadFile()
	//io.ReadAll()
}
