package main

import (
	"fmt"
)

// Polymorphism means that a piece of code changes its behavior depending on the
// concrete data it’s operating on // Tom Kurtz, Basic inventor

// "Don’t design with interfaces, discover them". - Rob Pike

// Bigger the interface weaker the abstraction // Rob Pike

type Reader interface {
	// only method signatures could be added to an interface
	read(b []byte) (int, error)
	// interfaces are automatically implemented if method signature is same as of interfaces
	//write(b []byte) (int, error) // all methods must be implemented to implement the interface
}

type File struct {
	name string
}

func (f File) read(b []byte) (int, error) {
	fmt.Println("reading files and processing them", f.name)
	return 0, nil
}

//func DoWork(f File, io IO) {
//	i, err := f.read(nil)
//	_, _ = i, err
//	i, err = io.read(nil)
//	_, _ = i, err
//
//}

func (f File) Print() {
	fmt.Println("printing file", f.name)
}

// DoWork can accept any type that implements the interface
func DoWork(r Reader) {

	//piece of code changes its behavior depending on the
	//// concrete data it’s operating on
	// if file is passed Read would be called from the file struct, or if Io is passed Read would be called from IO
	i, err := r.read(nil)
	_, _ = i, err
	// we can only access methods which are part of interface with interface variable
	//r.Print()
	fmt.Printf("%T\n", r) // by printing type of an interface we can see what concrete value interface is storing
	// type assertion // we can fetch a concrete type from an interface variable
	// checking if file struct is present in the interface and doing some file specific task

	f, ok := r.(File) // r.(File) // type assertion // checking if a type is present in the interface or not
	if ok {
		f.Print() // f is a concrete file object not an interface variable
	}

}

type IO struct {
	name string
}

func (i IO) read(b []byte) (int, error) {
	fmt.Println("reading and processing ", i.name)
	return 0, nil
}

func main() {
	f := File{"test.txt"}
	i := IO{"os.stdin"}

	// we can pass a file or io object to DoWork because both objects implement the interface
	DoWork(f)
	DoWork(i)
	//io.Reader()

}
