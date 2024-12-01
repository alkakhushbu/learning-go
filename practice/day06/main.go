package main

import "fmt"

type author struct {
	name  string
	books []book
}

type book struct {
	name string
}

func main() {
	b := book{"Introducing Go"}
	books := []book{b}
	a := author{name: "John Peterson", books: books}
	a.appendBook(b)
	a.printAuthor()
	books = a.getBooks()
	fmt.Println(books)
	fmt.Println(a)

	author := new(author)
	fmt.Println(author)

}

func (a *author) appendBook(b book) {
	a.books = append(a.books, b)
}

func (a *author) printAuthor() {
	fmt.Println(a)
}

func (a *author) getBooks() []book {
	return a.books
}

//use pointer for author,
func appendBookToAuthor(a *author, b book) {
	books := a.getBooks()
	books = append(books, b)
	a.books = books
}

func appendBookToBooks(books *[]book, b book) {
	*books = append(*books, b)
}
