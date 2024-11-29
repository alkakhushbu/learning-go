package stores

import "question1/stores/models"

type Database interface {
	Create(models.User) bool
	Update(int, string) bool
	Delete(int) bool
	FetchAll()
	FetchUser(int) bool
}

type Store struct {
	Database //embedded interface
}

func NewStore(db Database) *Store {
	if db == nil {
		panic("Database is nil")
	}
	return &Store{Database: db}
}
