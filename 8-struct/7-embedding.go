package main

import "fmt"

type User struct {
	name  string
	email string
}

func (u *User) UpdateEmail(email string) {
	u.email = email
}

func (u *User) Print() {
	fmt.Println(u)
}

type BookAuthor struct {
	User // embedding user struct to BookAuthor// anonymous field // a field without Name
	// anonymous field name would be same as type name
	bio   string
	books []string
}

func (b *BookAuthor) Print() {
	fmt.Println(b)
}

func (b *BookAuthor) AddBook(book string) {
	b.books = append(b.books, book)
}

type Actor struct {
	u      User // not embedding // creating a field of type User
	movies []string
}

func (a *Actor) Print() {
	fmt.Println(a)
}

func (a *Actor) AddMovie(movie string) {
	a.movies = append(a.movies, movie)
}

func main() {

	u := User{name: "Raj", email: "raj@email.com"}
	b := BookAuthor{
		User:  u,
		bio:   "random book bio",
		books: []string{"b1", "b2"},
	}

	a := Actor{
		u:      User{name: "ajay", email: "raj@email.com"},
		movies: []string{"m1"},
	}

	// using embedding, we can directly access methods of embedded struct
	b.UpdateEmail("raj@gmail.com")
	b.User.Print() // this is specifically accessing the method of User type

	// u is a field of actor struct , we need to access it first to call the updateEmail method
	a.u.UpdateEmail("ajay@gmail.com")
	b.Print()
	a.Print()

}
