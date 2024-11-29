package mysql

import (
	"fmt"
	"question1/stores"
)


type Mysql struct {
	userMap map[int]stores.User
	lastId  int
}

func NewConnection() *Mysql {
	mySql := &Mysql{userMap: make(map[int]stores.User, 5)}
	return mySql
}

func (m *Mysql) Create(user stores.User) bool {
	lastId := m.lastId + 1
	m.userMap[lastId] = user
	m.lastId = lastId
	return true
}

func (m *Mysql) Update(id int, name string) bool {
	user := m.userMap[id]
	user.Name = name
	m.userMap[id] = user
	return true
}
func (m *Mysql) Delete(id int) bool {
	delete(m.userMap, id)
	return true
}

func (m *Mysql) FetchAll() {
	fmt.Println(m.userMap)
}

// FetchUser implements stores.Database.
func (m *Mysql) FetchUser(id int) bool {
	user, ok := m.userMap[id]
	if ok {
		fmt.Println(user)
		return true
	} else {
		fmt.Println("User not found with id:", id)
		return false
	}
}
