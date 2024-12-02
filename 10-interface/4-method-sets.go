package main

import "fmt"

type Speaker interface {
	Speak() string
	Speak2() string
}

type Person struct {
	name string
}

func (p *Person) Speak() string {
	return fmt.Sprintf("Hi, my name is %s", p.name)
}
func (p *Person) Speak2() string {
	return fmt.Sprintf("Hi, my name is 2 %s", p.name)
}

func main() {
	p := Person{"John"}

	var s1 Speaker = &p
	fmt.Println(s1.Speak())
	fmt.Println(s1.Speak2())
}

/*
                          +-------------------------------------+
                          |              Method Set             |
                          +------------------+------------------+
                          |    Value Type    |   Pointer Type   |
+-------------------------+------------------+------------------+
| Function w/ Value Rec.  |        Yes       |       Yes        |
| ( func (t T) )          |                  |                  |
+-------------------------+------------------+------------------+
| Function w/ Ptr Rec.    |        No        |       Yes        |
| ( func (t *T) )         |                  |                  |
+-------------------------+------------------+------------------+

If a function is implemented with a value receiver (func (t T)), it can be called through a value or a pointer.
If a function is implemented with a pointer receiver (func (t *T)), it can only be called through a pointer.
*/
