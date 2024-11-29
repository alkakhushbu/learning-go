package postgres

import (
	"fmt"
	"question1/stores/models"
)

type Postgres struct {
	userMap map[int]models.User
	lastId  int
}

func NewConnection() *Postgres {
	mySql := &Postgres{userMap: make(map[int]models.User, 5)}
	return mySql
}

func (p *Postgres) Create(user models.User) bool {
	lastId := p.lastId + 1
	p.userMap[lastId] = user
	p.lastId = lastId
	return true
}

func (p *Postgres) Update(id int, name string) bool {
	user, ok := p.userMap[id]
	if ok {
		user.Name = name
		p.userMap[id] = user
		return true
	} else {
		fmt.Println("Can't update, user does not exist with id:", id)
		return false
	}
}
func (p *Postgres) Delete(id int) bool {
	_, ok := p.userMap[id]
	if ok {
		delete(p.userMap, id)
		return true
	} else {
		fmt.Println("Can't delete, user does not exist with id:", id)
		return false
	}
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
