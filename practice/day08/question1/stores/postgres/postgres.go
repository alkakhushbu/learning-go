package postgres

import (
	"fmt"
	"question1/stores"
)

type Postgres struct {
	userMap map[int]stores.User
	lastId  int
}

func NewConnection() *Postgres {
	mySql := &Postgres{userMap: make(map[int]stores.User, 5)}
	return mySql
}

func (p *Postgres) Create(user stores.User) bool {
	lastId := p.lastId + 1
	p.userMap[lastId] = user
	p.lastId = lastId
	return true
}

func (p *Postgres) Update(id int, name string) bool {
	user := p.userMap[id]
	user.Name = name
	p.userMap[id] = user
	return true
}
func (p *Postgres) Delete(id int) bool {
	delete(p.userMap, id)
	return true
}

func (p *Postgres) FetchAll() {
	fmt.Println(p.userMap)
}

func (m *Postgres) FetchUser(id int) bool {
	user, ok := m.userMap[id]
	if ok {
		fmt.Println(user)
		return true
	} else {
		fmt.Println("User not found with id:", id)
		return false
	}
}
