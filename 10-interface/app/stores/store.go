package stores

import "app/stores/models"

type DataBase interface {
	Create(u models.User) (models.User, bool)
	//CreateSimple(string) (models.User, bool)
	Update(int, string) (models.User, bool)
	Delete(int) bool
	FetchAll() (map[int]models.User, bool)
	FetchUser(int) (models.User, bool)
}

//var DB DataBase // not recommended // if someone changes this global value, then it is changed for all

type Store struct {
	DataBase
}

func NewStore(db DataBase) *Store {
	if db == nil {
		panic("db is nil")
	}
	return &Store{DataBase: db}
}
