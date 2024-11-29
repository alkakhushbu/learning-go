package stores

type User struct {
	Name string
}

type Database interface {
	Create(User) bool
	Update(int, string) bool
	Delete(int) bool
	FetchAll()
	FetchUser(int) bool
}
